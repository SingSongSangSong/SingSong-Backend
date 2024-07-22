package handler

import (
	"github.com/pinecone-io/go-pinecone/pinecone"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
	"math/rand"
)

func (pineconeHandler *PineconeHandler) queryPineconeWithTag(filterStruct *structpb.Struct) (*pinecone.QueryVectorsResponse, error) {
	// Define a dummy vector (e.g., zero vector) for the query
	dummyVector := make([]float32, 30) // Assuming the vector length is 1536, adjust as necessary
	for i := range dummyVector {
		dummyVector[i] = rand.Float32() //random vector
	}

	// 쿼리 요청을 보냅니다.
	values, err := pineconeHandler.pinecone.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          dummyVector,
		TopK:            100,
		Filter:          filterStruct,
		SparseValues:    nil,
		IncludeValues:   true,
		IncludeMetadata: true,
	})
	return values, err
}
