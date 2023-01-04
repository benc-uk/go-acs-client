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
	"log"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

var toAddress string
var ccAddress string
var fromAddress string
var endpoint string
var accessKey string

var fromNumber string
var toNumber string

const subject = "Test email via Azure Communication Services"
const emailBody = "<h1>Hello!</h1>This email was sent using Go and the Azure Communication Services REST API"

func TestMain(m *testing.M) {
	_ = godotenv.Load("./.env")
	_ = godotenv.Load("../.env")

	endpoint = os.Getenv("ACS_ENDPOINT")
	accessKey = os.Getenv("ACS_ACCESS_KEY")

	toAddress = os.Getenv("TO_ADDRESS")
	fromAddress = os.Getenv("FROM_ADDRESS")
	ccAddress = os.Getenv("CC_ADDRESS")

	fromNumber = os.Getenv("FROM_NUMBER")
	toNumber = os.Getenv("TO_NUMBER")

	if endpoint == "" || accessKey == "" {
		log.Fatal("Please set ACS_ENDPOINT and ACS_ACCESS_KEY")
	}

	if toAddress == "" || fromAddress == "" || ccAddress == "" {
		log.Fatal("Please set TO_ADDRESS, FROM_ADDRESS & CC_ADDRESS")
	}

	if fromNumber == "" || toNumber == "" {
		log.Fatal("Please set FROM_NUMBER, TO_NUMBER")
	}

	m.Run()
}

func TestSendSimple(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, subject, emailBody)
	_, err := client.SendEmail(e)

	if err != nil {
		t.Error(err)
	}
}

func TestSendStatus(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, subject, emailBody)

	id, err := client.SendEmail(e)
	if err != nil {
		t.Error(err)
	}

	status, err := client.GetEmailStatus(id)
	if err != nil {
		t.Error(err)
	}

	if status != "Queued" {
		t.Error("Expected status to be 'Queued', got: " + status)
	}
}

func TestSendCC(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, subject, emailBody)
	e.AddCC(ccAddress, "Some person")
	_, err := client.SendEmail(e)

	if err != nil {
		t.Error(err)
	}
}

func TestSendAttachmentText(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, "Testing text attachments", "Yo! Here are some attachments...")

	e.AddAttachmentRaw("hello.txt", []byte("Hello world!"), "txt")
	_ = e.AddAttachmentFile("testdata/test_file.txt")

	_, err := client.SendEmail(e)
	if err != nil {
		t.Error(err)
	}
}

func TestSendAttachmentImage(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, "Testing image attachments", "Yo! Here are some attachments...")

	_ = e.AddAttachmentFile("testdata/trek.gif")
	_ = e.AddAttachmentFile("testdata/moss.jpg")

	_, err := client.SendEmail(e)
	if err != nil {
		t.Error(err)
	}
}

func TestAttachmentMissing(t *testing.T) {
	e := NewHTMLEmail(fromAddress, toAddress, subject, emailBody)

	err := e.AddAttachmentFile("doesnt_exist.docx")
	if err == nil {
		t.Error(err)
	}
}

func TestSendCustomHeader(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewPlainEmail(fromAddress, toAddress, "Testing with custom header", "I wonder how my socks are?")
	e.AddCustomHeader("X-Sock-Status", "My socks are extremely smelly")
	_, err := client.SendEmail(e)

	if err != nil {
		t.Error(err)
	}
}

func TestInvalidToAddress(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, "lemon", subject, emailBody)
	_, err := client.SendEmail(e)

	if err == nil || !strings.Contains(err.Error(), "Error setting value to 'Email'") {
		t.Error("Expected error, but got:", err)
	}
}

func TestInvalidFromAddress(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail("sausages", toAddress, subject, emailBody)
	_, err := client.SendEmail(e)

	if err == nil || !strings.Contains(err.Error(), "Error setting value to 'Sender'") {
		t.Error("Expected error, but got:", err)
	}
}

func TestNoSubject(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, "", emailBody)
	_, err := client.SendEmail(e)

	if err == nil || !strings.Contains(err.Error(), "Email should contain a non-empty subject") {
		t.Error("Expected error, but got:", err)
	}
}

func TestNoBody(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, subject, "")
	_, err := client.SendEmail(e)

	if err == nil || !strings.Contains(err.Error(), "Email body validation error") {
		t.Error("Expected error, but got:", err)
	}
}

func TestInvalidImportance(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, subject, emailBody)
	e.Importance = "fishcake"
	_, err := client.SendEmail(e)

	if err == nil || !strings.Contains(err.Error(), "EmailImportance") {
		t.Error("Expected error, but got:", err)
	}
}

func TestInvalidReplyTo(t *testing.T) {
	client := New(accessKey, endpoint)
	e := NewHTMLEmail(fromAddress, toAddress, subject, emailBody)
	e.ReplyTo = []Address{
		{
			DisplayName: "dave",
			Email:       "cheese_burger",
		},
	}
	_, err := client.SendEmail(e)

	if err == nil || !strings.Contains(err.Error(), "Error setting value to 'Email'") {
		t.Error("Expected error, but got:", err)
	}
}
