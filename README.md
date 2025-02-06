Receipt Processor API

Overview

The Receipt Processor API is a web service that processes receipts and calculates reward points based on predefined rules. The API supports running as a standalone Go application or within a Docker container.

Features

Process Receipts: Submit a receipt and receive a unique ID.

Retrieve Points: Query the points awarded for a specific receipt ID.

In-Memory Storage: The system currently stores data in memory but is designed for future scalability.

Scalable Architecture: Built with modular components to easily swap storage methods.

Installation & Running the Application

Running with Go

Prerequisites:

Go 1.23.5 or later must be installed.

Dependencies:

github.com/google/uuid v1.6.0

github.com/sirupsen/logrus v1.9.3

golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8

Steps to Run:

# Clone the repository
https://github.com/debeshASU/go-receipt-processor.git
cd receipt-processor

# Install dependencies
go mod tidy

# Run the application
go run cmd/receipt-processor/main.go

The server will start on http://localhost:8080.

Running with Docker

Prerequisites:

Docker must be installed and running.

Steps to Build and Run the Docker Container:

# Build the Docker image
docker build --platform=linux/amd64 -t go-receipt-processor .

# Run the container
docker run -p 8080:8080 go-receipt-processor

API Endpoints

1. Process a Receipt

Endpoint: POST /receipts/process

Request Body:

{
   "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}

Response:

{
  "id": "72e64d15-a60c-41e4-afa3-32d9085eff6f"
}

2. Get Points for a Receipt

Endpoint: GET /receipts/{id}/points

Response:

{
  "points": 28
}

Running Tests

To ensure the application works correctly, run unit and integration tests.

# Run unit and integration tests
go test ./test/unit ./test/integration -v

Logging

The application uses logrus for structured logging.

Logs include structured fields such as receipt_id, status, and points.

Future Improvements

Database Integration (PostgreSQL, MongoDB, etc.)

Authentication and Authorization

Enhanced Error Handling

API Rate Limiting

License

This project is licensed under the MIT License.

