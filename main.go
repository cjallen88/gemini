package main

import (
	"fmt"
	"gemini/request"
	"gemini/response"
	"gemini/server"
	"log"
)

func main() {
	config := &server.Config{
		Host:            "localhost",
		Port:            "1965",
		CertificatePath: "localhost.crt",
		KeyPath:         "localhost.key",
	}

	err := server.Serve(config, handleRequest)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleRequest(request request.Request, cert *server.Certificate) response.Response {
	if cert == nil {
		message := "Client certificate is required for this request"
		return response.NewClientCertificatesResponse(response.CertificateRequired, &message)
	}
	return response.NewSuccessResponse("text/gemini", fmt.Sprintf("hello to %s, world on %s!\r\n with fingerprint %s", cert.Name, request.Url.String(), cert.Fingerprint))
}
