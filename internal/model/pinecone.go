package model

import (
	"SingSong-Backend/config"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"math/rand"
	"os"
)

type PineconeModel struct {
	idxConnection *pinecone.IndexConnection
}

// NewPineconeClient는 Pinecone 클라이언트를 초기화하는 함수입니다.
func NewPineconeClient(ctx context.Context, config *config.PineconeConfig) (*PineconeModel, error) {
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

	return &PineconeModel{idxConnection: idxConnection}, nil
}

func (pineconeModel *PineconeModel) QueryPineconeWithTag(filterStruct *structpb.Struct) (*pinecone.QueryVectorsResponse, error) {
	// Define a dummy vector (e.g., zero vector) for the query
	dummyVector := make([]float32, 30) // Assuming the vector length is 1536, adjust as necessary
	for i := range dummyVector {
		dummyVector[i] = rand.Float32() //random vector
	}

	// 쿼리 요청을 보냅니다.
	values, err := pineconeModel.idxConnection.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          dummyVector,
		TopK:            100,
		Filter:          filterStruct,
		SparseValues:    nil,
		IncludeValues:   true,
		IncludeMetadata: true,
	})
	return values, err
}

func (pineconeModel *PineconeModel) FetchVectors(c *gin.Context, songs []string) (*pinecone.FetchVectorsResponse, error) {
	vectors, err := pineconeModel.idxConnection.FetchVectors(c, songs)
	return vectors, err
}

func (pineconeModel *PineconeModel) QueryByVectorValues(c *gin.Context, p *pinecone.QueryByVectorValuesRequest) (*pinecone.QueryVectorsResponse, error) {
	values, err := pineconeModel.idxConnection.QueryByVectorValues(c, p)
	return values, err
}
