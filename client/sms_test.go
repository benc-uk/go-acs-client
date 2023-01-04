package client

// ==============================================================================
// Test suite for the client SDK and email
// Create .env file with the following:
//  ACS_ENDPOINT=https://<your-resource-name>.communication.azure.com
//  ACS_ACCESS_KEY=<your-acs-access-key>
//  TO_ADDRESS=<your-email-address>
//  FROM_ADDRESS=<valid-from-address>
//  CC_ADDRESS=<other-email-address>
// ==============================================================================

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

const smsMessage = "Test SMS from Azure Communication Services"

func TestSendSMS(t *testing.T) {
	client := New(accessKey, endpoint)

	s := NewSMS(fromNumber, toNumber, smsMessage)

	r, err := client.SendSingleSMS(s)
	if err != nil {
		t.Error(err)
	}

	if r.Successful != true {
		t.Error("SMS was not sent successfully")
	}
}

func TestSendSMSNoMessage(t *testing.T) {
	client := New(accessKey, endpoint)

	s := NewSMS(fromNumber, toNumber, "")

	_, err := client.SendSingleSMS(s)
	if err == nil {
		t.Error(err)
	}
}

func TestSendSMSBadFrom(t *testing.T) {
	client := New(accessKey, endpoint)

	s := NewSMS("hello", toNumber, smsMessage)

	_, err := client.SendSingleSMS(s)
	if err == nil {
		t.Error(err)
	}
}

func TestSendSMSBadTo(t *testing.T) {
	client := New(accessKey, endpoint)

	s := NewSMS(fromNumber, "goats", smsMessage)

	r, err := client.SendSingleSMS(s)

	// NOTE! The API does NOT return a HTTP error if the 'to' number is invalid
	// This is because there could be multiple recipients and some may be valid
	if r.Successful == true {
		t.Error("SMS with invalid number was sent successfully")
	}

	if err != nil {
		t.Error(err)
	}
}
