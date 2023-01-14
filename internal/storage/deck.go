package storage

import (
	cards2 "gitlab.com/high-creek-software/goscryfall/cards"
	"time"
)

type Deck struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	CoverImage  string
	Commander   *cards2.Card
	commanderID string
	Sideboard   []cards2.Card
	Main        []cards2.Card
}

type gormDeck struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	CreateAt    time.Time
	CoverImage  string
	CommanderID string
}

func (gormDeck) TableName() string {
	return "decks"
}

type associationType int

const (
	associationUnknown associationType = iota
	associationSideboard
	associationMain
	associationCommander
)

type gormDeckCard struct {
	ID              string `gorm:"primaryKey"`
	DeckID          string `gorm:"index:idx_deck_card_deck_id"`
	CardID          string `gorm:"index:idx_deck_card_card_id"`
	CardName        string
	AssociationType associationType
	Count           int
}

func (gormDeckCard) TableName() string {
	return "deck_cards"
}
