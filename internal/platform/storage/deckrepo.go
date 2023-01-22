package storage

import "gorm.io/gorm"

type gormDeckRepo struct {
	db *gorm.DB
}

func newGormDeckRepo(db *gorm.DB) *gormDeckRepo {
	repo := &gormDeckRepo{db: db}
	repo.db.AutoMigrate(&gormDeck{}, &gormDeckCard{})

	return repo
}

func (r *gormDeckRepo) create(gd gormDeck) error {
	err := r.db.Create(&gd).Error
	return err
}

func (r *gormDeckRepo) addCard(gdc gormDeckCard) error {
	err := r.db.Create(&gdc).Error
	return err
}

func (r *gormDeckRepo) addCards(gdcs []gormDeckCard) error {
	err := r.db.CreateInBatches(gdcs, 20).Error
	return err
}

func (r *gormDeckRepo) ListDecks() ([]Deck, error) {
	var gds []gormDeck
	err := r.db.Find(&gds).Error
	if err != nil {
		return nil, err
	}

	var ds []Deck
	for _, gd := range gds {
		ds = append(ds, Deck{ID: gd.ID, Name: gd.Name, CreatedAt: gd.CreateAt, CoverImage: gd.CoverImage, DeckType: gd.DeckType})
	}
	return ds, nil
}

func (r *gormDeckRepo) UpdateCover(deckID, path string) error {
	return r.db.Model(&gormDeck{}).Where("id", deckID).Update("cover_image", path).Error
}

func (r *gormDeckRepo) findDeck(id string) (Deck, error) {
	var gd gormDeck
	err := r.db.Where("id = ?", id).First(&gd).Error
	if err != nil {
		return Deck{}, err
	}

	d := Deck{
		ID:         gd.ID,
		Name:       gd.Name,
		CreatedAt:  gd.CreateAt,
		CoverImage: gd.CoverImage,
		DeckType:   gd.DeckType,
	}

	return d, nil
}

func (r *gormDeckRepo) listDeckCards(deckID string) ([]gormDeckCard, error) {
	var gdcs []gormDeckCard
	err := r.db.Order("created_at asc").Where("deck_id = ?", deckID).Find(&gdcs).Error
	return gdcs, err
}
