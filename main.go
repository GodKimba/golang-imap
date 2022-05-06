package main

import (
	"github.com/emersion/go-imap/client"
	"log"
)

//"github.com/emersion/go-imap"
// https://github.com/emersion/go-imap/wiki/Fetching-messages

const userMailAccount = "user@mail.acc"
const mailServer = "imap.gmail.com:993"
const userMailPassword = "password"

var c *client.Client

func connectToClientServer() {
	c, err := client.DialTLS(mailServer, nil)
	if err != nil {
		log.Fatalf("Couldn't login, %v", err)
	}

	log.Println("Connected")

	if err := c.Login(userMailAccount, userMailPassword); err != nil {
		log.Fatalf("Couldn't login, %v", err)
	}
	log.Println("Logged in!")
}

func selectMailBox() {
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for inbox: ", mbox.Flags)

}

func main() {
	connectToClientServer()
}
