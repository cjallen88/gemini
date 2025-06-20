package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"gemini/request"
	"gemini/response"
	"log"
	"net"
)

const (
	HOST = "localhost"
	PORT = "1965"
)

func main() {
	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
	if err != nil {
		log.Fatal(err)
	}
	// Configure the server to trust TLS client certs issued by a CA.
	// certPool := x509.SystemCertPool()

	listen, err := tls.Listen("tcp", HOST+":"+PORT, &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   HOST + ":" + PORT,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}

func writeResponseAndClose(conn net.Conn, resp *response.Response) {
	if resp != nil {
		_, err := (*resp).WriteToStream(conn) //conn.Write([]byte((*resp).String()))
		if err != nil {
			log.Println("Failed to write response", err.Error())
		}
	} else {
		log.Println("No response to write")
	}
	conn.Close()
}

func handleRequest(conn net.Conn) {
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
		log.Print(&message)
		resp = response.NewPermanentFailureResponse(response.PermanentFailureBadRequest, &message)
		return
	}

	requestStr := string(bytes.TrimSpace(buffer[:bytesRead]))
	request, err := request.ParseRequest(requestStr)

	if err != nil {
		message := fmt.Sprintf("Invalid request: %s\r\n", err)
		resp = response.NewPermanentFailureResponse(response.PermanentFailureBadRequest, &message)
	} else {
		resp = response.NewSuccessResponse("text/gemini", fmt.Sprintf("hello, world on %s!\r\n", request.Url.String()))
	}
}
