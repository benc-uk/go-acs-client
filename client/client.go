package client

// ==============================================================================
// Client SDK for the Azure Communication Services API
// An implementation of the API documented here:
// https://learn.microsoft.com/en-us/rest/api/communication/email/send
// ==============================================================================

const sendEmailEndpoint = "/emails:send"
const sendSMSEndpoint = "/sms"
const statusEmailEndpoint = "/emails/%s/status"
const clientTimeout = 20

// Client is used to send emails with Azure Communication Services
type Client struct {
	AccessKey       string
	Endpoint        string
	APIVersionEmail string // Defaults to 2021-10-01-preview
	APIVersionSMS   string // Defaults to 2021-03-07
}

// New creates a client with the given access key and endpoint
func New(accessKey, endpoint string) *Client {
	return &Client{
		AccessKey:       accessKey,
		Endpoint:        endpoint,
		APIVersionEmail: "2021-10-01-preview",
		APIVersionSMS:   "2021-03-07",
	}
}
