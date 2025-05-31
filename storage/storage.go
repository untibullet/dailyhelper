package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/untibullet/dailyhelper/tools/elog"
)

var (
	ErrNoSavedPages = errors.New("no saved pages for user")
)

type PageStorer interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(*Page) error
	Exists(*Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := md5.New()

	if _, err := h.Write([]byte(p.URL)); err != nil {
		return "", elog.Wrap("can`t calculate hash for URL", err)
	}

	if _, err := h.Write([]byte(p.UserName)); err != nil {
		return "", elog.Wrap("can`t calculate hash for user name", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
