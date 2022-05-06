package main

import (
	"github.com/emersion/go-imap/client"
	"log"
	"github.com/joho/godotenv"
	"os"
)

//"github.com/emersion/go-imap"
// https://github.com/emersion/go-imap/wiki/Fetching-messages

// Function to get the environment variable from .evn
func getEnvKey(key string) string {
	err := godoenv.Load()
	
	if err != nil {
		log.Fatalf("Error  loading .env file")
	}
	
	return os.Getenv(key)
}

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
