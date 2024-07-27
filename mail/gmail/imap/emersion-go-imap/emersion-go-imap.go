package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

func main() {
	// GetSendbox()
	Listmailbox()
}

func GetSendbox() {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Printf("env が読み込み出来ませんでした: %v", err)
		return
	}

	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")
	to := os.Getenv("To")
	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")

	fmt.Println(hostname)
	fmt.Println(port)
	fmt.Println(to)
	fmt.Println(password)
	fmt.Println(from)

	recipients := []string{""}

	auth := smtp.PlainAuth("", from, password, hostname)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Test Email2\r\n" +
		"\r\n" +
		"This is a test email9.\r\n")

	err = smtp.SendMail(hostname+":"+port, auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")

	// サーバーに接続
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("c: ", c)
	log.Println("Connected")

	// ログイン
	if err := c.Login(from, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// 送信済みフォルダを選択
	mbox, err := c.Select("[Gmail]/送信済みメール", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for Sent Mail:", mbox.Flags)

	// メールの検索
	criteria := imap.NewSearchCriteria()
	criteria.Since = time.Now().Add(-1 * time.Minute) // 直近1分以内のメールを検索
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ids: ", ids)

	if len(ids) == 0 {
		log.Println("No recent sent messages")
		return
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {
		log.Println("Recent sent message:", msg.Envelope.Subject)
		log.Println("Recent sent MessageId:", msg.Envelope.MessageId)
		log.Println("Recent sent Format():", msg.Envelope.Format())
	}

	// ログアウト
	if err := c.Logout(); err != nil {
		log.Fatal(err)
	}
}

func Listmailbox() {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Printf("env が読み込み出来ませんでした: %v", err)
		return
	}

	// hostname := os.Getenv("HOSTNAME")
	// port := os.Getenv("PORT")
	// to := os.Getenv("To")
	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")

	// サーバーに接続
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// ログイン
	if err := c.Login(from, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// メールボックス一覧を取得
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// ログアウト
	if err := c.Logout(); err != nil {
		log.Fatal(err)
	}
}
