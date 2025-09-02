📰 RSS Feed Aggregator
A simple RSS Feed Aggregator built with Go.
This project lets you:
 - 👤 Register and log in as a user
 - ➕ Add and follow RSS feeds
 - 📖 Fetch and view articles from feeds
 - 🔄 Reset and start fresh anytime

Currently supports XML-based RSS feeds.

Example supported feeds:
 - Hacker News (https://news.ycombinator.com/rss)
 - Boot.dev Blog (https://blog.boot.dev/index.xml)

🚀 Getting Started
1. Clone the repository
```bash
git clone https://github.com/Pradhyumna789/RSS_Feed_Aggregator.git
cd RSS_Feed_Aggregator
```
2. Run commands
The project is CLI-based. You can run commands using:
```bash
go run . <command>
```
📌 Available Commands
👤 User Management
```bash
go run . register "your-username"
go run . login "your-username"
go run . users
go run . reset
```
📡 Feed Management
```bash
go run . addfeed "feed-name" "feed-url"
# ⚠️ Currently supports only XML-based feeds

go run . following
```
⏳ Aggregating Feeds
```bash
go run . agg 2s
```
👉 You can customize how often feeds are fetched (intervals, number of feeds, etc.) in the scrapeFeeds() function inside rss.go.

📖 Example Usage
```bash
# Register a user
go run . register "alice"

# Login as that user
go run . login "alice"

# Add a feed
go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"

# Followed feeds
go run . following

# Aggregate feeds every 5s
go run . agg 5s
```
