package response

import (
	"bufio"
	"fmt"
	"io"
)

type SuccessResponse struct {
	MimeType string
	Body     bufio.Reader
}

func (r *SuccessResponse) WriteTo(w io.Writer) (int64, error) {
	n1, err := fmt.Fprintf(w, "20 %s\r\n", r.MimeType)
	if err != nil {
		return int64(n1), err
	}
	n2, err := r.Body.WriteTo(w)
	total := int64(n1) + n2
	if err != nil {
		return total, err
	}
	return total, nil
}

func NewSuccessResponse(mimeType string, body io.Reader) *SuccessResponse {
	bufferedReader := bufio.NewReader(body)
	return &SuccessResponse{
		MimeType: mimeType,
		Body:     *bufferedReader,
	}
}
