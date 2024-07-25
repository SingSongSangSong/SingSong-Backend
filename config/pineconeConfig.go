package config

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type PineconeConfig struct {
	ApiKey     string            // required - provide through NewClientParams or environment variable PINECONE_API_KEY
	Headers    map[string]string // optional
	Host       string            // optional
	RestClient *http.Client      // optional
	SourceTag  string            // optional
}

var (
	PConf *PineconeConfig
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("PINECONE_API_KEY")
	if apiKey == "" {
		log.Fatalf("Pincone API key is required")
	}

	PConf = &PineconeConfig{
		ApiKey: apiKey,
	}

	log.Printf("init pinecone config success")
}
