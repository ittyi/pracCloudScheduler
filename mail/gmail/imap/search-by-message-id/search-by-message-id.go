package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Printf("env が読み込み出来ませんでした: %v", err)
		return
	}

	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")

	// fmt.Println(hostname)
	// fmt.Println(port)
	// fmt.Println(to)
	fmt.Println(password)
	fmt.Println(from)

	// Gmail のサーバーに接続
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// ログイン（アプリパスワードを使用）
	if err := c.Login(from, password); err != nil {
		log.Fatal(err)
	}

	// メールボックスを選択
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for Sent Mail:", mbox.Flags)

	// メッセージIDによる検索
	criteria := imap.NewSearchCriteria()
	criteria.Header.Add("Message-ID", "<>") // 例: <unique-message-id@example.com>

	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	msgs := make(chan *imap.Message, len(uids))
	go func() {
		if err := c.Fetch(seqSet, items, msgs); err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range msgs {
		r := msg.GetBody(section)
		if r == nil {
			log.Fatal("Server didn't return message body")
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		header := mr.Header

		messageID, err := header.Text("Message-ID")
		if err != nil {
			log.Println("Message-ID ヘッダーが存在しません")
		} else {
			fmt.Println("Message-ID:", messageID)
		}

		references, err := header.Text("References")
		if err != nil {
			log.Println("References ヘッダーが存在しません")
		} else {
			fmt.Println("References:", references)
		}

		inReplyTo, err := header.Text("In-Reply-To")
		if err != nil {
			log.Println("In-Reply-To ヘッダーが存在しません")
		} else {
			fmt.Println("In-Reply-To:", inReplyTo)
		}

		// メッセージの各パートを解析
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch h := part.Header.(type) {
			case *mail.InlineHeader:
				// インラインパートの処理（本文など）
				buf := new(bytes.Buffer)
				buf.ReadFrom(part.Body)
				fmt.Println("本文:", buf.String())
			case *mail.AttachmentHeader:
				// 添付ファイルの処理
				filename, _ := h.Filename()
				file, err := os.Create(filepath.Join("attachments", filename))
				fmt.Println("file:", file)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				log.Println("part.Body:", part.Body)
				_, err = file.ReadFrom(part.Body)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("添付ファイルを保存しました: %s\n", filename)
			}
		}
	}
}
