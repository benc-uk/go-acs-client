package client

// ==============================================================================
// Client SDK for the Azure Communication Services SMS API for SMS
// An implementation of the API documented here:
// https://learn.microsoft.com/en-us/rest/api/communication/sms/send
// ==============================================================================

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/benc-uk/go-acs-client/auth"
	"github.com/google/uuid"
)

// NewSMS creates a new SMS message for sending
func NewSMS(from, to, msg string) *SMS {
	return &SMS{
		From: from,
		SMSRecipients: []SMSRecipient{
			{
				To:                     to,
				RepeatabilityRequestID: uuid.New().String(),
				RepeatabilityFirstSent: time.Now().UTC().Format(http.TimeFormat),
			},
		},
		Message: msg,
	}
}

// SendSingleSMS sends a single SMS and returns the API response and/or error
func (c *Client) SendSingleSMS(s *SMS) (smsResp *SMSSendResponseItem, err error) {
	postBody, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("sms failed JSON marshalling: %s", err)
	}

	bodyBuffer := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", c.Endpoint+sendSMSEndpoint+"?api-version="+c.APIVersionSMS, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("error creating API request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Sign the request using the ACS access key and HMAC-SHA256
	err = auth.SignRequestHMAC(c.AccessKey, req)
	if err != nil {
		return nil, fmt.Errorf("error signing API request: %s", err)
	}

	client := &http.Client{
		Timeout: time.Second * clientTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending API request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		// For some reason the API returns various body content on error (bad API design)
		// So we just return the raw body as a string
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("error %d sending sms: %s", resp.StatusCode, string(body))
	}

	smsRespList := &SMSSendResponse{}

	err = json.NewDecoder(resp.Body).Decode(smsRespList)
	if err != nil {
		return nil, err
	}

	return &smsRespList.Value[0], nil
}
