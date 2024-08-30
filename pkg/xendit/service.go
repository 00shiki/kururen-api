package xendit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kururen/entity"
	"log"
	"net/http"
	"os"
)

type Service interface {
	CreateInvoice(*entity.RentalHistory) error
}

type XenditClient struct{}

func NewXenditService() *XenditClient {
	return &XenditClient{}
}

func (service *XenditClient) CreateInvoice(rentalHistory *entity.RentalHistory) error {
	baseURL := "https://api.xendit.co/v2"
	apiKey := os.Getenv("XENDIT_KEY")

	cars := make([]map[string]interface{}, len(rentalHistory.Cars))
	for i, car := range rentalHistory.Cars {
		cars[i] = make(map[string]interface{})
		cars[i]["name"] = fmt.Sprintf("%s's %s %s", car.Year, car.Brand, car.Model)
		cars[i]["quantity"] = rentalHistory.EndDate.Sub(rentalHistory.StartDate).Hours() / 24
		cars[i]["price"] = car.RentalCost
	}
	payload := map[string]interface{}{
		"external_id":      "1",
		"amount":           rentalHistory.Payment.Amount,
		"description":      "Car Renting Invoice",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":  rentalHistory.User.Name,
			"email": rentalHistory.User.Email,
		},
		"currency": "IDR",
		"items":    cars,
		"payment_method": []interface{}{
			rentalHistory.Payment.Type,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL+"/invoices", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(apiKey, "")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	invoice := new(entity.Invoice)
	if err = json.NewDecoder(res.Body).Decode(&invoice); err != nil {
		return err
	}
	rentalHistory.Payment.InvoiceURL = invoice.InvoiceURL
	log.Printf("invoice url: %s", invoice.InvoiceURL)

	return nil
}
