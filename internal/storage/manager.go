package storage

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Application data directory: https://hub.jmonkeyengine.org/t/appdata-equivalent-on-macos-and-linux/43735

const (
	appName = "mtgstudio"
	sets    = "sets"
	cards   = "cards"
	symbols = "symbols"
)

// Manager handles the file access for the application
type Manager struct {
	applicationDirectory string
	setsDirectory        string
	cardsDirectory       string
	symbolsDirectory     string
	dbPath               string
	db                   *gorm.DB

	*gormSetRepo
	*GormCardRepo
}

func NewManager() *Manager {
	m := &Manager{applicationDirectory: getApplicationDirectory()}
	err := os.Mkdir(m.applicationDirectory, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	m.setsDirectory = filepath.Join(m.applicationDirectory, sets)
	err = os.Mkdir(m.setsDirectory, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	m.cardsDirectory = filepath.Join(m.applicationDirectory, cards)
	err = os.Mkdir(m.cardsDirectory, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	m.symbolsDirectory = filepath.Join(m.applicationDirectory, symbols)
	err = os.Mkdir(m.symbolsDirectory, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	m.dbPath = filepath.Join(m.applicationDirectory, "mtgstudio.db")
	m.db, err = gorm.Open(sqlite.Open(m.dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m.gormSetRepo = newGormSetRepo(m.db)
	m.GormCardRepo = newGormCardRepo(m.db)
	return m
}

func (m *Manager) LoadSetIcon(uri string) ([]byte, error) {

	res := path.Base(uri)
	reconfiguredName := m.reconfigureName(res)
	resourcePath := filepath.Join(m.setsDirectory, reconfiguredName)
	data, err := os.ReadFile(resourcePath)
	if err == nil {
		return data, nil
	}

	data, err = requestURL(uri)
	if err != nil {
		return nil, err
	}

	os.WriteFile(resourcePath, data, os.ModePerm)

	return data, nil
}

func (m *Manager) LoadCardImage(uri string) ([]byte, error) {
	res := path.Base(uri)
	reconfiguredName := m.reconfigureName(res)
	resourcePath := filepath.Join(m.cardsDirectory, reconfiguredName)
	data, err := os.ReadFile(resourcePath)
	if err == nil {
		return data, nil
	}

	data, err = requestURL(uri)
	if err != nil {
		return nil, err
	}

	os.WriteFile(resourcePath, data, os.ModePerm)

	return data, nil
}

func (m *Manager) LoadSymbolImage(uri string) ([]byte, error) {
	res := path.Base(uri)
	reconfiguredName := m.reconfigureName(res)
	resourcePath := filepath.Join(m.symbolsDirectory, reconfiguredName)
	data, err := os.ReadFile(resourcePath)
	if err == nil {
		return data, nil
	}

	data, err = requestURL(uri)
	if err != nil {
		return nil, err
	}

	os.WriteFile(resourcePath, data, os.ModePerm)

	return data, nil
}

func (m *Manager) reconfigureName(name string) string {
	if !strings.Contains(name, "?") {
		return name
	}

	idx := strings.Index(name, "?")
	ext := path.Ext(name[:idx])
	extIdx := strings.Index(name[:idx], ext)
	return name[:extIdx] + "_" + name[idx+1:] + ext
}
