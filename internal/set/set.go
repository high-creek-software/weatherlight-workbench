package set

type Set struct {
	Object        string `json:"object"`
	Id            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Uri           string `json:"uri"`
	ScryfallUri   string `json:"scryfall_uri"`
	SearchUri     string `json:"search_uri"`
	ReleasedAt    string `json:"released_at"`
	SetType       string `json:"set_type"`
	CardCount     int    `json:"card_count"`
	ParentSetCode string `json:"parent_set_code"`
	Digital       bool   `json:"digital"`
	NonfoilOnly   bool   `json:"nonfoil_only"`
	FoilOnly      bool   `json:"foil_only"`
	IconSvgUri    string `json:"icon_svg_uri"`
}
