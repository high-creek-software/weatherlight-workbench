package storage

import (
	"errors"
	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	CardID string `gorm:"index:idx_bookmark_card_id"`
}

type bookmarkRepo struct {
	db *gorm.DB
}

func NewBookmarkRepo(db *gorm.DB) *bookmarkRepo {
	db.AutoMigrate(&Bookmark{})

	repo := &bookmarkRepo{db: db}

	return repo
}

func (r *bookmarkRepo) AddBookmark(cardID string) error {
	b := Bookmark{CardID: cardID}

	err := r.db.Create(&b).Error
	return err
}

func (r *bookmarkRepo) RemoveBookmark(cardID string) error {
	err := r.db.Where("card_id = ?", cardID).Delete(&Bookmark{}).Error
	return err
}

func (r *bookmarkRepo) FindBookmark(cardID string) (*Bookmark, error) {
	var b Bookmark
	err := r.db.Where("card_id = ?", cardID).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &b, err
}

func (r *bookmarkRepo) List() ([]Bookmark, error) {
	var bs []Bookmark
	err := r.db.Find(&bs).Error
	return bs, err
}
