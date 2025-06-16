package response

import "fmt"

type SuccessResponse struct {
	MimeType string
	Body     string
}

func (r *SuccessResponse) String() string {
	return fmt.Sprintf("20 %s\r\n%s", r.MimeType, r.Body)
}

func NewSuccessResponse(mimeType, body string) *SuccessResponse {
	return &SuccessResponse{
		MimeType: mimeType,
		Body:     body,
	}
}
