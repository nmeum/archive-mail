package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	verbose = flag.Bool("v", false, "print out performed changes")
	dryrun  = flag.Bool("d", false, "don't perform any changes, combine with -v")
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
			if !isMaildirFn(info.Name()) {
				return handleError(fmt.Errorf("unexpected directory %q", info.Name()))
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

	for _, dir := range []string{"cur", "new"} {
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

		new := splitted[0]
		old := splitted[1]
		for _, dir := range []string{old, new} {
			if !isValidMaildir(dir) {
				return nil, fmt.Errorf("%q is not a valid maildir", dir)
			}
		}

		if _, ok := parsedArgs[new]; ok {
			return nil, fmt.Errorf("duplicate maildir %q", arg)
		}
		parsedArgs[new] = old
	}
	return parsedArgs, nil
}

func indexMsgs(args map[string]string) (*MailDatabase, error) {
	var wg sync.WaitGroup
	db := NewMailDatabase()

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

func archiveMsgs(args map[string]string, db *MailDatabase) error {
	for _, new := range db.newMsgs {
		if *verbose {
			fmt.Printf("new: %s\n", new)
		}
		if *dryrun {
			continue
		}

		err := new.CopyTo(args[new.maildir])
		if err != nil {
			return err
		}
	}
	for _, pair := range db.modMsgs {
		if *verbose {
			fmt.Printf("move: %s → %s\n", pair.old, pair.new)
		}
		if *dryrun {
			continue
		}

		destDir := args[pair.new.maildir]
		newFp := filepath.Join(destDir, pair.new.directory, pair.new.name)
		err := os.Rename(pair.old.Path(), newFp)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [-v] [-d] MAILDIR_CURRENT→MAILDIR_ARCHIVE...\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	args, err := parseArgs(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	db, err := indexMsgs(args)
	if err != nil {
		log.Fatal(err)
	}
	err = archiveMsgs(args, db)
	if err != nil {
		log.Fatal(err)
	}
}
