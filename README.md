# 📰 RSS Feed Aggregator

A simple **RSS Feed Aggregator** built with **Go**.  
This project lets you:

- 👤 Register and log in as a user  
- ➕ Add and follow RSS feeds  
- 📖 Fetch and view articles from feeds  
- 🔄 Reset and start fresh anytime  

Currently supports **XML-based RSS feeds**.  

Example supported feeds:  
- [Hacker News](https://news.ycombinator.com/rss)  
- [Boot.dev Blog](https://blog.boot.dev/index.xml)  

---

## 🚀 Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/Pradhyumna789/RSS_Feed_Aggregator.git
cd RSS_Feed_Aggregator
2. Run commands
The project is CLI-based. You can run commands using:

bash
Copy code
go run . <command>
📌 Available Commands
👤 User Management
Register a new user

bash
Copy code
go run . register "your-username"
Log in as an existing user

bash
Copy code
go run . login "your-username"
List all users
Displays all registered users and highlights the currently logged-in user.

bash
Copy code
go run . users
Reset the database
Removes all users, feeds, and subscriptions. Start fresh.

bash
Copy code
go run . reset
📡 Feed Management
Add a feed
Add a new RSS feed by providing a name and URL.

bash
Copy code
go run . addfeed "feed-name" "feed-url"
⚠️ Currently supports only XML-based feeds.

List followed feeds
Shows all feeds followed by the logged-in user.

bash
Copy code
go run . following
⏳ Aggregating Feeds
Fetch new posts from feeds
Example: fetch every 2 seconds

bash
Copy code
go run . agg 2s
👉 You can customize how often feeds are fetched (intervals, number of feeds, etc.) in the scrapeFeeds() function inside rss.go.

🛠️ Future Improvements
✅ Support for Atom and JSON feeds

✅ Persistent storage (DB)

✅ Improved CLI UX (colors, better formatting)

📖 Example Usage
bash
Copy code
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
