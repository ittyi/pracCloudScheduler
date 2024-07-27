package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

// SMTP server configuration.
const (
	smtpServer   = "smtp.gmail.com"
	smtpPort     = "465" // SSL/TLS port
	startTlsPort = "587" // STARTTLS port
)

// Email authentication credentials.
const (
	email    = ""
	password = ""
)

func sendEmailSSL() {
	// Recipient email address.
	to := ""

	// Email subject and body.
	subject := "Subject: Test Email\n"
	body := "This is a test email sent using Go and Gmail SMTP server over SSL/TLS."

	// Construct the email message.
	msg := []byte(subject + "\n" + body)

	// Set up authentication information.
	auth := smtp.PlainAuth("", email, password, smtpServer)

	// Set up the TLS configuration.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // In production, use a proper certificate verification
		ServerName:         smtpServer,
	}

	// Dial the SMTP server with TLS.
	conn, err := tls.Dial("tcp", smtpServer+":"+smtpPort, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a new SMTP client from the connection.
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	// Check the encryption used.
	state := conn.ConnectionState()
	fmt.Printf("SSL/TLS Connection: Version %x, Cipher %x\n", state.Version, state.CipherSuite)

	// Authenticate to the SMTP server.
	if err := client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient addresses.
	if err := client.Mail(email); err != nil {
		log.Fatal(err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Fatal(err)
	}

	// Get the data writer to send the email body.
	writer, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}

	// Write the email body and close the writer.
	if _, err := writer.Write(msg); err != nil {
		log.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully using SSL/TLS!")
}

func sendEmailSTARTTLS() {
	// Recipient email address.
	to := ""

	// Email subject and body.
	subject := "Subject: Test Email\n"
	body := "This is a test email sent using Go and Gmail SMTP server over STARTTLS."

	// Construct the email message.
	msg := []byte(subject + "\n" + body)

	// Set up authentication information.
	auth := smtp.PlainAuth("", email, password, smtpServer)

	// Connect to the SMTP server.
	conn, err := net.Dial("tcp", smtpServer+":"+startTlsPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a new SMTP client from the connection.
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Quit()

	// Start TLS encryption.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // In production, use a proper certificate verification
		ServerName:         smtpServer,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Fatal(err)
	}

	// Check the encryption used.
	state := conn.(*tls.Conn).ConnectionState()
	fmt.Printf("STARTTLS Connection: Version %x, Cipher %x\n", state.Version, state.CipherSuite)

	// Authenticate to the SMTP server.
	if err := client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient addresses.
	if err := client.Mail(email); err != nil {
		log.Fatal(err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Fatal(err)
	}

	// Get the data writer to send the email body.
	writer, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}

	// Write the email body and close the writer.
	if _, err := writer.Write(msg); err != nil {
		log.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully using STARTTLS!")
}

func main() {
	sendEmailSSL()
	sendEmailSTARTTLS()
}
