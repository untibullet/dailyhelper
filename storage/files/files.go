package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/untibullet/dailyhelper/storage"
	"github.com/untibullet/dailyhelper/tools/elog"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func NewStrorage(basePath string) (*Storage, error) {
	info, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("base path '%s' does not exist", basePath)
			return nil, elog.Wrap(msg, err)
		}
		msg := fmt.Sprintf("could not stat base path '%s'", basePath)
		return nil, elog.Wrap(msg, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("base path '%s' is not a directory", basePath)
	}

	return &Storage{
		basePath: basePath,
	}, nil
}

func (s *Storage) Save(page *storage.Page) (err error) {
	userDirPath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(userDirPath, defaultPerm); err != nil {
		return elog.Wrap(fmt.Sprintf("cannot create user directory '%s'", userDirPath), err)
	}

	pagePath, err := s.getPagePath(page)
	if err != nil {
		return elog.Wrap("cannot generate page path for save", err)
	}

	file, err := os.Create(pagePath)
	if err != nil {
		return elog.Wrap(fmt.Sprintf("cannot create page file '%s'", pagePath), err)
	}
	defer func() {
		_ = file.Close()
	}()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		msg := fmt.Sprintf("cannot encode page data to file '%s'", pagePath)
		return elog.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	userDirPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(userDirPath)
	if err != nil {
		msg := fmt.Sprintf("cannot read user directory '%s'", userDirPath)
		return nil, elog.Wrap(msg, err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(userDirPath, file.Name()))
}

func (s *Storage) Remove(page *storage.Page) error {
	pagePath, err := s.getPagePath(page)
	if err != nil {
		return elog.Wrap("cannot generate page path for removal", err)
	}

	if err := os.Remove(pagePath); err != nil {
		msg := fmt.Sprintf("can`t remove page for path: %s", pagePath)

		return elog.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) Exists(page *storage.Page) (bool, error) {
	pagePath, err := s.getPagePath(page)
	if err != nil {
		return false, elog.Wrap("cannot generate page path for existence check", err)
	}

	_, err = os.Stat(pagePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, elog.Wrap(fmt.Sprintf("error checking page existence for '%s'", pagePath), err)
	}

	return true, nil
}

func (s *Storage) getPagePath(page *storage.Page) (string, error) {
	fileNameHash, err := page.Hash()
	if err != nil {
		return "", elog.Wrap("cannot calculate page hash for path", err)
	}
	return filepath.Join(s.basePath, page.UserName, fileNameHash), nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, elog.Wrap("can`t open file with page", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var p storage.Page
	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, elog.Wrap("can`t decode page", err)
	}

	return &p, nil
}
