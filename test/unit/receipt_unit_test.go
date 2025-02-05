package unit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"receipt-processor/internal/handler"
	"receipt-processor/internal/repository"
	"receipt-processor/internal/service"
	"testing"
)

func TestProcessReceipt_Valid(t *testing.T) {
	repo := repository.NewInMemoryReceiptStore()
	service := service.NewReceiptService(repo)
	handler := handler.NewReceiptHandler(service)

	validJSON := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
		],
		"total": "6.49"
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(validJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ProcessReceipt(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", res.StatusCode)
	}
}

func TestProcessReceipt_MissingField(t *testing.T) {
	repo := repository.NewInMemoryReceiptStore()
	service := service.NewReceiptService(repo)
	handler := handler.NewReceiptHandler(service)

	invalidJSON := `{
		"retailer": "",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [],
		"total": "6.49"
	}`

	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ProcessReceipt(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422, got %d", res.StatusCode)
	}
}
