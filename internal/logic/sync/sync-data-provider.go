package sync

import "github.com/dimazhornyk/generic-proving-network/internal/logic"

type DataProvider struct {
	storage *logic.Storage
}

func NewDataProvider() *DataProvider {
	return &DataProvider{}
}
