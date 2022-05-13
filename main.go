package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

//"github.com/emersion/go-imap"
// https://github.com/emersion/go-imap/wiki/Fetching-messages

// Function to get the environment variable from .evn
func getEnvKey(key string) string {
	err := godotenv.Load()
	
	if err != nil {
		log.Fatalf("Error  loading .env file")
	}
	
	return os.Getenv(key)
}

var userMailAccount = getEnvKey("USERNAME")
var userMailPassword = getEnvKey("PASSWORD")
const mailServer = "imap.gmail.com:993"
var c *client.Client


func connectToClientServer(c *client.Client, err error) {
	c, err = client.DialTLS(mailServer, nil)
	if err != nil {
		log.Fatalf("Couldn't login, %v", err)
	}

	
	log.Println("Connected")
	defer c.Logout()
	
	if err := c.Login(userMailAccount, userMailPassword); err != nil {
		log.Fatalf("Couldn't login, %v", err)
	}
	log.Println("Logged in!")

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for inbox: ", mbox.Flags)
	
	// Searching criteria almost working, lines below are responsible for it
	keyword := "rafael"
	criteria := imap.NewSearchCriteria()
	criteria.Header.Get(keyword)
	c.Search(criteria)
	fmt.Println(criteria)
	
}

//func selectMailBox(c *client.Client, err error) {
//}

func main() {
	fmt.Println(userMailAccount, userMailPassword)
	connectToClientServer(c, nil)


	
}

	
