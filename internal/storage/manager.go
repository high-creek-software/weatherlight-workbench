package storage

import (
	"errors"
	"gitlab.com/high-creek-software/goscryfall"
	scryfallcards "gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/high-creek-software/goscryfall/rulings"
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
	userDBPath           string
	userDB               *gorm.DB
	client               *goscryfall.Client

	*gormSetRepo
	*gormCardRepo
	*bookmarkRepo
}

func NewManager(client *goscryfall.Client) *Manager {
	m := &Manager{applicationDirectory: getApplicationDirectory(), client: client}
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
	m.userDBPath = filepath.Join(m.applicationDirectory, "userdata.db")
	m.userDB, err = gorm.Open(sqlite.Open(m.userDBPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m.gormSetRepo = newGormSetRepo(m.db)
	m.gormCardRepo = newGormCardRepo(m.db)
	m.bookmarkRepo = NewBookmarkRepo(m.userDB)
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

func (m *Manager) CardSearch(sr SearchRequest) ([]scryfallcards.Card, error) {
	return m.gormCardRepo.Search(sr)
}

func (m *Manager) ListBookmarked() ([]scryfallcards.Card, error) {
	bs, err := m.bookmarkRepo.List()
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, b := range bs {
		ids = append(ids, b.CardID)
	}

	return m.gormCardRepo.ListByIds(ids)
}

func (m *Manager) LoadRulings(c *scryfallcards.Card) ([]rulings.Ruling, error) {
	return m.client.List(c.RulingsUri)
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
