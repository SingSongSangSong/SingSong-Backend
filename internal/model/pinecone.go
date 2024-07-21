package model

import (
	"context"
	"fmt"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"http-practice/config"
	"log"
	"os"
)

// NewPineconeClient는 Pinecone 클라이언트를 초기화하는 함수입니다.
func NewPineconeClient(ctx context.Context, config *config.PineconeConfig) (*pinecone.IndexConnection, error) {
	if config.ApiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// 실제 클라이언트 초기화 코드
	client, err := pinecone.NewClient(pinecone.NewClientParams{ApiKey: config.ApiKey})
	if err != nil {
		return nil, err
	}

	// context를 사용하여 DescribeIndex 호출
	idx, err := client.DescribeIndex(ctx, os.Getenv("PINECONE_INDEX"))
	if err != nil {
		log.Fatalf("Failed to describe index \"%s\". Error:%s", idx.Name, err)
	} else {
		fmt.Printf("Successfully found the \"%s\" index!\n", idx.Name)
	}

	idxConnection, err := client.Index(idx.Host)
	if err != nil {
		log.Fatalf("Failed to create IndexConnection for Host: %v. Error: %v", idx.Host, err)
	} else {
		log.Println("IndexConnection created successfully!")
	}

	return idxConnection, nil
}
