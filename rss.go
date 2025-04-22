package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"html"
	"github.com/Pradhyumna789/RSS/internal/config"
	"github.com/Pradhyumna789/RSS/internal/database"
	"github.com/google/uuid"
	_ "github.com/cweill/gotests"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandSystem map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
} 

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed ,error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error in creating a request to the url: %w", err)
	}

	client := &http.Client{}

	req.Header.Add("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error in getting a response: %w", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error converting the response's body into bytes of data: %w", err)
	}
	
	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error in unmarshlling the xml data into a go struct: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Link = html.UnescapeString(rssFeed.Channel.Link)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Link = html.UnescapeString(rssFeed.Channel.Item[i].Link)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
		rssFeed.Channel.Item[i].PubDate = html.UnescapeString(rssFeed.Channel.Item[i].PubDate)
	}

	return &rssFeed, nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("To fetch feeds continuously mention the time between requests along with the agg command: %w")
	} 

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error in parsing the duration string into time.Duration value: %w", err)
	}
	
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for range ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf("error scraping feeds: %w", err)
		}
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("enter the name and url of the feed to add a feed")
	}

	ctx := context.Background()
	name := cmd.args[0]	
	url := cmd.args[1]

	// Create the feed
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Createdat: time.Now(),
		Updatedat: time.Now(),
		FeedName:  name,
		FeedUrl:   url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("error in adding the feed to the database: %w", err)
	}

	// Create the feed follow relationship
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		Createdat: time.Now(),
		Updatedat: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return fmt.Errorf("error in creating a feed follow record: %w", err)
	}

	fmt.Println("Feed added successfully!")
	fmt.Println("Feed ID:", feed.ID)
	fmt.Println("Feed Name:", feed.FeedName)
	fmt.Println("Feed URL:", feed.FeedUrl)
	fmt.Println("Created at:", feed.Createdat)
	fmt.Println("Updated at:", feed.Updatedat)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("please enter the login command with the user name")
	}

	userName := cmd.args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, userName)
	if err != nil {
		return fmt.Errorf("user with user-name: %s doesn't exist please register first to login in, error logging in: %w", userName, err)
	}

	s.config.CurrentUserName = user.UserName 

	jsonData, err := json.Marshal(s.config)
	if err != nil {
		return fmt.Errorf("error in marshalling json to a config struct: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error fetching the user's home directory")
	}

	filePath := filepath.Join(homeDir, ".gatorconfig.json")

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error in writing and updating the gatorconfig.json file: %w", err)
	}

	fmt.Printf("user %s is been logged in ", user.UserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Please enter the register command with the user name")
	}

	q := s.db
	id := uuid.New()
	userName := cmd.args[0]

	ctx := context.Background()	
	parameters := database.CreateUserParams{
		ID: id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserName: userName,
	}

	user, err := q.CreateUser(ctx, parameters)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_user_name_key\"" {
			fmt.Printf("User with name '%s' already exists\n", userName)
			os.Exit(1) // Exit with code 1 as required
    	}
		return fmt.Errorf("error in creating the user: %w", err)
	}

	s.config.CurrentUserName = user.UserName

	jsonData, err := json.Marshal(s.config)
	if err != nil {
		return fmt.Errorf("error in marshalling json to a config struct: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error fetching the user's home directory")
	}

	filePath := filepath.Join(homeDir, ".gatorconfig.json")

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error in writing and updating the gatorconfig.json file: %w", err)
	}

	fmt.Println("user is created")	
	fmt.Println("user id: ", user.ID)
	fmt.Println("user is created at time: ", user.CreatedAt)
	fmt.Println("user is updated at time: ", user.UpdatedAt)
	fmt.Println("user's name is: ", user.UserName)

	return nil 
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	
	// We need to modify the DeleteUser function to handle the foreign key constraints
	// For now, we'll use a direct SQL query with the database connection from main
	dbURL := "postgres://postgres:Omvwsuv200319@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("error in opening connection to the database: ", err)
		os.Exit(1)
	}
	defer db.Close()
	
	// First, delete all feeds
	_, err = db.ExecContext(ctx, "DELETE FROM feeds")
	if err != nil {
		fmt.Println("error in deleting all feeds: ", err)
		os.Exit(1)
	}
	
	// Then delete all users
	_, err = db.ExecContext(ctx, "DELETE FROM users")
	if err != nil {
		fmt.Println("error in deleting all users: ", err)
		os.Exit(1)
	}

	fmt.Println("successfully deleted all users and feeds and reset the tables")
	os.Exit(0)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	query := s.db
	users, err := query.GetUsers(ctx)
	if err != nil {
		fmt.Println("error in fetching all the users from the database: ", err)
		os.Exit(1)
	}

	cf := s.config
	for _, user := range users {
		if cf.CurrentUserName == user.UserName {
			fmt.Println("*", user.UserName, "(current)")
		} 
		fmt.Println("* ", user.UserName)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		fmt.Println("error in fetching all the feeds from the database: ", err)
		os.Exit(1)
	}
	
	// printing the feeds and the user who owns the feed
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		fmt.Println("error in fetching all the users from the database: ", err)
		os.Exit(1)
	}
	
	userMap := make(map[uuid.UUID]string)
	for _, user := range users {
		userMap[user.ID] = user.UserName
	}

	for _, feed := range feeds {
		fmt.Println("Feed ID:", feed.ID)
		fmt.Println("Feed Name:", feed.FeedName)
		fmt.Println("Feed URL:", feed.FeedUrl)
		fmt.Println("Created at:", feed.Createdat)
		fmt.Println("Updated at:", feed.Updatedat)
		fmt.Println("User:", userMap[feed.UserID]) 		
		fmt.Println("--------------------------------")
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("please enter the follow command along with the feed url to follow the feed")
	}

	feed_url := cmd.args[0]
	ctx := context.Background()
	id := uuid.New()

	feedId, err := s.db.GetFeedByURL(ctx, feed_url)
	if err != nil {
		return fmt.Errorf("error in fetching feed's id: %w", err)
	}

	parameters := database.CreateFeedFollowParams{
		ID: id,
		Createdat: time.Now(),
		Updatedat: time.Now(),
		UserID: user.ID,
		FeedID: feedId,
	}

	feeds, err := s.db.CreateFeedFollow(ctx, parameters)
	if err != nil {
		return fmt.Errorf("error in creating feed follow: %w", err)
	}

	for _, feed := range feeds {
		fmt.Println("Feed Name:", feed.FeedName)
		fmt.Println("User Name:", feed.UserName)
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feed_follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error in fetching feed follows for the user: %w", err)
	}

	for _, val := range feed_follows {
		fmt.Println(val.FeedName)
	}

	return nil
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		userName := s.config.CurrentUserName

		user, err := s.db.GetUser(ctx, userName)
		if err != nil {
			return fmt.Errorf("user doesn't exist: %w", err)
		}


		return handler(s, cmd, user)
	}
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("enter the unfollow command along with the url of the feed that you want to unfollow")
	}

	feedUrl := cmd.args[0]
	ctx := context.Background()

	// First check if the user is actually following this feed
	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error checking feed follows: %w", err)
	}

	// Get the feed ID for the URL
	feedId, err := s.db.GetFeedByURL(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("feed with URL %s not found: %w", feedUrl, err)
	}

	// Check if user is following this feed
	isFollowing := false
	for _, follow := range feedFollows {
		feedName, err := s.db.GetFeedNameById(ctx, feedId)
		if err != nil {
			return fmt.Errorf("error getting feed name: %w", err)
		}
		if follow.FeedName == feedName {
			isFollowing = true
			break
		}
	}

	if !isFollowing {
		return fmt.Errorf("you are not following the feed with URL %s", feedUrl)
	}

	// If we get here, the user is following the feed, so we can unfollow
	params := database.DeleteFeedFollowRecordParams{
		UserID:  user.ID,
		FeedUrl: feedUrl,
	}

	err = s.db.DeleteFeedFollowRecord(ctx, params)
	if err != nil {
		return fmt.Errorf("error unfollowing the feed: %w", err)
	}

	fmt.Printf("Successfully unfollowed feed with URL: %s\n", feedUrl)
	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeeds, err := s.db.GetNextFeedToFetch(ctx)	
	if err != nil {
		return fmt.Errorf("error in fetching the next feed: %w", err)
	}

	for _, feed := range nextFeeds {
		fmt.Printf("\nProcessing feed: %s\n", feed.FeedName)
		rssFeed, err := fetchFeed(ctx, feed.FeedUrl)
		if err != nil {
			fmt.Println("Error fetching feed because it's not of an xml type")
			continue  // Skip this feed and continue with the next one
		}

		itemsToShow := 1 
		if len(rssFeed.Channel.Item) > itemsToShow {
			rssFeed.Channel.Item = rssFeed.Channel.Item[:itemsToShow]
		}

		for _, val := range rssFeed.Channel.Item {
			fmt.Println("Title:", val.Title)
			fmt.Println("Link:", val.Link)
			fmt.Println("----------------------------------------")
		}
	}

	err = s.db.MarkFeedFetched(ctx)
	if err != nil {
		return fmt.Errorf("error in marking the feed as fetched: %w", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandSystem[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, found := c.commandSystem[cmd.name]
	if !found {
		return fmt.Errorf("command not found: %s", cmd.name)
	}

	return handler(s, cmd)
}

func main() {
	dbURL := "postgres://postgres:Omvwsuv200319@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error in opening connection to the database: ", err)
	}

	q := database.New(db)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error in fetching the user's home directory ", err)
	}

	filepath := filepath.Join(homeDir, ".gatorconfig.json")
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("error in reading the file's content ", err)
	}

	var c config.Config
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		log.Fatal("error unmarshalling json into a config struct ", err)
	}

	s := state{
		db: q,
		config: &c,
	}

	commands := commands{
		commandSystem: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("agg", handlerAgg)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("enter the command name")
	}

	cmdName := args[1]
	cmdArg := args[2:]

	command := command{
		name: cmdName,
		args: cmdArg,
	}

	err = commands.run(&s, command)
	if err != nil {
		log.Fatal("error running the command ", err)
	}
}
