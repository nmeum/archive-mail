package main

import (
	"errors"
	"sync"
)

type MailPair struct {
	old, new *Mail
}

type MailDatabase struct {
	msgMtx *sync.Mutex

	newMsgs []*Mail
	modMsgs []*MailPair
	oldMsgs map[string]*Mail
}

func NewMailDatabase() *MailDatabase {
	db := new(MailDatabase)
	db.msgMtx = new(sync.Mutex)
	db.oldMsgs = make(map[string]*Mail)
	return db
}

func (db *MailDatabase) AddOldMessage(mail *Mail) error {
	csum, err := mail.Checksum()
	if err != nil {
		return err
	}

	db.msgMtx.Lock()
	defer db.msgMtx.Unlock()

	if _, ok := db.oldMsgs[csum]; ok {
		return errors.New("hash collision")
	}
	db.oldMsgs[csum] = mail

	return nil
}

func (db *MailDatabase) GetOldMessage(mail *Mail) (*Mail, error) {
	csum, err := mail.Checksum()
	if err != nil {
		return nil, err
	}

	db.msgMtx.Lock()
	defer db.msgMtx.Unlock()

	mail, ok := db.oldMsgs[csum]
	if !ok {
		return nil, nil
	}

	return mail, nil
}

func (db *MailDatabase) AddNewMessage(old *Mail, new *Mail) {
	db.msgMtx.Lock()
	if old == nil {
		db.newMsgs = append(db.newMsgs, new)
	} else {
		db.modMsgs = append(db.modMsgs, &MailPair{old, new})
	}
	db.msgMtx.Unlock()
}
