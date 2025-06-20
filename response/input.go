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

func (r *InputResponse) WriteToStream(w io.Writer) (int, error) {
	return fmt.Fprintf(w, "%d %s\r\n", r.Status, r.Prompt)
}

func NewInputResponse(status InputStatus, prompt string) *InputResponse {
	return &InputResponse{
		Status: status,
		Prompt: prompt,
	}
}
