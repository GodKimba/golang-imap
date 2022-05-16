package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
	"ioutil"
	"log"
	"os"
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

	

	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(mbox.Messages, mbox.Messages)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	log.Println("Last message:")
	msg := <-messages
	r := msg.GetBody(section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}

	header := m.Header
	log.Println("Date:", header.Get("Date"))
	log.Println("From:", header.Get("From"))
	log.Println("To:", header.Get("To"))
	log.Println("Subject:", header.Get("Subject"))

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(body)
}

// Expunge func
// item := imap.FormatFlagsOp(imap.AddFlags, true)
// flags := []interface{}{imap.DeletedFlag}
// if err := c.Store(seqset, item, flags, nil); err != nil {
// 	log.Fatal(err)
// }

// // Then delete it
// if err := c.Expunge(nil); err != nil {
// 	log.Fatal(err)
// }

// log.Println("Last message has been deleted")
