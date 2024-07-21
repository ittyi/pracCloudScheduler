package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	mail "github.com/xhit/go-simple-mail/v2"
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

	var portNum int
	portNum, _ = strconv.Atoi(port)

	// SMTPサーバーの設定
	server := mail.NewSMTPClient()
	server.Host = hostname
	server.Port = portNum
	server.Username = from
	server.Password = password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// SMTPサーバーに接続
	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// メールの設定
	email := mail.NewMSG()
	email.SetFrom("From " + "Ittyi" + " <" + from + ">")
	email.AddTo(to)
	email.SetSubject("by go-mail")
	email.SetBody(mail.TextPlain, "This is a test email sent from Go using go-mail package.")

	// メールの送信
	err = email.Send(smtpClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}
