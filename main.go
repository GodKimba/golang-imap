package main

import (
      "log"
      "github.com/emersion/go-imap/client"
)
      //"github.com/emersion/go-imap"

const userMailAccount = "user@mail.acc"
const mailServer = "imap.gmail.com:993"
const userMailPassword = "password"

func connectToClientServer() {
  c, err := client.DialTLS(mailServer, nil)
  if err != nil {
    log.Fatalf("Couldn't login, %v", err)
  }
  log.Println("Connected")

  defer c.Logout()

  if err := c.Login(userMailAccount, userMailPassword); err != nil {
    log.Fatalf("Couldn't login, %v", err)
  }
  log.Println("Logged in!")
}

func selectMailBox() {
  connectToClientServer()
  mbox, err := c.selectMailBox("INBOX", false)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Flags for inbox: ", mbox.Flags)

}

func main() {
  connectToClientServer()
}
