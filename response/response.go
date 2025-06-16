package response

type Status int

type Response interface {
	String() string
}

// func (r *SuccessResponse) Type() Type {
// 	return Success
// }

// type Response interface {
// 	String()
// }

// func (r *Response) String() string {
// 	return fmt.Sprintf("20 %d\r\n%d", r.MimeType, r.Body)
// }
