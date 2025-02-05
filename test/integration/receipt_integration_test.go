package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor/internal/handler"
	"receipt-processor/internal/repository"
	"receipt-processor/internal/service"
	"testing"
)

// Struct to decode JSON responses
type IDResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

func TestFullReceiptFlow(t *testing.T) {
	repo := repository.NewInMemoryReceiptStore()
	service := service.NewReceiptService(repo)
	handler := handler.NewReceiptHandler(service)

	receiptJSON := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
		],
		"total": "6.49"
	}`

	// Step 1: Process the Receipt
	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(receiptJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ProcessReceipt(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", res.StatusCode)
	}

	var idResp IDResponse
	if err := json.NewDecoder(res.Body).Decode(&idResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if idResp.ID == "" {
		t.Fatalf("Expected non-empty ID, got empty")
	}

	// Step 2: Retrieve Points for the Generated Receipt ID
	pointsReq := httptest.NewRequest("GET", "/receipts/"+idResp.ID+"/points", nil)
	w = httptest.NewRecorder()
	handler.GetPoints(w, pointsReq)

	pointsRes := w.Result()
	if pointsRes.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", pointsRes.StatusCode)
	}

	var pointsResp PointsResponse
	if err := json.NewDecoder(pointsRes.Body).Decode(&pointsResp); err != nil {
		t.Fatalf("Failed to decode points response: %v", err)
	}

	if pointsResp.Points <= 0 {
		t.Fatalf("Expected positive points, got %d", pointsResp.Points)
	}
}

func TestGetPoints_InvalidID(t *testing.T) {
	repo := repository.NewInMemoryReceiptStore()
	service := service.NewReceiptService(repo)
	handler := handler.NewReceiptHandler(service)

	req := httptest.NewRequest("GET", "/receipts/invalid-id/points", nil)
	w := httptest.NewRecorder()
	handler.GetPoints(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for invalid ID, got %d", res.StatusCode)
	}
}
