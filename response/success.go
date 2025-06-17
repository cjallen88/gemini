package response

import (
	"fmt"
	"io"
)

type SuccessResponse struct {
	MimeType string
	Body     string
}

func (r *SuccessResponse) WriteToStream(w io.Writer) (int, error) {
	// body will at some point come from a file
	return fmt.Fprintf(w, "20 %s\r\n%s", r.MimeType, r.Body)
}

func NewSuccessResponse(mimeType, body string) *SuccessResponse {
	return &SuccessResponse{
		MimeType: mimeType,
		Body:     body,
	}
}
