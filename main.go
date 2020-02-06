package main

import (
	"flag"
	"io/ioutil"
	"os"
	"log"
	"hash"
	"strings"
	"fmt"
	"path/filepath"
	"crypto/sha1"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
)

var (
	csum = flag.String("c", "sha1", "checksum algorithm to use for duplicate detection")
)

type ModMsg struct {
	path string
	newFlags string
}

var chkSum hash.Hash

var (
	oldMsgs = make(map[string]string)
	newMsgs = []string{}
	modMsgs = []ModMsg{}
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"USAGE: %s [FLAGS] OLDMAILDIR NEWMAILDIR\n\n"+
			"The following flags are supported:\n\n", os.Args[0])

	flag.PrintDefaults()
	os.Exit(2)
}

func strToHsh(algorithm string) *hash.Hash {
	var hash hash.Hash
	switch (strings.ToLower(algorithm)) {
	case "md5":
		hash = md5.New()
	case "sha1":
		hash = sha1.New()
	case "sha256":
		hash = sha256.New()
	case "sha512":
		hash = sha512.New()
	}
	return &hash
}

func isMaildir(name string) bool {
	return name == "new" || name == "cur" || name == "tmp"
}

func extractFlags(fn string) (string, error) {
	idx := strings.IndexByte(fn, ':')
	if idx == -1 || idx + 1 >= len(fn) {
		return "", fmt.Errorf("message has no flags")
	}

	return fn[idx:], nil
}

func indexOldMsgs(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	sum := string(chkSum.Sum(data))
	oldMsgs[sum] = path

	return nil
}

func indexNewMsgs(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	sum := string(chkSum.Sum(data))
	old, ok := oldMsgs[sum]
	if ok {
		if filepath.Base(old) == info.Name() {
			goto cont
		}
		newFlags, err := extractFlags(info.Name())
		if err != nil {
			goto cont
		}

		modMsgs = append(modMsgs, ModMsg{old, newFlags})
	} else {
		newMsgs = append(newMsgs, path)
	}

cont:
	delete(oldMsgs, sum)
	return nil
}

func walkMaildir(maildir string, walkFn filepath.WalkFunc) error {
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

		return walkFn(path, info, err)
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
	fmt.Printf("\n##\n# New Messages\n##\n\n")
	for _, new := range newMsgs {
		fmt.Println(new)
	}
	fmt.Printf("##\n# Changed Flags\n##\n\n")
	for _, msg := range modMsgs {
		fmt.Printf("%s, new flags: %s\n", msg.path, msg.newFlags)
	}

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Parse()
	if flag.NArg() < 2 {
		usage()
	}

	sum := strToHsh(*csum)
	if sum == nil {
		log.Fatalf("Unsupported checksum algorithm %q\n", *csum)
	} else {
		chkSum = *sum
	}

	mergeMsgs(flag.Arg(0), flag.Arg(1))
}
