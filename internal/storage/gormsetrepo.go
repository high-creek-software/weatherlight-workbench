package storage

import (
	"gorm.io/gorm"
)
import scryfallset "gitlab.com/high-creek-software/goscryfall/sets"

type gormSetRepo struct {
	db *gorm.DB
}

func newGormSetRepo(db *gorm.DB) *gormSetRepo {
	repo := &gormSetRepo{db: db}
	repo.db.AutoMigrate(&gormSet{})
	return repo
}

func (r *gormSetRepo) StoreSets(sets []scryfallset.Set) error {
	var insert []gormSet
	for _, set := range sets {
		gs := gormSet{
			Object:        set.Object,
			Id:            set.Id,
			Code:          set.Code,
			Name:          set.Name,
			Uri:           set.Uri,
			ScryfallUri:   set.ScryfallUri,
			ReleasedAt:    set.ReleasedAt,
			SetType:       set.SetType,
			CardCount:     set.CardCount,
			ParentSetCode: set.ParentSetCode,
			Digital:       set.Digital,
			NonfoilOnly:   set.NonfoilOnly,
			FoilOnly:      set.FoilOnly,
			IconSvgUri:    set.IconSvgUri,
		}
		insert = append(insert, gs)
	}

	return r.db.Create(&insert).Error
}

func (r *gormSetRepo) ListSets() ([]scryfallset.Set, error) {
	var all []gormSet
	r.db.Find(&all)

	var res []scryfallset.Set
	for _, a := range all {
		s := scryfallset.Set{
			Object:        a.Object,
			Id:            a.Id,
			Code:          a.Code,
			Name:          a.Name,
			Uri:           a.Uri,
			ScryfallUri:   a.ScryfallUri,
			ReleasedAt:    a.ReleasedAt,
			SetType:       a.SetType,
			CardCount:     a.CardCount,
			ParentSetCode: a.ParentSetCode,
			Digital:       a.Digital,
			NonfoilOnly:   a.NonfoilOnly,
			FoilOnly:      a.FoilOnly,
			IconSvgUri:    a.IconSvgUri,
		}
		res = append(res, s)
	}

	return res, nil
}

func (r *gormSetRepo) SetCount() int64 {
	var count int64
	r.db.Model(&gormSet{}).Count(&count)
	return count
}
