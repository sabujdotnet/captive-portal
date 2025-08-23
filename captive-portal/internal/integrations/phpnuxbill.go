package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type PHPNuxBillIntegration struct {
	baseURL string
	apiKey  string
}

func NewPHPNuxBillIntegration(cfg config.PHPNuxBillConfig) *PHPNuxBillIntegration {
	return &PHPNuxBillIntegration{
		baseURL: cfg.URL,
		apiKey:  cfg.APIKey,
	}
}

func (p *PHPNuxBillIntegration) ValidateVoucher(code string) (bool, map[string]interface{}, error) {
	formData := url.Values{}
	formData.Set("voucher", code)
	formData.Set("api_key", p.apiKey)
	
	resp, err := http.PostForm(p.baseURL+"/api/voucher/validate", formData)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, nil, err
	}
	
	valid, ok := result["valid"].(bool)
	if !ok {
		return false, nil, fmt.Errorf("invalid response format")
	}
	
	return valid, result, nil
}

func (p *PHPNuxBillIntegration) CreateCustomer(name, email, phone string) (string, error) {
	formData := url.Values{}
	formData.Set("name", name)
	formData.Set("email", email)
	formData.Set("phone", phone)
	formData.Set("api_key", p.apiKey)
	
	resp, err := http.PostForm(p.baseURL+"/api/customer/create", formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	customerID, ok := result["customer_id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to create customer")
	}
	
	return customerID, nil
}

func (p *PHPNuxBillIntegration) AddPayment(customerID string, amount float64, method string) error {
	formData := url.Values{}
	formData.Set("customer_id", customerID)
	formData.Set("amount", fmt.Sprintf("%.2f", amount))
	formData.Set("method", method)
	formData.Set("api_key", p.apiKey)
	
	resp, err := http.PostForm(p.baseURL+"/api/payment/add", formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	
	status, ok := result["status"].(string)
	if !ok || status != "success" {
		return fmt.Errorf("failed to add payment")
	}
	
	return nil
}