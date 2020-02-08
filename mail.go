package main

import (
	"errors"
	"path/filepath"
	"strings"
)

type Mail struct {
	maildir   string // absolute path to maildir
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

func (m *Mail) Path() string {
	return filepath.Join(m.maildir, m.directory, m.name)
}

func (m *Mail) IsSame(other *Mail) bool {
	return filepath.Base(m.maildir) == filepath.Base(other.maildir) &&
		m.directory == other.directory &&
		m.name == other.name
}
