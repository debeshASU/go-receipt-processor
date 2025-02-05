package main

import (
	"log"
	"net/http"
	"receipt-processor/internal/handler"
	"receipt-processor/internal/repository"
	"receipt-processor/internal/service"

	"github.com/sirupsen/logrus"
)

// Recover from panics to prevent server crashes
func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{"error": err}).Error("Panic recovered")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	repo := repository.NewInMemoryReceiptStore()
	service := service.NewReceiptService(repo)
	handler := handler.NewReceiptHandler(service)

	http.HandleFunc("/receipts/process", recoveryMiddleware(handler.ProcessReceipt))
	http.HandleFunc("/receipts/", recoveryMiddleware(handler.GetPoints))

	logrus.Info("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
