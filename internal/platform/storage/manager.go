package storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/high-creek-software/goscryfall"
	scryfallcards "github.com/high-creek-software/goscryfall/cards"
	"github.com/high-creek-software/goscryfall/decks"
	"github.com/high-creek-software/goscryfall/rulings"
	"github.com/rs/xid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Application data directory: https://hub.jmonkeyengine.org/t/appdata-equivalent-on-macos-and-linux/43735

const (
	appName = "weatherlightworkbench"
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
	*gormDeckRepo
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
	m.dbPath = filepath.Join(m.applicationDirectory, "weatherlightworkbench.db")
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
	m.gormDeckRepo = newGormDeckRepo(m.userDB)
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
	if strings.Contains(uri, "front") {
		reconfiguredName = "front_" + reconfiguredName //fmt.Sprint("front_%s", reconfiguredName)
	} else if strings.Contains(uri, "back") {
		reconfiguredName = "back_" + reconfiguredName //fmt.Sprintf("back_%s", reconfiguredName)
	}
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

func (m *Manager) ParseDeckDefinition(deckName, data string) (*decks.Deck, error) {
	buf := bytes.NewBufferString(data)
	return decks.Unmarshal(deckName, buf)
}

func (m *Manager) CreateDeck(name, deckType string) (*Deck, error) {
	deckID := xid.New().String()
	gd := gormDeck{ID: deckID, Name: name, CreateAt: time.Now(), CoverImage: "", DeckType: deckType}
	err := m.gormDeckRepo.create(gd)
	if err != nil {
		return nil, err
	}

	return &Deck{ID: deckID, Name: name, CreatedAt: gd.CreateAt, DeckType: deckType}, nil
}

func (m *Manager) AddCardToDeck(deck Deck, cardID, cardName string, count int, association AssociationType) error {
	gdc := gormDeckCard{ID: xid.New().String(), DeckID: deck.ID, CardID: cardID, CardName: cardName, AssociationType: association, Count: count}

	return m.gormDeckRepo.addCard(gdc)
}

func (m *Manager) LoadDeck(id string) (Deck, error) {
	d, err := m.gormDeckRepo.findDeck(id)
	if err != nil {
		return d, err
	}

	allCards, err := m.gormDeckRepo.listDeckCards(id)
	if err != nil {
		return d, err
	}

	for _, dc := range allCards {
		if cCard, err := m.gormCardRepo.findById(dc.CardID); err == nil {
			switch dc.AssociationType {
			case AssociationSideboard:
				d.Sideboard = append(d.Sideboard, DeckCard{ID: dc.ID, Count: dc.Count, Card: cCard, AssociationType: dc.AssociationType})
			case AssociationMain:
				d.Main = append(d.Main, DeckCard{ID: dc.ID, Count: dc.Count, Card: cCard, AssociationType: dc.AssociationType})
			case AssociationCommander:
				d.Commander = &DeckCard{ID: dc.ID, Count: dc.Count, Card: cCard, AssociationType: dc.AssociationType}
			}
		}
	}

	return d, nil
}

func (m *Manager) CopyDeck(d Deck, copyName string) (Deck, error) {
	deckID := xid.New().String()
	gd := gormDeck{ID: deckID, Name: copyName, CreateAt: time.Now(), CoverImage: "", DeckType: d.DeckType}
	err := m.gormDeckRepo.create(gd)
	if err != nil {
		return Deck{}, fmt.Errorf("error creating copy: %w", err)
	}
	newDeck := Deck{ID: deckID, Name: copyName, DeckType: d.DeckType, CreatedAt: gd.CreateAt}

	deckCards, err := m.gormDeckRepo.listDeckCards(d.ID)
	if err != nil {
		return newDeck, fmt.Errorf("error loading source deck cards: %w", err)
	}
	var gdcs []gormDeckCard
	for _, dc := range deckCards {
		gdcs = append(gdcs, gormDeckCard{ID: xid.New().String(), DeckID: deckID, CardID: dc.CardID, CardName: dc.CardName, AssociationType: dc.AssociationType, Count: dc.Count})
	}
	err = m.gormDeckRepo.addCards(gdcs)
	if err != nil {
		return newDeck, fmt.Errorf("error saving cards %w", err)
	}
	return newDeck, nil
}

func (m *Manager) RemoveDeck(d Deck) error {
	err := m.gormDeckRepo.removeDeckCards(d.ID)
	if err != nil {
		return err
	}
	return m.gormDeckRepo.removeDeck(d.ID)
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
