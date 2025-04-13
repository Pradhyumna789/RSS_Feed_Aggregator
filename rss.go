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

	html.UnescapeString(rssFeed.Channel.Title)
	html.UnescapeString(rssFeed.Channel.Link)
	html.UnescapeString(rssFeed.Channel.Description)

	for _, val := range rssFeed.Channel.Item {
		html.UnescapeString(val.Title)
		html.UnescapeString(val.Link)
		html.UnescapeString(val.Description)
		html.UnescapeString(val.PubDate)
	}

	return &rssFeed, nil

}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()

	rssFeed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error in fetching the feed: %w", err)
	}

	fmt.Println(rssFeed)	

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("enter the name and url of the feed to add a feed")
	}

	uuid := uuid.New()
	ctx := context.Background()
	name := cmd.args[0]	
	url := cmd.args[1]

	userId, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		fmt.Println("error in fetching the user", err)
		os.Exit(1)
	}

	feedParams := database.CreateFeedParams{
		ID: uuid,
		Createdat: time.Now(),
		Updatedat: time.Now(),
		FeedName: sql.NullString{String: name, Valid: true},
		FeedUrl: sql.NullString{String: url, Valid: true},
		UserID: sql.NullString{String: userId, Valid: true},
	}

	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		fmt.Println("error in adding the feed to the database", err)
		os.Exit(1)
	}

	fmt.Println(feed.ID)
	fmt.Println(feed.Createdat)
	fmt.Println(feed.Updatedat)
	fmt.Println(feed.FeedName)
	fmt.Println(feed.FeedUrl)
	fmt.Println(feed.UserID)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("please enter the login command with the user name")
	}

	userName := cmd.args[0]
	ctx := context.Background()

	queriedUserName, err := s.db.GetUser(ctx, userName)
	if err != nil {
		return fmt.Errorf("user with user-name: %s doesn't exist please register first to login in, error logging in: %w", userName, err)
	}

	s.config.CurrentUserName = queriedUserName 

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

	fmt.Printf("user %s is been logged in ", queriedUserName)

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
	err := s.db.DeleteUser(ctx) 
	if err != nil {
		fmt.Println("error in deleting all the users: ", err)
		os.Exit(1)
	}

	fmt.Println("successful in deleting all the users and resetting the table")
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
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)

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
