package storage

type gormSet struct {
	Object        string
	Id            string `gorm:"primaryKey"`
	Code          string
	Name          string
	Uri           string
	ScryfallUri   string
	ReleasedAt    string
	SetType       string
	CardCount     int
	ParentSetCode string
	Digital       bool
	NonfoilOnly   bool
	FoilOnly      bool
	IconSvgUri    string
}

func (gormSet) TableName() string {
	return "sets"
}
