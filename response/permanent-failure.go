package response

import (
	"fmt"
	"io"
)

const (
	Generic             Status = 50
	NotFound            Status = 51
	Gone                Status = 52
	ProxyRequestRefused Status = 53
	BadRequest          Status = 59
)

type PermanentFailureResponse struct {
	Status  Status
	Message string
}

func (r *PermanentFailureResponse) WriteToStream(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "%d %s\r\n", r.Status, r.Message)
}

func NewPermanentFailureResponse(status Status, message string) *PermanentFailureResponse {
	return &PermanentFailureResponse{
		Status:  status,
		Message: message,
	}
}
