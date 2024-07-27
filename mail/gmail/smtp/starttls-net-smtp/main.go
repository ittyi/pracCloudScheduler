package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// SMTP server configuration.
const (
	smtpServer = "smtp.gmail.com"
	smtpPort   = "587"
)

// Email authentication credentials.
const (
	email    = ""
	password = ""
)

func main() {
	// Recipient email address.
	to := ""

	// Email subject and body.
	subject := "Subject: Test Email\n"
	body := "This is a test email sent using Go and Gmail SMTP server with STARTTLS."

	// Construct the email message.
	msg := []byte(subject + "\n" + body)

	// Connect to the SMTP server.
	client, err := smtp.Dial(smtpServer + ":" + smtpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	printConnectionState(client)

	// Start TLS encryption.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set this to false in production
		ServerName:         smtpServer,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		log.Fatal(err)
	}

	// Authenticate to the SMTP server.
	auth := smtp.PlainAuth("", email, password, smtpServer)
	if err := client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	log.Println("二回目")
	printConnectionState(client)

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

	// Close the SMTP client.
	if err := client.Quit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}

func printConnectionState(conn *smtp.Client) {
	state, ok := conn.TLSConnectionState()
	if ok {
		fmt.Printf("Version: %x\n", state.Version)
		fmt.Printf("HandshakeComplete: %t\n", state.HandshakeComplete)
		fmt.Printf("CipherSuite: %x\n", state.CipherSuite)
		fmt.Printf("ServerName: %s\n", state.ServerName)
	}
}
