package response

import (
	"fmt"
	"io"
)

type PermanentFailureStatus Status

const (
	PermanentFailure                    PermanentFailureStatus = 50
	PermanentFailureNotFound            PermanentFailureStatus = 51
	PermanentFailureGone                PermanentFailureStatus = 52
	PermanentFailureProxyRequestRefused PermanentFailureStatus = 53
	PermanentFailureBadRequest          PermanentFailureStatus = 59
)

func (r *PermanentFailureStatus) DefaultMessage() string {
	switch *r {
	case PermanentFailureNotFound:
		return "This resource was not found"
	case PermanentFailureGone:
		return "This resource is no longer available"
	case PermanentFailureProxyRequestRefused:
		return "The proxy server rejected the request"
	case PermanentFailureBadRequest:
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
