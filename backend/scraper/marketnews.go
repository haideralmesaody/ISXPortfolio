package scraper

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Attachment Datatype
type Attachment struct {
	URL      string
	Filename string
	IsLoaded bool
}

// NewsItem Datatype
type NewsItem struct {
	Title       string       `json:"title"`
	Link        string       `json:"link"`
	Date        string       `json:"date"`
	Ticker      string       `json:"ticker"`
	IsNew       bool         `json:"is_new"`
	Attachments []Attachment `json:"attachments"`
}

// News URLs
var newsURLs = []string{
	"http://www.isx-iq.net/isxportal/portal/storyList.html?currLanguage=ar&activeTab=0",
	"http://www.isx-iq.net/isxportal/portal/storyList.html?currLanguage=ar&activeTab=1",
}

// Add new type to manage the entire scraping process
type MarketNewsScraper struct {
	CSVPath       string
	PDFDir        string
	BaseURL       string
	ExistingItems []NewsItem
	NewItems      []NewsItem
	AllItems      []NewsItem
}

// Constructor for our scraper
func NewMarketNewsScraper(csvPath, pdfDir string) *MarketNewsScraper {
	return &MarketNewsScraper{
		CSVPath: "/app/data/market_news.csv",
		PDFDir:  "/app/data/pdfs",
		BaseURL: "http://www.isx-iq.net",
	}
}

// Main method to run the entire process
func (s *MarketNewsScraper) Run() error {
	log.Println("=== Starting Market News Scraper ===")

	// Phase 1: Gathering News Items
	if err := s.gatherNewsItems(); err != nil {
		return fmt.Errorf("error gathering news items: %w", err)
	}

	// Phase 2: Processing Items
	if err := s.processNewsItems(); err != nil {
		return fmt.Errorf("error processing news items: %w", err)
	}

	// Phase 3: Save Results
	if err := s.saveResults(); err != nil {
		return fmt.Errorf("error saving results: %w", err)
	}

	log.Println("=== Market News Scraper Finished ===")
	return nil
}

// Phase 1: Gather news items
func (s *MarketNewsScraper) gatherNewsItems() error {
	log.Println("=== Phase 1: Gathering News Items ===")

	// Read existing items
	existing, err := s.readExistingCSV()
	if err != nil {
		log.Printf("Error reading existing CSV: %v", err)
		existing = []NewsItem{} // Start fresh if error
	}
	s.ExistingItems = existing
	log.Printf("Read %d existing items from CSV", len(existing))

	// Get basic info for all news items
	newItems, err := s.getNewsItemsList()
	if err != nil {
		return fmt.Errorf("error getting news list: %w", err)
	}
	log.Printf("Found %d items on website", len(newItems))

	// Merge and identify new items
	s.AllItems = s.mergeNewsItems(s.ExistingItems, newItems)

	// Count new vs existing
	var newCount int
	for _, item := range s.AllItems {
		if item.IsNew {
			newCount++
		}
	}
	log.Printf("After merge: %d total items (%d new, %d existing)",
		len(s.AllItems), newCount, len(s.AllItems)-newCount)

	return nil
}

// Phase 2: Process items
func (s *MarketNewsScraper) processNewsItems() error {
	log.Println("=== Phase 2: Processing Individual News Items ===")

	// Count new items first
	var newItems int
	for _, item := range s.AllItems {
		if item.IsNew {
			newItems++
		}
	}

	if newItems == 0 {
		log.Println("No new items to process, skipping details retrieval")
		return nil
	}

	log.Printf("Processing %d new items...", newItems)

	for i := range s.AllItems {
		if !s.AllItems[i].IsNew {
			continue // Skip existing items
		}

		log.Printf("Processing new item %d/%d: %s", i+1, newItems, s.AllItems[i].Title)

		if err := s.processNewsItem(&s.AllItems[i]); err != nil {
			log.Printf("Error processing item: %v", err)
			continue
		}
	}

	return nil
}

// Process single news item
func (s *MarketNewsScraper) processNewsItem(item *NewsItem) error {
	log.Printf("Processing item: %s", item.Title)

	// Get details
	if err := s.getNewsItemDetails(item); err != nil {
		return fmt.Errorf("error getting details: %w", err)
	}

	// Process attachments
	if len(item.Attachments) > 0 {
		if err := s.processAttachments(item); err != nil {
			return fmt.Errorf("error processing attachments: %w", err)
		}
	}

	return nil
}

// Phase 3: Save results
func (s *MarketNewsScraper) saveResults() error {
	log.Println("=== Phase 3: Saving Results ===")

	// Sort items
	log.Println("Sorting items by date...")
	s.sortNewsByDateTime()

	// Save to CSV
	if err := s.saveToCSV(); err != nil {
		return fmt.Errorf("error saving to CSV: %w", err)
	}

	// Verify attachments
	log.Println("=== Verifying All Attachments ===")
	s.verifyAllAttachments()

	return nil
}

// GetAllMarketNews scrapes all news from predefined URLs
func GetAllMarketNews() ([]NewsItem, error) {
	log.Println("Starting GetAllMarketNews...")
	var allNewsItems []NewsItem

	for i, url := range newsURLs {
		log.Printf("Processing URL %d of %d: %s", i+1, len(newsURLs), url)
		newsItems, err := ScrapeMarketNews(url)
		if err != nil {
			log.Printf("Error scraping news from %s: %v", url, err)
			continue
		}
		log.Printf("Successfully got %d news items from URL %d", len(newsItems), i+1)
		allNewsItems = append(allNewsItems, newsItems...)
	}

	log.Printf("Finished GetAllMarketNews, total items: %d", len(allNewsItems))
	return allNewsItems, nil
}

// ScrapeMarketNews remains as an internal function
func ScrapeMarketNews(url string) ([]NewsItem, error) {
	log.Printf("Starting to scrape URL: %s", url)

	log.Println("Setting up Edge browser options...")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.ExecPath("C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"),
	)

	log.Println("Creating allocator context...")
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	log.Println("Creating browser context...")
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	log.Println("Starting navigation and table extraction...")
	var rows []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Page loaded, waiting for table to be visible...")
			return nil
		}),
		chromedp.WaitVisible(`.indnews-datarow`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("News rows are visible, getting data...")
			return nil
		}),
		chromedp.Nodes(`.indnews-datarow`, &rows),
	)
	if err != nil {
		log.Printf("Error during navigation/extraction: %v", err)
		return nil, fmt.Errorf("failed to get news rows: %w", err)
	}

	log.Printf("Found %d news items", len(rows))

	var newsItems []NewsItem
	baseURL := "http://www.isx-iq.net/isxportal/portal/"

	for i, row := range rows {
		log.Printf("Processing news item %d...", i)

		var item NewsItem
		err := chromedp.Run(ctx,
			chromedp.Text(".table-newsdata", &item.Date, chromedp.ByQuery, chromedp.FromNode(row)),
			chromedp.Text(".indnews-title", &item.Title, chromedp.ByQuery, chromedp.FromNode(row)),
			chromedp.AttributeValue(".indnews-title a", "href", &item.Link, nil, chromedp.ByQuery, chromedp.FromNode(row)),
		)
		if err != nil {
			log.Printf("Error extracting data from row %d: %v", i, err)
			continue
		}

		if dateEnd := strings.Index(item.Date, "\u00a0"); dateEnd != -1 {
			item.Date = strings.TrimSpace(item.Date[:dateEnd])
		}

		log.Printf("Extracted item %d: Date=%s, Title=%s, Link=%s", i, item.Date, item.Title, item.Link)

		// Get the detail page
		detailURL := baseURL + item.Link
		var detailHTML string
		err = chromedp.Run(ctx,
			chromedp.Navigate(detailURL),
			chromedp.OuterHTML("html", &detailHTML),
		)
		if err != nil {
			log.Printf("Error getting detail page %s: %v", detailURL, err)
			continue
		}

		// Extract ticker
		item.Ticker = ExtractTickerFromHTML(detailHTML)

		// Extract PDF attachments
		pdfLinks := ExtractPDFLinks(detailHTML)
		for _, pdfURL := range pdfLinks {
			filename := strings.TrimPrefix(pdfURL, "/isxportal/files/")
			item.Attachments = append(item.Attachments, Attachment{
				URL:      pdfURL,
				Filename: filename,
				IsLoaded: false, // Will be set to true after downloading
			})
		}

		newsItems = append(newsItems, item)
	}

	log.Printf("Successfully extracted %d news items", len(newsItems))
	return newsItems, nil
}

// SplitDateTime splits a datetime string into date and time components
func SplitDateTime(datetime string) (date, time string) {
	parts := strings.Split(datetime, " ")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return datetime, ""
}

// Add this function to parse the Iraqi date format
func parseDateTime(dateStr string) (time.Time, error) {
	// Format: "31/12/2024 10:14"
	return time.Parse("02/01/2006 15:04", dateStr)
}

// Add this function to sort news items
func SortNewsByDateTime(items []NewsItem) {
	sort.Slice(items, func(i, j int) bool {
		timeI, errI := parseDateTime(items[i].Date)
		timeJ, errJ := parseDateTime(items[j].Date)

		// If there's an error parsing either date, move the error item to the end
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		// Sort in descending order (newest first)
		return timeI.After(timeJ)
	})
}

// Add to your existing NewsItem type
func (n NewsItem) Equals(other NewsItem) bool {
	return n.Title == other.Title && n.Link == other.Link
}

// MergeNewsItems combines two slices of news items and removes duplicates
func MergeNewsItems(existing, new []NewsItem) []NewsItem {
	seen := make(map[string]bool)
	var merged []NewsItem

	// Helper function to add items if not seen
	addIfNotSeen := func(item NewsItem, isNew bool) {
		if !seen[item.Link] {
			seen[item.Link] = true
			item.IsNew = isNew
			merged = append(merged, item)
		}
	}

	// Add new items first (they take precedence)
	for _, item := range new {
		addIfNotSeen(item, true)
	}

	// Add existing items that weren't in new items
	for _, item := range existing {
		addIfNotSeen(item, false)
	}

	return merged
}

// Add function to extract ticker from HTML
func ExtractTickerFromHTML(html string) string {
	// Look for company profile link with companyCode parameter
	re := regexp.MustCompile(`companyprofilecontainer\.html\?companyCode=([A-Z]+)`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Add function to extract PDF links
func ExtractPDFLinks(html string) []string {
	re := regexp.MustCompile(`/isxportal/files/story[0-9]+_[0-9_]+\.pdf`)
	matches := re.FindAllString(html, -1)
	return matches
}

// DownloadPDF downloads a PDF file and saves it to the specified directory
func DownloadPDF(baseURL, pdfURL, outputDir string) error {
	fullURL := baseURL + pdfURL
	log.Printf("Downloading PDF from: %s", fullURL)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get the filename from the URL
	filename := filepath.Base(pdfURL)
	outputPath := filepath.Join(outputDir, filename)

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		log.Printf("PDF already exists: %s", outputPath)
		return nil
	}

	// Download the file
	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write PDF to file: %w", err)
	}

	log.Printf("Successfully downloaded PDF to: %s", outputPath)
	return nil
}

// GetNewsItemsList gets basic info for all news items without details
func GetNewsItemsList() ([]NewsItem, error) {
	log.Println("Getting list of all news items...")
	var allNewsItems []NewsItem

	for i, url := range newsURLs {
		log.Printf("Processing URL %d of %d: %s", i+1, len(newsURLs), url)
		newsItems, err := getNewsItemsFromPage(url)
		if err != nil {
			log.Printf("Error getting news from %s: %v", url, err)
			continue
		}
		log.Printf("Got %d items from URL %d", len(newsItems), i+1)
		allNewsItems = append(allNewsItems, newsItems...)
	}

	return allNewsItems, nil
}

// getNewsItemsFromPage gets basic info from a single page
func getNewsItemsFromPage(url string) ([]NewsItem, error) {
	// Setup Chrome
	ctx, cancel := setupChrome()
	defer cancel()

	var rows []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.indnews-datarow`),
		chromedp.Nodes(`.indnews-datarow`, &rows),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get news rows: %w", err)
	}

	var items []NewsItem
	for _, row := range rows {
		var item NewsItem
		err := chromedp.Run(ctx,
			chromedp.Text(".table-newsdata", &item.Date, chromedp.ByQuery, chromedp.FromNode(row)),
			chromedp.Text(".indnews-title", &item.Title, chromedp.ByQuery, chromedp.FromNode(row)),
			chromedp.AttributeValue(".indnews-title a", "href", &item.Link, nil, chromedp.ByQuery, chromedp.FromNode(row)),
		)
		if err != nil {
			continue
		}

		if dateEnd := strings.Index(item.Date, "\u00a0"); dateEnd != -1 {
			item.Date = strings.TrimSpace(item.Date[:dateEnd])
		}

		items = append(items, item)
	}

	return items, nil
}

// GetNewsItemDetails gets full details for a single news item
func GetNewsItemDetails(item *NewsItem) error {
	ctx, cancel := setupChrome()
	defer cancel()

	baseURL := "http://www.isx-iq.net/isxportal/portal/"
	detailURL := baseURL + item.Link

	log.Printf("Getting details for news item: %s", item.Title)

	var detailHTML string
	err := chromedp.Run(ctx,
		chromedp.Navigate(detailURL),
		chromedp.OuterHTML("html", &detailHTML),
	)
	if err != nil {
		return fmt.Errorf("failed to get detail page: %w", err)
	}

	// Extract ticker
	item.Ticker = ExtractTickerFromHTML(detailHTML)
	log.Printf("Found ticker: %s", item.Ticker)

	// Extract PDF attachments
	pdfLinks := ExtractPDFLinks(detailHTML)
	log.Printf("Found %d PDF attachments", len(pdfLinks))

	for i, pdfURL := range pdfLinks {
		filename := strings.TrimPrefix(pdfURL, "/isxportal/files/")
		log.Printf("Adding attachment %d/%d: %s", i+1, len(pdfLinks), filename)
		item.Attachments = append(item.Attachments, Attachment{
			URL:      pdfURL,
			Filename: filename,
			IsLoaded: false,
		})
	}

	return nil
}

// Helper function to setup Chrome
func setupChrome() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.ExecPath("/usr/bin/chromium"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	return ctx, cancel
}

// VerifyAttachments checks if all PDFs were downloaded correctly
func VerifyAttachments(item NewsItem, pdfDir string) []string {
	var missing []string

	log.Printf("Verifying attachments for: %s", item.Title)
	for i, att := range item.Attachments {
		path := filepath.Join(pdfDir, att.Filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("Warning: Attachment %d/%d missing: %s",
				i+1, len(item.Attachments), att.Filename)
			missing = append(missing, att.Filename)
		} else {
			log.Printf("Verified attachment %d/%d: %s",
				i+1, len(item.Attachments), att.Filename)
		}
	}

	return missing
}

// Add these methods to MarketNewsScraper

func (s *MarketNewsScraper) readExistingCSV() ([]NewsItem, error) {
	// Check if file exists
	if _, err := os.Stat(s.CSVPath); os.IsNotExist(err) {
		return nil, nil
	}

	file, err := os.Open(s.CSVPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	var items []NewsItem
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	portalURL := s.BaseURL + "/isxportal/portal/"
	for _, record := range records {
		if len(record) < 6 {
			continue
		}

		link := record[3]
		link = strings.TrimPrefix(link, portalURL)

		item := NewsItem{
			Date:   record[0] + " " + record[1],
			Title:  record[2],
			Link:   link,
			Ticker: record[4],
		}

		if len(record) > 6 {
			attachmentsStr := record[6]
			if attachmentsStr != "" {
				attachmentsList := strings.Split(attachmentsStr, ";")
				for _, attStr := range attachmentsList {
					parts := strings.Split(attStr, "|")
					if len(parts) == 3 {
						item.Attachments = append(item.Attachments, Attachment{
							URL:      parts[0],
							Filename: parts[1],
							IsLoaded: parts[2] == "true",
						})
					}
				}
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *MarketNewsScraper) getNewsItemsList() ([]NewsItem, error) {
	return GetNewsItemsList()
}

func (s *MarketNewsScraper) mergeNewsItems(existing, new []NewsItem) []NewsItem {
	// Use map to track existing items by link
	existingMap := make(map[string]NewsItem)
	for _, item := range existing {
		existingMap[item.Link] = item
	}

	var merged []NewsItem
	var newItems []NewsItem

	// First check which items are actually new
	for _, item := range new {
		if existingItem, exists := existingMap[item.Link]; exists {
			// Item exists, keep the existing one with its attachments
			item = existingItem
			item.IsNew = false
		} else {
			// This is a new item
			item.IsNew = true
			newItems = append(newItems, item)
		}
		merged = append(merged, item)
	}

	// Log new items found
	if len(newItems) > 0 {
		log.Printf("Found %d new items:", len(newItems))
		for _, item := range newItems {
			log.Printf("- %s (%s)", item.Title, item.Date)
		}
	} else {
		log.Println("No new items found")
	}

	return merged
}

func (s *MarketNewsScraper) getNewsItemDetails(item *NewsItem) error {
	return GetNewsItemDetails(item)
}

func (s *MarketNewsScraper) processAttachments(item *NewsItem) error {
	log.Printf("Processing %d attachments for: %s", len(item.Attachments), item.Title)

	// Create PDFs directory if it doesn't exist
	if err := os.MkdirAll(s.PDFDir, 0755); err != nil {
		return fmt.Errorf("error creating PDF directory: %w", err)
	}

	for i, att := range item.Attachments {
		if att.IsLoaded {
			log.Printf("Attachment %d/%d already loaded: %s",
				i+1, len(item.Attachments), att.Filename)
			continue
		}

		log.Printf("Downloading attachment %d/%d: %s",
			i+1, len(item.Attachments), att.URL)

		err := DownloadPDF(s.BaseURL, att.URL, s.PDFDir)
		if err != nil {
			log.Printf("Error downloading attachment %d/%d: %v",
				i+1, len(item.Attachments), err)
			continue
		}

		item.Attachments[i].IsLoaded = true
		log.Printf("Successfully downloaded attachment %d/%d: %s",
			i+1, len(item.Attachments), att.Filename)
	}
	return nil
}

func (s *MarketNewsScraper) sortNewsByDateTime() {
	sort.Slice(s.AllItems, func(i, j int) bool {
		timeI, _ := time.Parse("02/01/2006 15:04", s.AllItems[i].Date)
		timeJ, _ := time.Parse("02/01/2006 15:04", s.AllItems[j].Date)
		return timeI.After(timeJ)
	})
}

func (s *MarketNewsScraper) saveToCSV() error {
	log.Printf("Writing %d items to CSV: %s", len(s.AllItems), s.CSVPath)

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("/app/data", 0755); err != nil {
		return fmt.Errorf("error creating data directory: %w", err)
	}

	file, err := os.Create(s.CSVPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write UTF-8 BOM
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Date", "Time", "Description", "Link", "Ticker", "Is New", "Attachments"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	for i, item := range s.AllItems {
		date, time := SplitDateTime(item.Date)

		var attachmentsStr []string
		for _, att := range item.Attachments {
			attachmentsStr = append(attachmentsStr,
				fmt.Sprintf("%s|%s|%v", att.URL, att.Filename, att.IsLoaded))
		}
		attachments := strings.Join(attachmentsStr, ";")

		row := []string{
			date,
			time,
			item.Title,
			s.BaseURL + "/isxportal/portal/" + item.Link,
			item.Ticker,
			fmt.Sprintf("%v", item.IsNew),
			attachments,
		}

		if err := writer.Write(row); err != nil {
			log.Printf("Error writing row %d: %v", i, err)
			continue
		}
		log.Printf("Wrote item %d to CSV", i)
	}

	return nil
}

func (s *MarketNewsScraper) verifyAllAttachments() {
	for _, item := range s.AllItems {
		if missing := VerifyAttachments(item, s.PDFDir); len(missing) > 0 {
			log.Printf("Missing attachments for %s: %v", item.Title, missing)
		}
	}
}
