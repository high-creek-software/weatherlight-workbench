package storage

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
)
import scryfallcards "gitlab.com/high-creek-software/goscryfall/cards"

type gormCardRepo struct {
	db *gorm.DB
}

func newGormCardRepo(db *gorm.DB) *gormCardRepo {
	repo := &gormCardRepo{db: db}
	repo.db.AutoMigrate(&gormCard{}, &gormCardPrices{})
	return repo
}

func (r *gormCardRepo) Store(cs []scryfallcards.Card) error {
	var insert []gormCard
	var prices []gormCardPrices
	for _, c := range cs {

		white, blue, black, red, green := false, false, false, false, false
		for _, clr := range c.Colors {
			switch clr {
			case "W":
				white = true
			case "U":
				blue = true
			case "B":
				black = true
			case "R":
				red = true
			case "G":
				green = true
			}
		}
		i := gormCard{
			Object:               c.Object,
			Id:                   c.Id,
			OracleId:             c.OracleId,
			MultiverseIds:        marshal(c.MultiverseIds),
			MtgoId:               c.MtgoId,
			TcgplayerId:          c.TcgplayerId,
			CardmarketId:         c.CardmarketId,
			Name:                 c.Name,
			Lang:                 c.Lang,
			ReleasedAt:           c.ReleasedAt,
			Uri:                  c.Uri,
			ScryfallUri:          c.ScryfallUri,
			Layout:               c.Layout,
			HighresImage:         c.HighresImage,
			ImageStatus:          c.ImageStatus,
			ImageUris:            marshal(c.ImageUris),
			ManaCost:             c.ManaCost,
			Cmc:                  c.Cmc,
			TypeLine:             c.TypeLine,
			FlavorText:           c.FlavorText,
			OracleText:           c.OracleText,
			Power:                c.Power,
			Toughness:            c.Toughness,
			White:                white,
			Blue:                 blue,
			Black:                black,
			Red:                  red,
			Green:                green,
			Keywords:             marshal(c.Keywords),
			ProducedMana:         marshal(c.ProducedMana),
			CardFaces:            marshal(c.CardFaces),
			AllParts:             marshal(c.AllParts),
			LegalStandard:        c.Legalities["standard"].String(),
			LegalFuture:          c.Legalities["future"].String(),
			LegalHistoric:        c.Legalities["historic"].String(),
			LegalGladiator:       c.Legalities["gladiator"].String(),
			LegalPioneer:         c.Legalities["pioneer"].String(),
			LegalExplorer:        c.Legalities["explorer"].String(),
			LegalModern:          c.Legalities["modern"].String(),
			LegalLegacy:          c.Legalities["legacy"].String(),
			LegalPauper:          c.Legalities["pauper"].String(),
			LegalVintage:         c.Legalities["vintage"].String(),
			LegalPenny:           c.Legalities["penny"].String(),
			LegalCommander:       c.Legalities["commander"].String(),
			LegalBrawl:           c.Legalities["brawl"].String(),
			LegalHistoricBrawl:   c.Legalities["historicbrawl"].String(),
			LegalAlchemy:         c.Legalities["alchemy"].String(),
			LegalPauperCommander: c.Legalities["paupercommander"].String(),
			LegalDuel:            c.Legalities["duel"].String(),
			LegalOldschool:       c.Legalities["oldschool"].String(),
			LegalPremodern:       c.Legalities["premodern"].String(),
			Games:                marshal(c.Games),
			Reserved:             c.Reserved,
			Foil:                 c.Foil,
			Nonfoil:              c.Nonfoil,
			Finishes:             marshal(c.Finishes),
			Oversized:            c.Oversized,
			Promo:                c.Promo,
			Reprint:              c.Reprint,
			Variation:            c.Variation,
			SetId:                c.SetId,
			Set:                  c.Set,
			SetName:              c.SetName,
			SetType:              c.SetType,
			SetUri:               c.SetUri,
			SetSearchUri:         c.SetSearchUri,
			ScryfallSetUri:       c.ScryfallSetUri,
			RulingsUri:           c.RulingsUri,
			PrintsSearchUri:      c.PrintsSearchUri,
			CollectorNumber:      c.CollectorNumber,
			Digital:              c.Digital,
			Rarity:               c.Rarity,
			CardBackId:           c.CardBackId,
			Artist:               c.Artist,
			ArtistIds:            marshal(c.ArtistIds),
			IllustrationId:       c.IllustrationId,
			BorderColor:          c.BorderColor,
			Frame:                c.Frame,
			FrameEffects:         marshal(c.FrameEffects),
			SecurityStamp:        c.SecurityStamp,
			FullArt:              c.FullArt,
			Textless:             c.Textless,
			Booster:              c.Booster,
			StorySpotlight:       c.StorySpotlight,
			EdhrecRank:           c.EdhrecRank,
			Preview:              marshal(c.Preview),
			Prices:               marshal(c.Prices),
			RelatedUris:          marshal(c.RelatedUris),
			PurchaseUris:         marshal(c.PurchaseUris),
		}
		insert = append(insert, i)

		price := gormCardPrices{CardID: c.Id}
		hasPrice := false
		if c.Prices.Usd != "" {
			hasPrice = true
			price.USD, _ = strconv.ParseFloat(c.Prices.Usd, 64)
		}
		if c.Prices.UsdFoil != "" {
			hasPrice = true
			price.USDFoil, _ = strconv.ParseFloat(c.Prices.UsdFoil, 64)
		}
		if ue, ok := c.Prices.UsdEtched.(string); ok && ue != "" {
			hasPrice = true
			price.USDEtched, _ = strconv.ParseFloat(ue, 64)
		}
		if c.Prices.Eur != "" {
			hasPrice = true
			price.EUR, _ = strconv.ParseFloat(c.Prices.Eur, 64)
		}
		if c.Prices.Tix != "" {
			hasPrice = true
			price.TIX, _ = strconv.ParseFloat(c.Prices.Tix, 64)
		}
		if hasPrice {
			prices = append(prices, price)
		}

	}
	pricesBatchErr := r.db.CreateInBatches(&prices, 50).Error
	if pricesBatchErr != nil {
		log.Println("Error inserting prices", pricesBatchErr)
	}
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(&insert, 50).Error
}

func (r *gormCardRepo) ListBySet(set string) ([]scryfallcards.Card, error) {
	var gcs []gormCard
	r.db.Where("`set` = ?", set).Find(&gcs)

	var res []scryfallcards.Card
	for _, gc := range gcs {
		c := internalToExternal(gc)
		res = append(res, c)
	}

	return res, nil
}

func (r *gormCardRepo) FindByName(name string) ([]scryfallcards.Card, error) {
	var cs []gormCard
	err := r.db.Where("name = ?", name).Find(&cs).Error
	if err != nil || len(cs) == 0 {

		err = r.db.Where("name LIKE ?", "%"+name+"%").Find(&cs).Error
		if err != nil {
			return nil, err
		}
	}

	var resp []scryfallcards.Card
	for _, c := range cs {
		resp = append(resp, internalToExternal(c))
	}

	return resp, err
}

func (r *gormCardRepo) findById(id string) (scryfallcards.Card, error) {
	var c gormCard
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		return scryfallcards.Card{}, err
	}

	return internalToExternal(c), err
}

func (r *gormCardRepo) ListByIds(ids []string) ([]scryfallcards.Card, error) {
	var gcs []gormCard
	err := r.db.Where("id IN ?", ids).Order("name asc").Find(&gcs).Error
	if err != nil {
		return nil, err
	}

	var res []scryfallcards.Card
	for _, gc := range gcs {
		c := internalToExternal(gc)
		res = append(res, c)
	}

	return res, nil
}

func (r *gormCardRepo) Search(sr SearchRequest) ([]scryfallcards.Card, error) {
	queryDB := r.db.Session(&gorm.Session{})
	if sr.Name != "" {
		queryDB = queryDB.Where("name LIKE ?", "%"+sr.Name+"%")
	}
	if sr.TypeLine != "" {
		queryDB = queryDB.Where("type_line LIKE ?", "%"+sr.TypeLine+"%")
	}
	if sr.OracleText != "" {
		queryDB = queryDB.Where("oracle_text LIKE ?", "%"+sr.OracleText+"%")
	}
	if sr.White {
		queryDB = queryDB.Where("white = true")
	}
	if sr.Blue {
		queryDB = queryDB.Where("blue = true")
	}
	if sr.Black {
		queryDB = queryDB.Where("black = true")
	}
	if sr.Red {
		queryDB = queryDB.Where("red = true")
	}
	if sr.Green {
		queryDB = queryDB.Where("green = true")
	}

	if sr.StandardLegal {
		queryDB = queryDB.Where("legal_standard = ?", scryfallcards.Legal.String())
	}
	if sr.FutureLegal {
		queryDB = queryDB.Where("legal_future = ?", scryfallcards.Legal.String())
	}
	if sr.HistoricLegal {
		queryDB = queryDB.Where("legal_historic = ?", scryfallcards.Legal.String())
	}
	if sr.GladiatorLegal {
		queryDB = queryDB.Where("legal_gladiator = ?", scryfallcards.Legal.String())
	}
	if sr.PioneerLegal {
		queryDB = queryDB.Where("legal_pioneer = ?", scryfallcards.Legal.String())
	}
	if sr.ExplorerLegal {
		queryDB = queryDB.Where("legal_explorer = ?", scryfallcards.Legal.String())
	}
	if sr.ModernLegal {
		queryDB = queryDB.Where("legal_modern = ?", scryfallcards.Legal.String())
	}
	if sr.LegacyLegal {
		queryDB = queryDB.Where("legal_legacy = ?", scryfallcards.Legal.String())
	}
	if sr.PauperLegal {
		queryDB = queryDB.Where("legal_pauper = ?", scryfallcards.Legal.String())
	}
	if sr.VintageLegal {
		queryDB = queryDB.Where("legal_vintage = ?", scryfallcards.Legal.String())
	}
	if sr.PennyLegal {
		queryDB = queryDB.Where("legal_penny = ?", scryfallcards.Legal.String())
	}
	if sr.CommanderLegal {
		queryDB = queryDB.Where("legal_commander = ?", scryfallcards.Legal.String())
	}
	if sr.BrawlLegal {
		queryDB = queryDB.Where("legal_brawl = ?", scryfallcards.Legal.String())
	}
	if sr.HistoricBrawlLegal {
		queryDB = queryDB.Where("legal_historic_brawl = ?", scryfallcards.Legal.String())
	}
	if sr.AlchemyLegal {
		queryDB = queryDB.Where("legal_alchemy = ?", scryfallcards.Legal.String())
	}
	if sr.PauperCommanderLegal {
		queryDB = queryDB.Where("legal_pauper_commander = ?", scryfallcards.Legal.String())
	}
	if sr.DuelLegal {
		queryDB = queryDB.Where("legal_duel = ?", scryfallcards.Legal.String())
	}
	if sr.OldschoolLegal {
		queryDB = queryDB.Where("legal_oldschool = ?", scryfallcards.Legal.String())
	}
	if sr.PremodernLegal {
		queryDB = queryDB.Where("legal_premodern = ?", scryfallcards.Legal.String())
	}

	queryDB = queryDB.Order("name asc")

	var gcs []gormCard
	err := queryDB.Find(&gcs).Error
	if err != nil {
		return nil, err
	}

	var res []scryfallcards.Card
	for _, gc := range gcs {
		c := internalToExternal(gc)
		res = append(res, c)
	}

	return res, nil
}

func internalToExternal(gc gormCard) scryfallcards.Card {
	c := scryfallcards.Card{
		Object:          gc.Object,
		Id:              gc.Id,
		OracleId:        gc.OracleId,
		MultiverseIds:   unmarshal[[]int](gc.MultiverseIds),
		MtgoId:          gc.MtgoId,
		TcgplayerId:     gc.TcgplayerId,
		CardmarketId:    gc.CardmarketId,
		Name:            gc.Name,
		Lang:            gc.Lang,
		ReleasedAt:      gc.ReleasedAt,
		Uri:             gc.Uri,
		ScryfallUri:     gc.ScryfallUri,
		Layout:          gc.Layout,
		HighresImage:    gc.HighresImage,
		ImageStatus:     gc.ImageStatus,
		ImageUris:       unmarshal[scryfallcards.ImageUris](gc.ImageUris),
		ManaCost:        gc.ManaCost,
		Cmc:             gc.Cmc,
		TypeLine:        gc.TypeLine,
		FlavorText:      gc.FlavorText,
		OracleText:      gc.OracleText,
		Power:           gc.Power,
		Toughness:       gc.Toughness,
		Keywords:        unmarshal[[]interface{}](gc.Keywords),
		ProducedMana:    unmarshal[[]string](gc.ProducedMana),
		CardFaces:       unmarshal[[]scryfallcards.CardFace](gc.CardFaces),
		AllParts:        unmarshal[[]scryfallcards.AllParts](gc.AllParts),
		Games:           unmarshal[[]string](gc.Games),
		Reserved:        gc.Reserved,
		Foil:            gc.Foil,
		Nonfoil:         gc.Nonfoil,
		Finishes:        unmarshal[[]string](gc.Finishes),
		Oversized:       gc.Oversized,
		Promo:           gc.Promo,
		Reprint:         gc.Reprint,
		Variation:       gc.Variation,
		SetId:           gc.SetId,
		Set:             gc.Set,
		SetName:         gc.SetName,
		SetType:         gc.SetType,
		SetUri:          gc.SetUri,
		SetSearchUri:    gc.SetSearchUri,
		ScryfallSetUri:  gc.ScryfallSetUri,
		RulingsUri:      gc.RulingsUri,
		PrintsSearchUri: gc.PrintsSearchUri,
		CollectorNumber: gc.CollectorNumber,
		Digital:         gc.Digital,
		Rarity:          gc.Rarity,
		CardBackId:      gc.CardBackId,
		Artist:          gc.Artist,
		ArtistIds:       unmarshal[[]string](gc.ArtistIds),
		IllustrationId:  gc.IllustrationId,
		BorderColor:     gc.BorderColor,
		Frame:           gc.Frame,
		FrameEffects:    unmarshal[[]string](gc.FrameEffects),
		SecurityStamp:   gc.SecurityStamp,
		FullArt:         gc.FullArt,
		Textless:        gc.Textless,
		Booster:         gc.Booster,
		StorySpotlight:  gc.StorySpotlight,
		EdhrecRank:      gc.EdhrecRank,
		Preview:         unmarshal[scryfallcards.Preview](gc.Preview),
		Prices:          unmarshal[scryfallcards.Prices](gc.Prices),
		RelatedUris:     unmarshal[scryfallcards.RelatedUris](gc.RelatedUris),
		PurchaseUris:    unmarshal[scryfallcards.PurchaseUris](gc.RelatedUris),
	}

	if gc.White {
		c.Colors = append(c.Colors, "W")
	}
	if gc.Blue {
		c.Colors = append(c.Colors, "U")
	}
	if gc.Black {
		c.Colors = append(c.Colors, "B")
	}
	if gc.Red {
		c.Colors = append(c.Colors, "R")
	}
	if gc.Green {
		c.Colors = append(c.Colors, "G")
	}

	c.Legalities = make(map[string]scryfallcards.Legality)
	c.Legalities["standard"] = scryfallcards.LegalityFromString(gc.LegalStandard)
	c.Legalities["future"] = scryfallcards.LegalityFromString(gc.LegalFuture)
	c.Legalities["historic"] = scryfallcards.LegalityFromString(gc.LegalHistoric)
	c.Legalities["gladiator"] = scryfallcards.LegalityFromString(gc.LegalGladiator)
	c.Legalities["pioneer"] = scryfallcards.LegalityFromString(gc.LegalPioneer)
	c.Legalities["explorer"] = scryfallcards.LegalityFromString(gc.LegalExplorer)
	c.Legalities["modern"] = scryfallcards.LegalityFromString(gc.LegalModern)
	c.Legalities["legacy"] = scryfallcards.LegalityFromString(gc.LegalLegacy)
	c.Legalities["pauper"] = scryfallcards.LegalityFromString(gc.LegalPauper)
	c.Legalities["vintage"] = scryfallcards.LegalityFromString(gc.LegalVintage)
	c.Legalities["penny"] = scryfallcards.LegalityFromString(gc.LegalPenny)
	c.Legalities["commander"] = scryfallcards.LegalityFromString(gc.LegalCommander)
	c.Legalities["brawl"] = scryfallcards.LegalityFromString(gc.LegalBrawl)
	c.Legalities["historicbrawl"] = scryfallcards.LegalityFromString(gc.LegalHistoricBrawl)
	c.Legalities["alchemy"] = scryfallcards.LegalityFromString(gc.LegalAlchemy)
	c.Legalities["paupercommander"] = scryfallcards.LegalityFromString(gc.LegalPauperCommander)
	c.Legalities["duel"] = scryfallcards.LegalityFromString(gc.LegalDuel)
	c.Legalities["oldschool"] = scryfallcards.LegalityFromString(gc.LegalOldschool)
	c.Legalities["premodern"] = scryfallcards.LegalityFromString(gc.LegalPremodern)

	return c
}

func marshal(input any) string {
	res, err := json.Marshal(input)
	if err != nil {
		res = []byte("")
	}
	return string(res)
}

func unmarshal[T any](data string) T {
	var res T
	json.Unmarshal([]byte(data), &res)
	return res
}
