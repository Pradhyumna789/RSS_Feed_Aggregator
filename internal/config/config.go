package config

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

/*
	SetUser Method Working:
		=> Get the file path of the file "gatorconfig.json" and initialize it to the variable "filePath" 		
		=> Read the "gatorconfig.json" file's json content
		=> Unmarshal the jsonData that is read into the config struct
		=> Change the contents of the struct - in this case SetUser method sets the CurrentUserName field 
		=> Marshal the config struct back to json - this is the updated json content
		=> Update the "gatorconfig.json" file with the updated json content
*/
func (config *Config) SetUser(current_user_name string) {
	homeDir, err := os.UserHomeDir();
	if err != nil {
		fmt.Println("Error fetching the user's directory")
		log.Fatal(err)
	}

	filePath := filepath.Join(homeDir, "gatorconfig.json")

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file gatorconfig.json")
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		fmt.Println("Error in unmarshalling json to a config struct")
		log.Fatal(err)
	}

	config.CurrentUserName = current_user_name

	updatedJson, err := json.Marshal(config)

	err = os.WriteFile(filePath, updatedJson, 0644)
	if err != nil {
		fmt.Println("Error writing to the file gatorconfig.json")
		log.Fatal(err)
	}

}

// The Read() function takes the gatorconfig.json file and reads from it
func Read() Config { 
	homeDir, err := os.UserHomeDir() 
	if err != nil {
		fmt.Println("Error fetching the user's directory")
		log.Fatal(err)
	}

	filePath := filepath.Join(homeDir, "gatorconfig.json")

	jsonData, err := os.ReadFile(filePath) 
	if err != nil {
		fmt.Println("Error fetching data from gatorconfig.json file")
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		fmt.Println("Error converting json into bytes")
		log.Fatal(err)
	}

	return config

}

