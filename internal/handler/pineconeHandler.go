package handler

import (
	"github.com/pinecone-io/go-pinecone/pinecone"
)

type PineconeHandler struct {
	pinecone *pinecone.IndexConnection
}

// NewPineconeHandler는 PineconeHandler를 초기화하는 함수입니다.
func NewPineconeHandler(pcIndex *pinecone.IndexConnection) (*PineconeHandler, error) {
	pcHandler := &PineconeHandler{pinecone: pcIndex}
	return pcHandler, nil
}
