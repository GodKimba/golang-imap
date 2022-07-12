package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

func getEnvKey(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error  loading .env file")
	}

	return os.Getenv(key)
}

var userMailAccount string
var userMailPassword string
var deletionType string
var deletionSpecify string
var c *client.Client

var deleteBySubject = "SUBJECT"
var deleteBySender = "FROM"

const mailServer = "imap.gmail.com:993"

type User struct {
	c *client.Client
}

func NewClient() *User {
	return &User{c: c}
}

func (u *User) createEnvFile() {
	var username string
	var password string

	f, err := os.Create(".env")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	fmt.Println("Enter your Mail Address: ")
	fmt.Scanln(&username)
	fmt.Println("Enter your Mail App Password")
	fmt.Scanln(&password)
	_, err2 := f.WriteString("USERNAME:" + username + "\nPASSWORD: " + password)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("File created")
}

func (u *User) checkIfEnvFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (u *User) chooseDeletionType(err error) (string, string) {
	fmt.Println("Do you want to delete by subject(s) or sender(sd)?")
	fmt.Scanln(&deletionType)
	if deletionType == "s" {
		fmt.Println("Enter the keyword:")
		fmt.Scanln(&deletionSpecify)

		deletionType = deleteBySubject

	} else if deletionType == "sd" {
		fmt.Println("Enter the sender mail:")
		fmt.Scanln(&deletionSpecify)

		deletionType = deleteBySender

	} else {
		log.Println("Try again.")
		return u.chooseDeletionType(nil)
	}
	return deletionType, deletionSpecify
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

func (u *User) searchingCriteria(deletionType, deletionSpecify string, err error) []uint32 {
	criteria := imap.NewSearchCriteria()
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
	}
}
func (u *User) flagAndDelete(ids []uint32, err error) {
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}
		if err := u.c.Store(seqset, item, flags, nil); err != nil {
			log.Fatal(err)
		}

		if err := u.c.Expunge(nil); err != nil {
			log.Fatal(err)
		}
		log.Println("Last message has been deleted")
	}
}

func main() {
	u := NewClient()
	if u.checkIfEnvFileExists(".env") {
		fmt.Println("User credentials selected")

	} else {
		fmt.Println("Gmail requires you to use a specific password for apps.\nYou only need to create one time, follow instructions in this link:\nhttps://support.google.com/accounts/answer/185833?hl=en")
		u.createEnvFile()

	}

	userMailAccount = getEnvKey("USERNAME")
	userMailPassword = getEnvKey("PASSWORD")

	u.conectToMailServer(nil)
	u.loginToMailServer(nil)
	defer u.c.Logout()

	u.selectMailBox(nil)
	u.chooseDeletionType(nil)
	ids := u.searchingCriteria(deletionType, deletionSpecify, nil)
	u.showMessages(ids, nil)
	u.flagAndDelete(ids, nil)
}
//test
