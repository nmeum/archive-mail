package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		// TODO: delete old mail from database?
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

func indexMsgs(olddir, newdir string) (*MailDatabase, error) {
	db := NewMailDatabase()
	err := walkMaildir(olddir, db, indexOldMsgs)
	if err != nil {
		return nil, err
	}
	err = walkMaildir(newdir, db, indexNewMsgs)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func mergeMsgs(olddir, newdir string) error {
	db, err := indexMsgs(olddir, newdir)
	if err != nil {
		return err
	}

	// TODO: merge them
	fmt.Printf("##\n# New Messages\n##\n\n")
	for _, new := range db.newMsgs {
		fmt.Println(new)
	}
	fmt.Printf("\n##\n# Changed Messages\n##\n\n")
	for _, msg := range db.modMsgs {
		fmt.Printf("%s → %s\n", msg.old, msg.new)
	}

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) <= 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s OLD_MAILDIR NEW_MAILDIR\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	// TODO: Handle moves between different maildirs

	mergeMsgs(os.Args[1], os.Args[2])
}
