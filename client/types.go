package client

// Types for the ACS Email API
// See: https://learn.microsoft.com/en-us/rest/api/communication/email/send?tabs=HTTP

const ImportanceLow = "low"
const ImportanceNormal = "normal"
const ImportanceHigh = "high"

// ==== Request types ====

// Email is the main request type
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

// Recipients contains the To, CC and BCC recipients of the email
type Recipients struct {
	To  []Address `json:"to"`
	CC  []Address `json:"cc"`
	BCC []Address `json:"bcc"`
}

// Address contains the email address and display name
type Address struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

// Content contains the subject, HTML and/or plain text body of the email
type Content struct {
	Subject   string `json:"subject"`
	HTML      string `json:"html"`
	PlainText string `json:"plainText"`
}

// CustomHeader contains a custom header name and value
type CustomHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Attachment is used to attach files to the email
type Attachment struct {
	Content        string `json:"contentBytesBase64"`
	AttachmentType string `json:"attachmentType"`
	Name           string `json:"name"`
}

// ==== Error types ====

// ErrorResponse wraps the error response from the API
type ErrorResponse struct {
	Error CommunicationError `json:"error"`
}

// CommunicationError contains the error code and message
type CommunicationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ==== Status types ====

// SendStatusResult contains the message ID and status of the email
type SendStatusResult struct {
	MessageID string `json:"messageId"`
	Status    string `json:"status"`
}
