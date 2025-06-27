package response

import "io"

type Status int

type Response interface {
	WriteTo(w io.Writer) (int64, error)
}
