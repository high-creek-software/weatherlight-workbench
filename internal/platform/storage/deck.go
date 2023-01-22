package storage

import (
	scryfallcards "gitlab.com/high-creek-software/goscryfall/cards"
	"time"
)

type Deck struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	CoverImage  string
	DeckType    string
	Commander   *DeckCard
	commanderID string
	Sideboard   []DeckCard
	Main        []DeckCard
}

type DeckCard struct {
	ID              string
	Count           int
	Card            scryfallcards.Card
	CreatedAt       time.Time
	AssociationType AssociationType
}

type gormDeck struct {
	ID         string `gorm:"primaryKey"`
	Name       string `gorm:"index:idx_deck_name"`
	CreateAt   time.Time
	CoverImage string
	DeckType   string `gorm:"index:idx_deck_type"`
}

func (gormDeck) TableName() string {
	return "decks"
}

type AssociationType int

const (
	AssociationUnknown AssociationType = iota
	AssociationSideboard
	AssociationMain
	AssociationCommander
	AssociationCompanion
)

type gormDeckCard struct {
	ID              string `gorm:"primaryKey"`
	DeckID          string `gorm:"index:idx_deck_card_deck_id"`
	CardID          string `gorm:"index:idx_deck_card_card_id"`
	CardName        string
	AssociationType AssociationType
	Count           int
	CreatedAt       time.Time
}

func (gormDeckCard) TableName() string {
	return "deck_cards"
}
