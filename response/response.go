package response

import "io"

type Status int

type Response interface {
	WriteToStream(w io.Writer) (int, error)
}
