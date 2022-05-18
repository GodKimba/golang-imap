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

// getting environment variables
func getEnvKey(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error  loading .env file")
	}

	return os.Getenv(key)
}

var userMailAccount = getEnvKey("USERNAME")
var userMailPassword = getEnvKey("PASSWORD")
var deletionType string
var deletionSpecify string
var c *client.Client

const mailServer = "imap.gmail.com:993"

type User struct {
	c *client.Client
}

func NewClient() *User {
	return &User{c: c}
}

func (u *User) conectToMailServer(err error) {
	u.c, err = client.DialTLS(mailServer, nil)
	if err != nil {
		log.Fatalf("Couldn't login, %v", err)
	}
	log.Println("Connected")
}

func (u *User) loginToMailServer(err error) {
	if err := u.c.Login(userMailAccount, userMailPassword); err != nil {
		log.Fatalf("Couldn't login, %v", err)
	}
	log.Println("Logged in!")
}

func (u *User) selectMailBox(err error) {
	mbox, err := u.c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for inbox:", mbox.Flags)
}

func (u *User) searchingCriteria(err error) []uint32 {
	criteria := imap.NewSearchCriteria()
	// Requesting user input
	fmt.Println("Do you want to delete by subject or sender?") // This line is responsable for subject selection
	fmt.Scanln(&deletionType)
	fmt.Println("Enter the keyword/mail address")
	fmt.Scanln(&deletionSpecify)

	criteria.Header.Add(deletionType, deletionSpecify)
	ids, err := u.c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("IDs found:", ids)
	return ids
}

func (u *User) showMessages(ids []uint32, err error) {
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- u.c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()
		log.Println("To be deleted messages:")
		// change here to apply deletion
		for msg := range messages {
			log.Println("* " + msg.Envelope.Subject)
		}
		if err != nil {
			log.Fatal(err)
		}
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}
		if err := u.c.Store(seqset, item, flags, nil); err != nil {
			log.Fatal(err)
		}
		// Then delete it
		if err := u.c.Expunge(nil); err != nil {
			log.Fatal(err)
		}
		log.Println("Last message has been deleted")
	}
}

func main() {
	u := NewClient()
	u.conectToMailServer(nil)
	u.loginToMailServer(nil)
	u.selectMailBox(nil)
	ids := u.searchingCriteria(nil)
	u.showMessages(ids, nil)
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
