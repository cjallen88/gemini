package server

import (
	"crypto/tls"
	"fmt"
	"gemini/server/request"
	"gemini/server/response"
	"log"
	"net"
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

// NewServer creates a new Gemini server instance with the provided configuration.
func NewServer(config Config) *Server {
	return &Server{
		handlers: make(map[string]RequestHandler),
		Config:   config,
	}
}

// HandlePath registers a custom request handler for a specific path.
// This takes precedence over serving static files, and can be used to implement dynamic content or custom logic.
func (s *Server) HandlePath(requestPath string, handler RequestHandler) {
	s.handlers[requestPath] = handler
}

// Serve starts the Gemini server, listening for incoming TLS connections.
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
	}
}

// writeResponseAndClose function to write a response to the connection and close it
// should be called in a defer statement to ensure it always runs
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
