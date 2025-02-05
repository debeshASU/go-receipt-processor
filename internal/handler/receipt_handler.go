package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"receipt-processor/internal/model"
	"receipt-processor/internal/service"

	"github.com/sirupsen/logrus"
)

type ReceiptHandler struct {
	service *service.ReceiptService
}

func NewReceiptHandler(service *service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{service: service}
}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var receipt model.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	id, err := h.service.ProcessReceipt(receipt)
	if err != nil {
		// Change this from 400 to 422 for validation failures
		if strings.Contains(err.Error(), "missing required fields") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	logrus.WithFields(logrus.Fields{"receipt_id": id}).Info("Receipt processed successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.IDResponse{ID: id})
}

func (h *ReceiptHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.Split(id, "/")[0]

	if id == "" {
		http.Error(w, "Missing receipt ID", http.StatusBadRequest)
		return
	}

	points, err := h.service.GetReceiptPoints(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"receipt_id": id, "status": "not found"}).Error("Receipt lookup failed")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	logrus.WithFields(logrus.Fields{"receipt_id": id, "points": points}).Info("Points retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.PointsResponse{Points: points})
}
