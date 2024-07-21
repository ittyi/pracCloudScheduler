package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
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

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("x-mm-message-id", "test")
	m.SetHeader("Subject", "Test email")
	m.SetBody("text/plain", "This is a test email sent from Go using gomail package.")

	var portNum int
	portNum, _ = strconv.Atoi(port)
	d := gomail.NewDialer(hostname, portNum, from, password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("Email sent successfully!")
}
