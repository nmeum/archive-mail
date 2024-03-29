package main

import (
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var chkSum hash.Hash = sha256.New()

type Mail struct {
	maildir   string // path to maildir
	directory string // new, cur, or tmp
	name      string // basename of message
}

func NewMail(maildir string, fp string) (*Mail, error) {
	cdir := filepath.Clean(maildir)
	cmsg := filepath.Clean(fp)
	if !strings.HasPrefix(cmsg, cdir) {
		return nil, errors.New("mail is not in given maildir")
	}

	return &Mail{
		maildir:   maildir,
		directory: getDir(fp),
		name:      filepath.Base(fp),
	}, nil
}

func (m *Mail) String() string {
	return filepath.Join(filepath.Base(m.maildir), m.directory, m.name)
}

func (m *Mail) Path() string {
	return filepath.Join(m.maildir, m.directory, m.name)
}

func (m *Mail) Checksum() (string, error) {
	data, err := os.ReadFile(m.Path())
	if err != nil {
		return "", err
	}

	return string(chkSum.Sum(data)), nil
}

func (m *Mail) IsSame(other *Mail) bool {
	return filepath.Base(m.maildir) == filepath.Base(other.maildir) &&
		m.directory == other.directory &&
		m.name == other.name
}

// TODO: fsync
func (m *Mail) CopyTo(maildir string) error {
	file, err := os.Open(m.Path())
	if err != nil {
		return err
	}
	defer file.Close()

	tmpFp := filepath.Join(maildir, "tmp", m.name)
	newFile, err := os.Create(tmpFp)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}

	newFp := filepath.Join(maildir, m.directory, m.name)
	return os.Rename(tmpFp, newFp)
}
