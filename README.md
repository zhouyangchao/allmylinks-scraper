# allmylinks-scraper

A Go-based scraper for extracting user information from allmylinks.com profiles.

## Features

- Scrapes user profile information including:
  - Username
  - Avatar URL
  - Display Name
  - Birthday
  - Bio
  - Content
  - Location
  - Profile Views
  - Last Online
  - QR Code URL
  - Links (with title, URL, URL text, and connected status)
- Handles both username and full URL inputs
- Removes duplicate links

## Installation

1. Ensure you have Go 1.22.5 or later installed.
2. Clone the repository:
   ```
   git clone https://github.com/zhouyangchao/allmylinks-scraper.git
   ```
3. Navigate to the project directory:
   ```
   cd allmylinks-scraper
   ```
4. Install dependencies:
   ```
   go mod download
   ```

## Usage

Run the scraper from the command line:
```
go run cmd/allmylinks/main.go <username or full URL>
```

## Example:

```
go run cmd/allmylinks/main.go johndoe
go run cmd/allmylinks/main.go https://allmylinks.com/johndoe
```

## Code Structure

- `cmd/allmylinks/main.go`: Entry point of the application
- `allmylinks/allmylinks.go`: Core scraping logic and data structures

## Dependencies

- golang.org/x/net/html: HTML parsing library

## License

[Add your chosen license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.