package client

// ==============================================================================
// Client SDK for the Azure Communication Services Email API
// An implementation of the API documented here:
// https://learn.microsoft.com/en-us/rest/api/communication/email/send
// ==============================================================================

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/benc-uk/go-acs-client/auth"

	"github.com/google/uuid"
)

// SendEmail sends an email and returns the message ID and any error
func (c *Client) SendEmail(e *Email) (messageID string, err error) {
	postBody, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("email failed JSON marshalling: %s", err)
	}

	bodyBuffer := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", c.Endpoint+sendEmailEndpoint+"?api-version="+c.APIVersionEmail, bodyBuffer)
	if err != nil {
		return "", fmt.Errorf("error creating API request: %s", err)
	}

	// Sign the request using the ACS access key and HMAC-SHA256
	err = auth.SignRequestHMAC(c.AccessKey, req)
	if err != nil {
		return "", fmt.Errorf("error signing API request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Important, without these headers the request will fail
	req.Header.Set("repeatability-request-id", uuid.New().String())
	req.Header.Set("repeatability-first-sent", time.Now().UTC().Format(http.TimeFormat))

	client := &http.Client{
		Timeout: time.Second * clientTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending API request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		commError := ErrorResponse{}

		err = json.NewDecoder(resp.Body).Decode(&commError)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("error sending email: %s", commError.Error.Message)
	}

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("error sending email: status: %d", resp.StatusCode)
	}

	// This header seems to be the message ID
	messageID = resp.Header.Get("x-ms-request-id")

	return messageID, nil
}

// GetStatus gets the status of an email message sent using SendEmail()
func (c *Client) GetEmailStatus(messageID string) (status string, err error) {
	req, err := http.NewRequest("GET", c.Endpoint+fmt.Sprintf(statusEmailEndpoint, messageID)+"?api-version="+c.APIVersionEmail, nil)
	if err != nil {
		return "", err
	}

	// Sign the request using the ACS access key and HMAC-SHA256
	err = auth.SignRequestHMAC(c.AccessKey, req)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: time.Second * clientTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		commError := ErrorResponse{}

		err = json.NewDecoder(resp.Body).Decode(&commError)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		return "", fmt.Errorf("error getting status: %s", commError.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting status: status: %d", resp.StatusCode)
	}

	statusResult := &SendStatusResult{}

	err = json.NewDecoder(resp.Body).Decode(statusResult)
	if err != nil {
		return "", err
	}

	return statusResult.Status, nil
}

func newEmail(from, to, subject string) *Email {
	return &Email{
		Recipients: Recipients{
			To: []Address{
				{
					DisplayName: to,
					Email:       to,
				},
			},
		},
		Sender: from,
		Content: Content{
			Subject: subject,
		},
		Importance: ImportanceNormal,
	}
}

// NewHTMLEmail creates a new email with HTML content
func NewHTMLEmail(from, to, subject, body string) *Email {
	e := newEmail(from, to, subject)
	e.Content.HTML = body

	return e
}

// NewPlainEmail creates a new email with plain text content
func NewPlainEmail(from, to, subject, body string) *Email {
	e := newEmail(from, to, subject)
	e.Content.PlainText = body

	return e
}

// AddCC adds a CC recipient
func (e *Email) AddCC(address, displayName string) {
	e.Recipients.CC = append(e.Recipients.CC, Address{
		DisplayName: displayName,
		Email:       address,
	})
}

// AddBCC adds a BCC recipient
func (e *Email) AddBCC(address, displayName string) {
	e.Recipients.BCC = append(e.Recipients.BCC, Address{
		DisplayName: displayName,
		Email:       address,
	})
}

// AddReplyTo adds a reply to address
func (e *Email) AddReplyTo(address, displayName string) {
	e.ReplyTo = append(e.ReplyTo, Address{
		DisplayName: displayName,
		Email:       address,
	})
}

// AddCustomHeader adds a custom header to the email
func (e *Email) AddCustomHeader(key, value string) {
	e.Headers = append(e.Headers, CustomHeader{
		Name:  key,
		Value: value,
	})
}

// EnableUserEngagementTracking enables user engagement tracking
func (e *Email) EnableUserEngagementTracking() {
	e.Tracking = false
}

// DisableUserEngagementTracking disables user engagement tracking
func (e *Email) DisableUserEngagementTracking() {
	e.Tracking = true
}

// AddAttachmentRaw adds an attachment to the email as raw bytes
func (e *Email) AddAttachmentRaw(name string, content []byte, attachmentType string) {
	b64content := base64.StdEncoding.EncodeToString(content)

	if attachmentType == "jpg" {
		attachmentType = "jpeg"
	}

	e.Attachments = append(e.Attachments, Attachment{
		Content:        b64content,
		AttachmentType: attachmentType,
		Name:           name,
	})
}

// AddAttachmentFile attaches a file from the filesystem to the email
func (e *Email) AddAttachmentFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filePath)
	if len(ext) == 0 {
		return fmt.Errorf("file extension not found")
	}

	name := filepath.Base(filePath)

	e.AddAttachmentRaw(name, data, ext[1:])

	return nil
}
