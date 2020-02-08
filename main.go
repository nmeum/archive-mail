package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type MailWalkFn func(mail *Mail, db *MailDatabase, err error) error

func indexOldMsgs(mail *Mail, db *MailDatabase, err error) error {
	if err != nil {
		panic(err)
	}

	db.AddOldMessage(mail)
	return nil
}

func indexNewMsgs(mail *Mail, db *MailDatabase, err error) error {
	if err != nil {
		panic(err)
	}

	oldMail, err := db.GetOldMessage(mail)
	if err != nil {
		return err
	}

	if oldMail == nil {
		db.AddNewMessage(nil, mail)
	} else if !oldMail.IsSame(mail) {
		db.AddNewMessage(oldMail, mail)
	}

	return nil
}

func walkMaildir(maildir string, db *MailDatabase, walkFn MailWalkFn) error {
	wrapFn := func(path string, info os.FileInfo, err error) error {
		handleError := func(err error) error { return walkFn(nil, nil, err) }
		if err != nil {
			return handleError(err)
		}

		if info.IsDir() {
			if !isMaildir(info.Name()) {
				return handleError(fmt.Errorf("unexpected folder %q", info.Name()))
			} else {
				return nil
			}
		}

		mail, err := NewMail(maildir, path)
		if err != nil {
			return handleError(err)
		}

		return walkFn(mail, db, err)
	}

	for _, dir := range []string{"cur", "new", "tmp"} {
		err := filepath.Walk(filepath.Join(maildir, dir), wrapFn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns mapping new maildir → old maildir.
func parseArgs(args []string) (map[string]string, error) {
	parsedArgs := make(map[string]string)
	for _, arg := range args {
		splitted := strings.Split(arg, "→")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("invalid argument %q", arg)
		}

		old := splitted[0]
		new := splitted[1]

		if _, ok := parsedArgs[new]; ok {
			return nil, fmt.Errorf("duplicate maildir %q", arg)
		}
		parsedArgs[new] = old
	}
	return parsedArgs, nil
}

func indexMsgs() (*MailDatabase, error) {
	var wg sync.WaitGroup
	db := NewMailDatabase()
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		return nil, err
	}

	wfn := func(dir string, mfn MailWalkFn) {
		defer wg.Done()
		err := walkMaildir(dir, db, mfn)
		if err != nil {
			log.Fatal(err)
		}
	}

	wg.Add(len(args))
	for _, old := range args {
		go wfn(old, indexOldMsgs)
	}
	wg.Wait()

	wg.Add(len(args))
	for new, _ := range args {
		go wfn(new, indexNewMsgs)
	}
	wg.Wait()

	return db, nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s OLD_MAILDIR→NEW_MAILDIR ...\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	db, err := indexMsgs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("##\n# New Messages\n##\n\n")
	for _, new := range db.newMsgs {
		fmt.Println(new)
	}
	fmt.Printf("\n##\n# Changed Messages\n##\n\n")
	for _, msg := range db.modMsgs {
		fmt.Printf("%s → %s\n", msg.old, msg.new)
	}
}
