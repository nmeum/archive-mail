package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ModMsg struct {
	old, new string
}

type MailInfo struct {
	checksum string
	os.FileInfo
}

type mailWalkFn func(path string, info *MailInfo, err error) error

// SHA1 should be good enough for this purpose
var chkSum hash.Hash = sha1.New()

var (
	oldMsgs = make(map[string]string)
	newMsgs = []string{}
	modMsgs = []ModMsg{}
)

func isMaildir(name string) bool {
	return name == "new" || name == "cur" || name == "tmp"
}

func getDir(path string) string {
	dir := filepath.Base(filepath.Dir(path))
	if !isMaildir(dir) {
		panic("unexpected non-maildir folder")
	}

	return dir
}

func indexOldMsgs(path string, info *MailInfo, err error) error {
	if err != nil {
		panic(err)
	}

	oldMsgs[info.checksum] = path
	return nil
}

func indexNewMsgs(path string, info *MailInfo, err error) error {
	if err != nil {
		panic(err)
	}

	old, ok := oldMsgs[info.checksum]
	if ok {
		newDir := getDir(path)
		if getDir(old) == newDir && filepath.Base(old) == info.Name() {
			goto cont
		}

		newPath := filepath.Join(filepath.Dir(old), "..", newDir, info.Name())
		modMsgs = append(modMsgs, ModMsg{old, filepath.Clean(newPath)})
	} else {
		newMsgs = append(newMsgs, path)
	}

cont:
	delete(oldMsgs, info.checksum)
	return nil
}

func walkMaildir(maildir string, walkFn mailWalkFn) error {
	wrapFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return walkFn(path, nil, err)
		}

		if info.IsDir() {
			if !isMaildir(info.Name()) {
				return fmt.Errorf("unexpected folder %q", info.Name())
			} else {
				return nil
			}
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		minfo := MailInfo{string(chkSum.Sum(data)), info}
		return walkFn(path, &minfo, err)
	}

	for _, dir := range []string{"cur", "new", "tmp"} {
		err := filepath.Walk(filepath.Join(maildir, dir), wrapFn)
		if err != nil {
			return err
		}
	}

	return nil
}

func indexMsgs(olddir, newdir string) error {
	err := walkMaildir(olddir, indexOldMsgs)
	if err != nil {
		return err
	}

	err = walkMaildir(newdir, indexNewMsgs)
	if err != nil {
		return err
	}

	return nil
}

func mergeMsgs(olddir, newdir string) error {
	err := indexMsgs(olddir, newdir)
	if err != nil {
		return err
	}

	// TODO: merge them
	fmt.Printf("##\n# New Messages\n##\n\n")
	for _, new := range newMsgs {
		fmt.Println(new)
	}
	fmt.Printf("\n##\n# Changed Messages\n##\n\n")
	for _, msg := range modMsgs {
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
