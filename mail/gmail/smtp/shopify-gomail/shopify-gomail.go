package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
)

func getDomainFromEmail(email string) (string, error) {
	// @ 記号でメールアドレスを分割
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email address format")
	}
	return parts[1], nil
}

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

	// メールの設定
	m := mail.NewMessage()

	domain, err := getDomainFromEmail(from)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Domain:", domain)
	// return

	// ユニークなMessage-IDを生成
	messageID := fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), "unique", domain)

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Bcc", from)
	m.SetHeader("Message-ID", messageID)
	m.SetHeader("Subject", "by Shopify/gomail")
	m.SetBody("text/plain", "This is a test email sent from Go using Shopify/gomail package.")

	var portNum int
	portNum, _ = strconv.Atoi(port)

	// ダイヤラの設定
	d := mail.NewDialer(hostname, portNum, from, password)
	// d.TLSConfig = nil // Disable TLS
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// メールの送信
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}
