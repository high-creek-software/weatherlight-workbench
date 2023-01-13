package storage

import "gorm.io/gorm"

type gormCard struct {
	Object               string
	Id                   string `gorm:"primaryKey"`
	OracleId             string
	MultiverseIds        string
	MtgoId               int
	TcgplayerId          int
	CardmarketId         int
	Name                 string `gorm:"index:idx_card_name"`
	Lang                 string
	ReleasedAt           string
	Uri                  string
	ScryfallUri          string
	Layout               string
	HighresImage         bool
	ImageStatus          string
	ImageUris            string
	ManaCost             string
	Cmc                  float64
	TypeLine             string
	FlavorText           string
	OracleText           string
	Power                string
	Toughness            string
	White                bool `gorm:"index:idx_card_white"`
	Blue                 bool `gorm:"index:idx_card_blue"`
	Black                bool `gorm:"index:idx_card_black"`
	Red                  bool `gorm:"index:idx_card_red"`
	Green                bool `gorm:"index:idx_card_green"`
	Keywords             string
	AllParts             string
	LegalStandard        string `gorm:"index:idx_card_legal_standard"`
	LegalFuture          string `gorm:"index:idx_card_legal_future"`
	LegalHistoric        string `gorm:"index:idx_card_legal_historic"`
	LegalGladiator       string `gorm:"index:idx_card_legal_gladiator"`
	LegalPioneer         string `gorm:"index:idx_card_legal_pioneer"`
	LegalExplorer        string `gorm:"index:idx_card_legal_explorer"`
	LegalModern          string `gorm:"index:idx_card_legal_modern"`
	LegalLegacy          string `gorm:"index:idx_card_legal_legacy"`
	LegalPauper          string `gorm:"index:idx_card_legal_pauper"`
	LegalVintage         string `gorm:"index:idx_card_legal_vintage"`
	LegalPenny           string `gorm:"index:idx_card_legal_penny"`
	LegalCommander       string `gorm:"index:idx_card_legal_commander"`
	LegalBrawl           string `gorm:"index:idx_card_legal_brawl"`
	LegalHistoricBrawl   string `gorm:"index:idx_card_legal_historicbrawl"`
	LegalAlchemy         string `gorm:"index:idx_card_legal_alchemy"`
	LegalPauperCommander string `gorm:"index:idx_card_legal_paupercommander"`
	LegalDuel            string `gorm:"index:idx_card_legal_duel"`
	LegalOldschool       string `gorm:"index:idx_card_legal_oldschool"`
	LegalPremodern       string `gorm:"index:idx_card_legal_premodern"`
	Games                string
	Reserved             bool
	Foil                 bool
	Nonfoil              bool
	Finishes             string
	Oversized            bool
	Promo                bool
	Reprint              bool
	Variation            bool
	SetId                string `gorm:"index:idx_card_set_id"`
	Set                  string `gorm:"index:idx_card_set"`
	SetName              string `gorm:"index:idx_card_set_name"`
	SetType              string
	SetUri               string
	SetSearchUri         string
	ScryfallSetUri       string
	RulingsUri           string
	PrintsSearchUri      string
	CollectorNumber      string
	Digital              bool
	Rarity               string `gorm:"index:idx_card_rarity"`
	CardBackId           string
	Artist               string
	ArtistIds            string
	IllustrationId       string
	BorderColor          string
	Frame                string
	FrameEffects         string
	SecurityStamp        string
	FullArt              bool
	Textless             bool
	Booster              bool
	StorySpotlight       bool
	EdhrecRank           int
	Preview              string
	Prices               string
	RelatedUris          string
	PurchaseUris         string
}

func (gormCard) TableName() string {
	return "cards"
}

type gormCardPrices struct {
	gorm.Model
	CardID    string `gorm:"index:idx_card_prices_card"`
	USD       float64
	USDFoil   float64
	USDEtched float64
	EUR       float64
	TIX       float64
}

func (gormCardPrices) TableName() string {
	return "card_prices"
}
