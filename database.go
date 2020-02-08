package main

import (
	"crypto/sha1"
	"errors"
	"hash"
	"io/ioutil"
)

// SHA1 should be good enough for this purpose
var chkSum hash.Hash = sha1.New()

type MailPair struct {
	old, new *Mail
}

type MailDatabase struct {
	newMsgs []*Mail
	modMsgs []*MailPair

	oldMsgs map[string]*Mail
}

func NewMailDatabase() *MailDatabase {
	db := new(MailDatabase)
	db.oldMsgs = make(map[string]*Mail)
	return db
}

func (db *MailDatabase) AddOldMessage(mail *Mail) error {
	data, err := ioutil.ReadFile(mail.Path())
	if err != nil {
		return err
	}

	sum := string(chkSum.Sum(data))
	if _, ok := db.oldMsgs[sum]; ok {
		return errors.New("hash collision")
	}

	db.oldMsgs[sum] = mail
	return nil
}

func (db *MailDatabase) GetOldMessage(mail *Mail) (*Mail, error) {
	data, err := ioutil.ReadFile(mail.Path())
	if err != nil {
		return nil, err
	}
	sum := string(chkSum.Sum(data))

	mail, ok := db.oldMsgs[sum]
	if !ok {
		return nil, nil
	}

	return mail, nil
}

func (db *MailDatabase) AddNewMessage(old *Mail, new *Mail) {
	if old == nil {
		db.newMsgs = append(db.newMsgs, new)
	} else {
		db.modMsgs = append(db.modMsgs, &MailPair{old, new})
	}
}
