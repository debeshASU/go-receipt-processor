package repository

import (
	"receipt-processor/internal/model"
	"sync"
)

type ReceiptStore interface {
	SaveReceipt(id string, receipt model.Receipt)
	GetReceipt(id string) (model.Receipt, bool)
}

type InMemoryReceiptStore struct {
	mu       sync.RWMutex
	receipts map[string]model.Receipt
}

func NewInMemoryReceiptStore() *InMemoryReceiptStore {
	return &InMemoryReceiptStore{
		receipts: make(map[string]model.Receipt),
	}
}

func (s *InMemoryReceiptStore) SaveReceipt(id string, receipt model.Receipt) {
	s.mu.Lock()
	s.receipts[id] = receipt
	s.mu.Unlock()
}

func (s *InMemoryReceiptStore) GetReceipt(id string) (model.Receipt, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	receipt, exists := s.receipts[id]
	return receipt, exists
}
