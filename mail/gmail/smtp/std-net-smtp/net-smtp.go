package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Printf("env が読み込み出来ませんでした: %v", err)
		return
	}

	hostname := os.Getenv("SMTPHOSTNAME")
	port := os.Getenv("PORT")
	to := os.Getenv("To")
	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")

	fmt.Println(hostname)
	fmt.Println(port)
	fmt.Println(to)
	fmt.Println(password)
	fmt.Println(from)

	recipients := []string{to}

	auth := smtp.PlainAuth("", from, password, hostname)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Test Email2\r\n" +
		"\r\n" +
		"This is a test email4.\r\n")

	err = smtp.SendMail(hostname+":"+port, auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")
}
