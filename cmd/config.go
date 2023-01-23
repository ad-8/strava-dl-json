package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path"
)

// All paths set here must be relative to $HOME,
// e.g. if the JSON file with the activities should be saved in /Users/alex/strava-data
// then set appDataDir = "strava-data".
//
// The appDataDir folder must contain a valid .env file (see README).
//
// The name of the JSON file can be set via the jsonFile constant.
const ( // TODO load via .env file (or use toml for everything)
	appDataDir = "strava-data"
	dotEnvFile = ".env"
	jsonFile   = "current.json"
)

// loadDotEnv reads the values that are required to generate a Strava token from an .env file.
func loadDotEnv() (clientID, clientSecret, refreshToken string) {
	dotEnvPath, err := prefixDataDir(dotEnvFile)
	if err != nil {
		log.Fatalf("error creating path of .env file: %v", err)
	}
	err = godotenv.Load(dotEnvPath)
	if err != nil {
		log.Fatal("error loading .env file")
	}
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	refreshToken = os.Getenv("REFRESH_TOKEN")
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		fmt.Println("ensure your .env file contains only the following lines:")
		fmt.Println(`
				CLIENT_ID = 123
				CLIENT_SECRET = "foo"
				REFRESH_TOKEN = "bar"
			`)
		log.Fatal("Incorrect .env file")
	}
	return
}

// checkDataDirExists makes sure that the apps data dir exists, creating it if necessary.
func checkDataDirExists() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("checkDataDirExists: %w", err)
	}
	if err := os.MkdirAll(path.Join(home, appDataDir), 0750); err != nil {
		return fmt.Errorf("checkDataDirExists: %w", err)
	}
	return nil
}

// prefixDataDir prefixes a file with the content of $HOME (e.g. /Users/alex)
// and the value set for appDataDir and returns a valid path if successful.
func prefixDataDir(file string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("prefixDataDir: %w", err)
	}
	return path.Join(home, appDataDir, file), nil
}
