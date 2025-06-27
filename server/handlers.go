package server

import (
	"bytes"
	"errors"
	"fmt"
	"gemini/request"
	"gemini/response"
	"log"
	"net"
	"os"
	"path"
)

// handleRequest processes incoming TCP connections, processing the request and
// attempting to handle it with either a custom handler or by serving a static file.
func (s *Server) handleRequest(conn net.Conn) {
	var resp response.Response
	defer writeResponseAndClose(conn, &resp)

	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		message := fmt.Sprintf("Couldn't read: %s", err)
		resp = response.NewPermanentFailureResponse(response.PermanentFailureBadRequest, &message)
		return
	}
	if bytesRead > 1024 {
		message := "Request too long"
		resp = response.NewPermanentFailureResponse(response.PermanentFailureBadRequest, &message)
		return
	}

	requestStr := string(bytes.TrimSpace(buffer[:bytesRead]))
	request, err := request.ParseRequest(requestStr)

	if err != nil {
		message := fmt.Sprintf("Invalid request: %s", err)
		resp = response.NewPermanentFailureResponse(response.PermanentFailureBadRequest, &message)
		return
	}

	customResponse := s.handleCustomRequest(conn, request)
	if customResponse != nil {
		resp = *customResponse
	} else {
		resp = s.handleFileRequest(request)
	}
}

// handleCustomRequest executes a handler matching the request path
// It also retrieves the client certificate if available, passing it to the handler
// If retrieving the certificate fails, it logs the error and passes nil to the handler
func (s *Server) handleCustomRequest(conn net.Conn, request request.Request) *response.Response {
	cert, err := ClientCert(conn)
	if err != nil {
		log.Printf("Error getting client certificate: %s\n", err)
	}
	// custom handlers
	for path, handler := range s.handlers {
		if path == request.Url.Path {
			response := handler(request, cert)
			return &response
		}
	}
	return nil
}

// handleFileRequest serves static files from the configured StaticFilesPath.
// It checks if the requested file exists and returns it if found, or a 404 error
// It may return a PermanentFailure response if the file cannot be read.
func (s *Server) handleFileRequest(request request.Request) response.Response {
	gemfileDirPath := path.Join(s.Config.StaticFilesPath, path.Clean(request.Url.Path))
	hasExtension := path.Ext(gemfileDirPath) != ""
	if !hasExtension {
		gemfileDirPath = path.Join(gemfileDirPath, "/index.gmi")
	}
	if _, err := os.Stat(gemfileDirPath); errors.Is(err, os.ErrNotExist) {
		return response.NewPermanentFailureResponse(response.PermanentFailureNotFound, nil)
	}

	file, err := os.Open(gemfileDirPath)
	if err != nil {
		log.Printf("Could not open file: %s", err)
		message := "Resource could not be read"
		return response.NewPermanentFailureResponse(response.PermanentFailure, &message)
	}
	return response.NewSuccessResponse("text/gemini", file)
}
