package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

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
	// c, err := client.DialTLS("imap.gmail.com:993", nil)
	// c, err := client.DialTLS(
	// 	"imap-mail.outlook.com:143",
	// 	nil,
	// )
	c, err := client.Dial("imap-mail.outlook.com:143")
	if err != nil {
		log.Fatal(err)
	}

	// STARTTLS を使用して接続をアップグレード
	if err := c.StartTLS(&tls.Config{}); err != nil {
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
	criteria.Header.Add("Message-ID", "") // 例: <unique-message-id@example.com>

	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	msgs := make(chan *imap.Message, len(uids))
	if err := c.Fetch(seqSet, items, msgs); err != nil {
		log.Fatal(err)
	}
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

		// references, err := header.Text("References")
		// if err != nil {
		// 	log.Println("References ヘッダーが存在しません")
		// } else {
		// 	fmt.Println("References:", references)
		// }

		// inReplyTo, err := header.Text("In-Reply-To")
		// if err != nil {
		// 	log.Println("In-Reply-To ヘッダーが存在しません")
		// } else {
		// 	fmt.Println("In-Reply-To:", inReplyTo)
		// }

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
				// 添付ファイル名取得
				filename, _ := h.Filename()
				fmt.Println("Found attachment:", filename)
				fmt.Println("part.Header:", part.Header)

				// 添付ファイルのContent-Lengthを取得
				contentLengthStr := part.Header.Get("Content-Length")
				var contentLength int64
				if contentLengthStr != "" {
					var err error
					contentLength, err = strconv.ParseInt(contentLengthStr, 10, 64)
					if err != nil {
						log.Fatalf("error parsing Content-Length: %v", err)
					}
				}
				fmt.Println("contentLength:", contentLength)

				// 添付ファイルの内容を取得してデコード
				// Base64エンコードされたデータを取得
				attachmentData, err := io.ReadAll(part.Body)
				if err != nil {
					log.Fatal(err)
				}

				// サイズが一致するか確認
				fmt.Println("サイズが一致するか確認:", contentLength > 0 && int64(len(attachmentData)) != contentLength)
				if contentLength > 0 && int64(len(attachmentData)) != contentLength {
					log.Fatalf("file size mismatch: expected %d bytes, got %d bytes", contentLength, len(attachmentData))
				}

				err = os.WriteFile(filename, attachmentData, 0644)
				if err != nil {
					log.Fatalf("error writing file: %v", err)
				}

				fmt.Printf("添付ファイルを保存しました: %s\n", filename)
			}
		}
	}
}
