package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 環境変数の読み込み
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("env が読み込み出来ませんでした: %v\n", err)
		return
	}

	// 環境変数から設定を取得
	hostname := os.Getenv("SMTPHOSTNAME")
	port := "465" // SSLの場合、ポート465を使用
	to := os.Getenv("To")
	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")

	// 確認用に設定内容を出力
	fmt.Println("SMTP Hostname:", hostname)
	fmt.Println("SMTP Port:", port)
	fmt.Println("To:", to)
	fmt.Println("Password:", password)
	fmt.Println("From:", from)

	// SMTPサーバーの設定
	smtpServer := hostname
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SSL/TLS設定の作成
	tlsconfig := &tls.Config{
		// 証明書の検証を有効にする
		InsecureSkipVerify: false,
		ServerName:         smtpServer,
	}

	// 接続の作成
	addr := smtpServer + ":" + port
	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// SMTPクライアントの作成
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// 認証の実行
	if err := client.Auth(auth); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// メールの送信
	if err := client.Mail(from); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Fatalf("Failed to add recipient: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to get writer for email data: %v", err)
	}
	defer w.Close()

	msg := []byte("Subject: Test Email\n\nThis is the email body.")
	if _, err := w.Write(msg); err != nil {
		log.Fatalf("Failed to write email data: %v", err)
	}

	log.Println("Email sent successfully!")
}
