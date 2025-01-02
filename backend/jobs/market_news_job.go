package jobs

import (
	"isxportfolio-backend/scraper"
	"log"
	"time"
)

type MarketNewsJob struct {
	scraper *scraper.MarketNewsScraper
	done    chan bool
}

func NewMarketNewsJob() *MarketNewsJob {
	return &MarketNewsJob{
		scraper: scraper.NewMarketNewsScraper(
			"/app/data/market_news.csv",
			"/app/data/pdfs",
		),
		done: make(chan bool),
	}
}

func (j *MarketNewsJob) Start() {
	go j.run()
}

func (j *MarketNewsJob) run() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Run immediately if within business hours
	if j.isBusinessHours() {
		log.Println("Initial market news update...")
		if err := j.scraper.Run(); err != nil {
			log.Printf("Error in initial market news update: %v", err)
		}
	}

	for {
		select {
		case <-j.done:
			return
		case t := <-ticker.C:
			if j.isBusinessHours() {
				log.Printf("Running scheduled market news update at %s...", t.Format("15:04:05"))
				if err := j.scraper.Run(); err != nil {
					log.Printf("Error updating market news: %v", err)
				}
			}
		}
	}
}

func (j *MarketNewsJob) Stop() {
	if j.done != nil {
		j.done <- true
		close(j.done)
	}
}

// isBusinessHours checks if current time is within trading hours
func (j *MarketNewsJob) isBusinessHours() bool {
	now := time.Now()

	// Check if it's Friday (5) or Saturday (6)
	if now.Weekday() == time.Friday || now.Weekday() == time.Saturday {
		return false
	}

	// Get current hour in Baghdad time
	loc, err := time.LoadLocation("Asia/Baghdad")
	if err != nil {
		log.Printf("Error loading Baghdad timezone: %v, using local time", err)
		loc = time.Local
	}
	baghdadTime := now.In(loc)
	hour := baghdadTime.Hour()

	// Check if within 9 AM to 3 PM (15:00)
	return hour >= 9 && hour < 15
}
