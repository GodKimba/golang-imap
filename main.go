package main

import (
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
	criteria := imap.NewSearchCriteria()
	criteria.Header.Add("SUBJECT", "Rafael")
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("IDs found:", ids)

	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		log.Println("Unseen messages:")
		for msg := range messages {
			log.Println("* " + msg.Envelope.Subject)
		}

		if err := <-done; err != nil {
			log.Fatal(err)
		}
	}
}
func main() {
	connectToClientServer(c, nil)
}

//func selectMailBox(c *client.Client, err error) {
//}

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
