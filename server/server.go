package server

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"gemini/request"
	"gemini/response"
	"log"
	"net"
	"os"
	"path"
)

type Config struct {
	Host            string
	Port            string
	CertificatePath string
	KeyPath         string
	StaticFilesPath string
}

type RequestHandler func(request request.Request, clientCert *Certificate) response.Response

type Server struct {
	handlers map[string]RequestHandler
	Config   Config
}

func NewServer(config Config) *Server {
	return &Server{
		handlers: make(map[string]RequestHandler),
		Config:   config,
	}
}

func (s *Server) CustomHandler(requestPath string, handler RequestHandler) {
	s.handlers[requestPath] = handler
}

func (s *Server) Serve() {
	cert, err := tls.LoadX509KeyPair(s.Config.CertificatePath, s.Config.KeyPath)
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %s", err)
	}

	serverName := fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port)
	listen, err := tls.Listen("tcp", serverName, &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
		ClientAuth:   tls.RequestClientCert, // validate?
	})
	if err != nil {
		log.Fatalf("Failed to start TLS listener: %s", err)
	}
	defer listen.Close()

	log.Printf("Gemini server listening on %s\n", s.Config.Port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go s.handleRequest(conn)
		// TODO prevent more handlers being added while handling requests
	}
}

func writeResponseAndClose(conn net.Conn, resp *response.Response) {
	if resp != nil {
		_, err := (*resp).WriteTo(conn)
		if err != nil {
			log.Println("Failed to write response", err.Error())
		}
	} else {
		log.Println("No response to write")
	}
	conn.Close()
}

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
		resp = s.handleFileRequest(conn, request)
	}
}

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

func (s *Server) handleFileRequest(conn net.Conn, request request.Request) response.Response {
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
