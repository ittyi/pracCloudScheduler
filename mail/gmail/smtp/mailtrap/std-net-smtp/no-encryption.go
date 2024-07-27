package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")

	if err != nil {
		fmt.Printf("env が読み込み出来ませんでした: %v", err)
		return
	}

	hostname := os.Getenv("SMTPHOSTNAME")
	port := os.Getenv("PORT")
	TO := os.Getenv("To")
	from := os.Getenv("FROM")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	fmt.Println(hostname)
	fmt.Println(port)
	fmt.Println(TO)
	fmt.Println(from)
	fmt.Println(username)
	fmt.Println(password)

	// MailtrapのSMTPサーバーの設定
	smtpServer := hostname

	fmt.Println("PlainAuth: ")
	auth := smtp.PlainAuth("", username, password, smtpServer)
	fmt.Println("auth: ", auth)

	// メールの内容
	to := []string{TO}
	subject := "Subject: Test Email\n"
	body := "This is the email body."
	msg := []byte(subject + "\n" + body)

	// メールの送信
	err = smtp.SendMail(smtpServer+":"+port, auth, from, to, msg)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
