package response

import (
	"fmt"
	"io"
)

type TemporaryFailureStatus Status

const (
	TemporaryFailure  TemporaryFailureStatus = 40
	ServerUnavailable TemporaryFailureStatus = 41
	CGIError          TemporaryFailureStatus = 42
	ProxyError        TemporaryFailureStatus = 43
	SlowDown          TemporaryFailureStatus = 44
)

func (r *TemporaryFailureStatus) DefaultMessage() string {
	switch *r {
	case ServerUnavailable:
		return "The server is currently unavailable"
	case CGIError:
		return "The server encountered an error while processing the request via CGI"
	case ProxyError:
		return "The proxy server encountered an error while processing the request"
	case SlowDown:
		return "The server is currently overloaded, please slow down your requests"
	}
	return "The server has encountered an error, please try again later"
}

type TemporaryFailureResponse struct {
	Status  TemporaryFailureStatus
	Message *string
}

func (r *TemporaryFailureResponse) WriteTo(w io.Writer) (int64, error) {
	var msg string
	if r.Message == nil {
		msg = r.Status.DefaultMessage()
	} else {
		msg = *r.Message
	}
	n, error := fmt.Fprintf(w, "%d %s\r\n", r.Status, msg)
	return int64(n), error
}

func NewTemporaryFailureResponse(status TemporaryFailureStatus, message *string) *TemporaryFailureResponse {
	return &TemporaryFailureResponse{
		Status:  status,
		Message: message,
	}
}
