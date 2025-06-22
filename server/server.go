package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"gemini/request"
	"gemini/response"
	"log"
	"net"
)

type Config struct {
	Host            string
	Port            string
	CertificatePath string
	KeyPath         string
}

type RequestHandler func(request request.Request, clientCert *Certificate) response.Response

func Serve(config *Config, handler RequestHandler) error {
	if config == nil {
		return fmt.Errorf("server configuration cannot be nil")
	}

	cert, err := tls.LoadX509KeyPair(config.CertificatePath, config.KeyPath)
	if err != nil {
		return fmt.Errorf("Failed to load TLS certificate: %w", err)
	}

	serverName := fmt.Sprintf("%s:%s", config.Host, config.Port)
	listen, err := tls.Listen("tcp", serverName, &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
		ClientAuth:   tls.RequestClientCert, // validate?
	})
	if err != nil {
		return fmt.Errorf("Failed to start TLS listener: %w", err)
	}
	defer listen.Close()

	log.Printf("Gemini server listening on %s\n", config.Port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleRequest(conn, handler)
	}
}

func writeResponseAndClose(conn net.Conn, resp *response.Response) {
	if resp != nil {
		_, err := (*resp).WriteToStream(conn)
		if err != nil {
			log.Println("Failed to write response", err.Error())
		}
	} else {
		log.Println("No response to write")
	}
	conn.Close()
}

func handleRequest(conn net.Conn, handle RequestHandler) {
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

	cert, err := ClientCert(conn)
	if err != nil {
		message := fmt.Sprintf("Error getting client certificate: %s", err)
		log.Print(message)
		resp = response.NewPermanentFailureResponse(response.PermanentFailureBadRequest, &message)
		return
	}

	resp = handle(request, cert)
}
