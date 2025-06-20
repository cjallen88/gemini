package response

import (
	"fmt"
	"io"
)

type PermanentFailureStatus Status

const (
	PermenentFailure                    PermanentFailureStatus = 50
	PermenentFailureNotFound            PermanentFailureStatus = 51
	PermenentFailureGone                PermanentFailureStatus = 52
	PermenentFailureProxyRequestRefused PermanentFailureStatus = 53
	PermenentFailureBadRequest          PermanentFailureStatus = 59
)

func (r *PermanentFailureStatus) DefaultMessage() string {
	switch *r {
	case PermenentFailureNotFound:
		return "This resource was not found"
	case PermenentFailureGone:
		return "This resource is no longer available"
	case PermenentFailureProxyRequestRefused:
		return "The proxy server rejected the request"
	case PermenentFailureBadRequest:
		return "The server was unable to understand the request"
	}
	return "The server has encountered an error"
}

type PermanentFailureResponse struct {
	Status  PermanentFailureStatus
	Message *string
}

func (r *PermanentFailureResponse) WriteToStream(w io.Writer) (int, error) {
	msg := *r.Message
	if r.Message == nil {
		msg = r.Status.DefaultMessage()
	}
	return fmt.Fprintf(w, "%d %s\r\n", r.Status, msg)
}

func NewPermanentFailureResponse(status PermanentFailureStatus, message *string) *PermanentFailureResponse {
	return &PermanentFailureResponse{
		Status:  status,
		Message: message,
	}
}
