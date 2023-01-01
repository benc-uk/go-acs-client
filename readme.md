# Client SDK for the Azure Communication Services Email API

This is a SDK client library for use with Azure Communication Services email service. It is a an implementation of the [REST API documented here](https://learn.microsoft.com/en-us/rest/api/communication/email/send?tabs=HTTP)

It's simple to use, the only pre-req is a Email Communication Service resource [deployed and configured](https://learn.microsoft.com/en-us/azure/communication-services/quickstarts/email/create-email-communication-resource)  in Azure

## Simple Usage

Let's send an email to Bob, he'd love to hear from us...

```go
import "github.com/benc-uk/go-acs-email/client"

endpoint := "https://blahblah.communication.azure.com"
accessKey := os.Getenv("ACS_ACCESS_KEY") // Keep the ACS access key secret

acsClient := client.New(accessKey, endpoint)

email := client.NewHTMLEmail("DoNotReply@blah.net", "bob@bob.com", "Hello!", "<h1>Hi bob!</h1>")

msgID, err := acsClient.SendEmail(email)
if err != nil {
  log.Fatal(err)
}

// Optional - check or poll the status with client.GetStatus(msgID)
```

See the `client_tests.go` file for more detailed examples

## Quick Docs

A simplified summary of the main data types and functions 

### Type: `Client`

```go
// NewClient creates a client with the given access key and endpoint
func NewClient(accessKey, endpoint string) *Client

// SendEmail sends an email and returns the message ID and any error
func (c *Client) SendEmail(e *Email) (messageID string, err error)

// GetStatus gets the status of an email message sent using SendEmail()
func (c *Client) GetStatus(messageID string) (status string, err error)
```

### Type: `Email`

```go
// Email is the main request type to send emails
type Email struct {
        Recipients  Recipients     `json:"recipients"`
        Sender      string         `json:"sender"`
        Content     Content        `json:"content"`
        Headers     []CustomHeader `json:"headers"`
        Tracking    bool           `json:"disableUserEngagementTracking"`
        Importance  string         `json:"importance"`
        ReplyTo     []Address      `json:"replyTo"`
        Attachments []Attachment   `json:"attachments"`
}

// NewHTMLEmail creates a new email with HTML content
func NewHTMLEmail(from, to, subject, body string) *Email

// NewPlainEmail creates a new email with plain text content
func NewPlainEmail(from, to, subject, body string) *Email

// AddAttachmentFile attaches a file from the filesystem to the email
func (e *Email) AddAttachmentFile(filePath string) error

// AddAttachmentRaw adds an attachment to the email as raw bytes
func (e *Email) AddAttachmentRaw(name string, content []byte, attachmentType string)

// AddBCC adds a BCC recipient
func (e *Email) AddBCC(address, displayName string)

// AddCC adds a CC recipient
func (e *Email) AddCC(address, displayName string)

// AddCustomHeader adds a custom header to the email
func (e *Email) AddCustomHeader(key, value string)

// AddReplyTo adds a reply to address
func (e *Email) AddReplyTo(address, displayName string)

// DisableUserEngagementTracking disables user engagement tracking
func (e *Email) DisableUserEngagementTracking()

// EnableUserEngagementTracking enables user engagement tracking
func (e *Email) EnableUserEngagementTracking()
```
