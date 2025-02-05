package service

import (
	"errors"
	"math"
	"receipt-processor/internal/model"
	"receipt-processor/internal/repository"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ReceiptService struct {
	repo repository.ReceiptStore
}

func NewReceiptService(repo repository.ReceiptStore) *ReceiptService {
	return &ReceiptService{repo: repo}
}

func (s *ReceiptService) ProcessReceipt(receipt model.Receipt) (string, error) {
	if err := s.ValidateReceipt(receipt); err != nil {
		logrus.WithError(err).Error("Validation failed for receipt")
		return "", err
	}

	id := uuid.New().String()
	s.repo.SaveReceipt(id, receipt)

	logrus.WithFields(logrus.Fields{
		"receipt_id": id,
		"retailer":   receipt.Retailer,
		"items":      len(receipt.Items),
	}).Info("Receipt processed successfully")

	return id, nil
}

func (s *ReceiptService) GetReceiptPoints(id string) (int, error) {
	receipt, exists := s.repo.GetReceipt(id)
	if !exists {
		err := errors.New("no receipt found for that ID")
		logrus.WithFields(logrus.Fields{"receipt_id": id}).Error(err)
		return 0, err
	}

	points := s.CalculatePoints(receipt)

	logrus.WithFields(logrus.Fields{
		"receipt_id": id,
		"points":     points,
	}).Info("Points calculated successfully")

	return points, nil
}

func (s *ReceiptService) ValidateReceipt(receipt model.Receipt) error {
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
		return errors.New("missing required fields")
	}

	if matched, _ := regexp.MatchString(`^[\w\s\-&]+$`, receipt.Retailer); !matched {
		return errors.New("invalid retailer name")
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return errors.New("invalid purchase date format (expected YYYY-MM-DD)")
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return errors.New("invalid purchase time format (expected HH:MM)")
	}

	for _, item := range receipt.Items {
		if matched, _ := regexp.MatchString(`^[\w\s\-]+$`, item.ShortDescription); !matched {
			return errors.New("invalid item description")
		}
		if matched, _ := regexp.MatchString(`^\d+\.\d{2}$`, item.Price); !matched {
			return errors.New("invalid item price format (expected XX.XX)")
		}
	}

	if matched, _ := regexp.MatchString(`^\d+\.\d{2}$`, receipt.Total); !matched {
		return errors.New("invalid total amount format (expected XX.XX)")
	}

	return nil
}

func (s *ReceiptService) CalculatePoints(receipt model.Receipt) int {
	points := 0

	points += len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(receipt.Retailer, -1))
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && total == float64(int(total)) {
		points += 50
	}

	if err == nil && math.Mod(total*100, 25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		price, err := strconv.ParseFloat(item.Price, 64)
		if err == nil && trimmedLength%3 == 0 {
			points += int(math.Ceil(price * 0.2))
		}
	}

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		points += 6
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}
