package response

import (
	"fmt"
	"io"
)

type ClientCertificateStatus Status

const (
	CertificateRequired      ClientCertificateStatus = 60
	CertificateNotAuthorized ClientCertificateStatus = 62
	CertificateNotValid      ClientCertificateStatus = 63
)

func (r *ClientCertificateStatus) DefaultMessage() string {
	switch *r {
	case CertificateNotAuthorized:
		return "Client certificate not authorized"
	case CertificateNotValid:
		return "Client certificate not valid"
	}
	return "Client certificate required"
}

type ClientCertificatesResponse struct {
	Status  ClientCertificateStatus
	Message *string
}

func (r *ClientCertificatesResponse) WriteToStream(w io.Writer) (int, error) {
	msg := *r.Message
	if r.Message == nil {
		msg = r.Status.DefaultMessage()
	}
	return fmt.Fprintf(w, "%d %s\r\n", r.Status, msg)
}

func NewClientCertificatesResponse(status ClientCertificateStatus, message *string) *ClientCertificatesResponse {
	return &ClientCertificatesResponse{
		Status:  status,
		Message: message,
	}
}
