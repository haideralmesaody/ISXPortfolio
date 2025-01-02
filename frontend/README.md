# ISXPortfolio

ISXPortfolio is a web application that provides users with tools and information to manage their investments in the Iraq Stock Exchange (ISX).

## Features

* **Market Overview:**
    * View the latest news from the Iraq Stock Market.
    * Browse market tickers with real-time data.
    * Access basic ticker information (news, candle chart, financial analysis).
* **Portfolio Management:**
    * Create and manage up to 3 portfolios.
    * Record and track transactions.
    * Generate performance reports with time-weighted or money-weighted returns.
    * Monitor portfolio growth and individual position weights.
* **Advanced Analysis:**
    * Access detailed technical and financial analysis for each company.
    * Get AI-powered investment suggestions.
* **Watchlist:**
    * Create and manage a watchlist of stocks you're interested in.
* **Stock Screener:**
    * Filter stocks based on specific criteria.
* **Price Alerts:**
    * Set up alerts for price movements in your portfolio or watchlist.

## Future Development

* **Social Community:** Connect and share insights with other investors.
* **Risk Analysis:** Assess the risk profile of your portfolio.
* **Goal Setting and Tracking:** Set financial goals and track your progress.

## Technology Stack

* **Backend:** Go
* **Frontend:** Flutter
* **Database:** SQLite (potentially migrating to a more scalable solution like PostgreSQL in the future)
* **Deployment:** Google Cloud Platform (GCP) or Amazon Web Services (AWS)
* **Authentication:** Google Authentication

## Getting Started

1. Clone the repository: `git clone https://github.com/haideralmesaody/isxportfolio.git`
2. Install dependencies:
   * Backend: `cd backend && go mod download`
   * Frontend: `cd frontend && flutter pub get`
3. Set up environment variables:
   * Create a `.env` file in the `backend` directory.
   * Refer to `.env.example` for the required variables.
4. Run the application:
   * Backend: `go run main.go`
   * Frontend: `flutter run -d web`

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines.

## License

This project is licensed under the [MIT License](LICENSE).