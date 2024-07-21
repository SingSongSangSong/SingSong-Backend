package config

import "net/http"

type PineconeConfig struct {
	ApiKey     string            // required - provide through NewClientParams or environment variable PINECONE_API_KEY
	Headers    map[string]string // optional
	Host       string            // optional
	RestClient *http.Client      // optional
	SourceTag  string            // optional
}

func NewPineconeConfig(ApiKey string) *PineconeConfig {
	return &PineconeConfig{ApiKey: ApiKey}
}
