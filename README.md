# RSS_Feed_Aggregator

# Usage:

- git clone the project
- go run . register "enter your username"
- go run . login "enter your username"
- go run . users => lists all the registered users and the currently logged in user
- go run . reset => Resets the whole database and you are now back to square one
- go run . addfeed "enter your feed's name" "enter your feed's url" => Currently the parser only supports parsing for "xml" content type
- go run . following => lists all of the feeds that the current logged in user is following
- go run . agg 2s => every two seconds you can fetch one rss feed links; If you want to change the time between the rss feeds are fetched or even the number of feeds fetched you can change the code in the scrapeFeeds() function => present in the rss.go file
