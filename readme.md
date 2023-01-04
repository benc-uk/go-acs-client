# Client SDK for the Azure Communication Services API

This is a SDK client library for use with Azure Communication Services, providing support for sending email and SMS.
It is a an implementation of the REST APIs documented [here](https://learn.microsoft.com/en-us/rest/api/communication/sms/send?tabs=HTTP) & [here](https://learn.microsoft.com/en-us/rest/api/communication/email/send?tabs=HTTP)

It's simple to use, the only pre-req is a Communication Service resource [deployed and configured](https://learn.microsoft.com/en-gb/azure/communication-services/quickstarts/create-communication-resource?tabs=windows&pivots=platform-azp)  in Azure

## Simple Usage

### Sending Email

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

### Sending a SMS

Maybe Bob would like a coffee, let's send him a SMS

```go
import "github.com/benc-uk/go-acs-email/client"

endpoint := "https://blahblah.communication.azure.com"
accessKey := os.Getenv("ACS_ACCESS_KEY") // Keep the ACS access key secret

acsClient := client.New(accessKey, endpoint)

sms := client.NewSMS("+18551111111", "+441234567890", "Fancy a coffee?")

smsResp, err := acsClient.SendSingleSMS(sms)
if err != nil {
  log.Fatal(err)
}

// Validate sending by checking the smsResp here
```

See the `email_test.go` & `sms_test.go` files for more detailed examples

## Quick Docs

A simplified summary of the main data types and functions 

### Type: `Client`

```go
// NewClient creates a client with the given access key and endpoint
func NewClient(accessKey, endpoint string) *Client

// SendEmail sends an email and returns the message ID and any error
func (c *Client) SendEmail(e *Email) (messageID string, err error)

// GetStatus gets the status of an email message sent using SendEmail()
func (c *Client) GetEmailStatus(messageID string) (status string, err error)

// SendSingleSMS sends a single SMS and returns the API response and/or error
func (c *Client) SendSingleSMS(s *SMS) (smsResp *SMSSendResponseItem, err error)
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

### Type `SMS`

```go
type SMS struct {
	From           string         `json:"from"`
	Message        string         `json:"message"`
	SMSRecipients  []SMSRecipient `json:"smsRecipients"`
	SMSSendOptions SMSOptions     `json:"smsSendOptions"`
}

// NewSMS creates a new SMS message for sending
func NewSMS(from, to, msg string) *SMS {
```