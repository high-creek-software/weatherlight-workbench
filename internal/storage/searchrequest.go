package storage

type SearchRequest struct {
	Name     string
	TypeLine string
	White    bool
	Blue     bool
	Black    bool
	Red      bool
	Green    bool

	BrawlLegal bool
}
