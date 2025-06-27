package main

import (
	"fmt"
	"gemini/request"
	"gemini/response"
	"gemini/server"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	staticFilesPath := path.Join(basepath, "/static")

	server := server.NewServer(server.Config{
		Host:            "localhost",
		Port:            "1965",
		CertificatePath: "localhost.crt",
		KeyPath:         "localhost.key",
		StaticFilesPath: staticFilesPath,
	})

	server.CustomHandler("/game", handleGame)
	server.Serve()
}

func handleGame(request request.Request, cert *server.Certificate) response.Response {
	if cert == nil {
		message := "Client certificate is required for this request"
		return response.NewClientCertificatesResponse(response.CertificateRequired, &message)
	}
	message := fmt.Sprintf("hello to %s, world on %s!", cert.Name, request.Url.String())
	return response.NewSuccessResponse("text/gemini", strings.NewReader(message))
}
