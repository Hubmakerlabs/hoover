package main

import (
	"encoding/json"
	"fmt"
	"os"

	farcaster "github.com/ertan/go-farcaster/pkg"
	"github.com/spf13/viper"
)

func writeToFile(st interface{}, filename string) {
	// Marshal the cast data to JSON
	stJson, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		panic(err)
	}

	// Open the file, or create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(stJson)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Casts written to %s\n", filename)
}

func main() {
	// Load the environment variables
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiUrl := viper.Get("FARCASTER_API_URL").(string)
	mnemonic := viper.Get("FARCASTER_MNEMONIC").(string)
	providerWs := viper.Get("ETHEREUM_PROVIDER_WS").(string)

	// Create a Farcaster client
	farcaster := farcaster.NewFarcasterClient(apiUrl, mnemonic, providerWs)
	fmt.Println("Farcaster client created")

	// Fetch recent casts
	casts, _, err := farcaster.Casts.GetRecentCasts(10)
	if err != nil {
		panic(err)
	}
	fmt.Println("Recent casts fetched")

	// Write the fetched casts to a JSON file
	writeToFile(casts, "recent_casts.json")
}
