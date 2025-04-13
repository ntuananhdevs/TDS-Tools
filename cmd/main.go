package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"tds/utils"
)

func loadEnvVariables() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Access Token and Cookie from environment variables
	tdsToken := os.Getenv("TDS_TOKEN")
	if tdsToken == "" {
		log.Fatal("TDS_TOKEN is missing in .env file")
	}

	facebookCookie := os.Getenv("FACEBOOK_COOKIE")
	if facebookCookie == "" {
		log.Fatal("FACEBOOK_COOKIE is missing in .env file")
	}

	// Print loaded variables for validation (optional)
	fmt.Printf("TDS Token: %s\n", tdsToken)
	fmt.Printf("Facebook Cookie: %s\n", facebookCookie)

	// Use tdsToken and facebookCookie for further actions
	utils.Info("🚀 Đang khởi động tool với TDS Token và Cookie")
}

func main() {
	// Load .env file and get token and cookie
	loadEnvVariables()

	// Remaining logic to handle the application...
}
