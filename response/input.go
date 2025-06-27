package response

import (
	"fmt"
	"io"
)

type InputStatus Status

const (
	InputBasic     InputStatus = 10
	InputSensitive InputStatus = 11
)

type InputResponse struct {
	Status InputStatus
	Prompt string
}

func (r *InputResponse) WriteTo(w io.Writer) (int64, error) {
	n, error := fmt.Fprintf(w, "%d %s\r\n", r.Status, r.Prompt)
	return int64(n), error
}

func NewInputResponse(status InputStatus, prompt string) *InputResponse {
	return &InputResponse{
		Status: status,
		Prompt: prompt,
	}
}
