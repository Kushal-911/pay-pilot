package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL, clientID, secret string
	http                      *http.Client
}

type OrderResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Links  []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

type CaptureResult struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PurchaseUnits []struct {
		Payments struct {
			Captures []struct {
				ID     string `json:"id"`
				Status string `json:"status"`
			} `json:"captures"`
		} `json:"payments"`
	} `json:"purchase_units"`
}

type WebhookHeaders struct {
	AuthAlgo         string
	CertURL          string
	TransmissionID   string
	TransmissionSig  string
	TransmissionTime string
}

type webhookVerificationRequest struct {
	AuthAlgo         string          `json:"auth_algo"`
	CertURL          string          `json:"cert_url"`
	TransmissionID   string          `json:"transmission_id"`
	TransmissionSig  string          `json:"transmission_sig"`
	TransmissionTime string          `json:"transmission_time"`
	WebhookID        string          `json:"webhook_id"`
	Event            json.RawMessage `json:"-"`
}

func (r webhookVerificationRequest) MarshalJSON() ([]byte, error) {
	if !json.Valid(r.Event) {
		return nil, fmt.Errorf("invalid webhook event JSON")
	}
	headers, err := json.Marshal(map[string]string{"auth_algo": r.AuthAlgo, "cert_url": r.CertURL, "transmission_id": r.TransmissionID, "transmission_sig": r.TransmissionSig, "transmission_time": r.TransmissionTime, "webhook_id": r.WebhookID})
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, len(headers)+len(r.Event)+20)
	result = append(result, headers[:len(headers)-1]...)
	result = append(result, `,"webhook_event":`...)
	result = append(result, r.Event...)
	result = append(result, '}')
	return result, nil
}

func (c *Client) VerifyWebhook(ctx context.Context, webhookID string, headers WebhookHeaders, event json.RawMessage) (bool, error) {
	var result struct {
		VerificationStatus string `json:"verification_status"`
	}
	body := webhookVerificationRequest{AuthAlgo: headers.AuthAlgo, CertURL: headers.CertURL, TransmissionID: headers.TransmissionID, TransmissionSig: headers.TransmissionSig, TransmissionTime: headers.TransmissionTime, WebhookID: webhookID, Event: event}
	if err := c.do(ctx, http.MethodPost, "/v1/notifications/verify-webhook-signature", body, &result, ""); err != nil {
		return false, err
	}
	return result.VerificationStatus == "SUCCESS", nil
}

func New(baseURL, clientID, secret string, client *http.Client) *Client {
	if baseURL == "" {
		baseURL = "https://api-m.sandbox.paypal.com"
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{strings.TrimRight(baseURL, "/"), clientID, secret, client}
}

func (c *Client) CreateOrder(ctx context.Context, reference, description, currency string, minor int64, returnURL, cancelURL, requestID string) (OrderResult, error) {
	var out OrderResult
	amount := fmt.Sprintf("%d.%02d", minor/100, minor%100)
	body := map[string]any{"intent": "CAPTURE", "purchase_units": []any{map[string]any{"reference_id": reference, "description": description, "amount": map[string]string{"currency_code": currency, "value": amount}}}, "application_context": map[string]string{"return_url": returnURL, "cancel_url": cancelURL, "user_action": "PAY_NOW"}}
	err := c.do(ctx, http.MethodPost, "/v2/checkout/orders", body, &out, requestID)
	return out, err
}

func (c *Client) Capture(ctx context.Context, orderID, requestID string) (CaptureResult, error) {
	var out CaptureResult
	err := c.do(ctx, http.MethodPost, "/v2/checkout/orders/"+url.PathEscape(orderID)+"/capture", map[string]string{}, &out, requestID)
	return out, err
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (CaptureResult, error) {
	var out CaptureResult
	err := c.do(ctx, http.MethodGet, "/v2/checkout/orders/"+url.PathEscape(orderID), nil, &out, "")
	return out, err
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any, requestID string) error {
	if c.clientID == "" || c.secret == "" {
		return fmt.Errorf("PayPal credentials are not configured")
	}
	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return err
	}
	tokenReq.SetBasicAuth(c.clientID, c.secret)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenResp, err := c.http.Do(tokenReq)
	if err != nil {
		return fmt.Errorf("PayPal OAuth request: %w", err)
	}
	defer tokenResp.Body.Close()
	if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
		return fmt.Errorf("PayPal OAuth returned HTTP %d", tokenResp.StatusCode)
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&token); err != nil {
		return fmt.Errorf("decode PayPal token: %w", err)
	}
	var data []byte
	if body != nil {
		if webhookBody, ok := body.(webhookVerificationRequest); ok {
			data, err = webhookBody.MarshalJSON()
		} else {
			data, err = json.Marshal(body)
		}
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("PayPal API request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PayPal API returned HTTP %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode PayPal response: %w", err)
	}
	return nil
}
