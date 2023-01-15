package storage

type SearchRequest struct {
	Name       string
	TypeLine   string
	OracleText string

	White bool
	Blue  bool
	Black bool
	Red   bool
	Green bool

	StandardLegal        bool
	FutureLegal          bool
	HistoricLegal        bool
	GladiatorLegal       bool
	PioneerLegal         bool
	ExplorerLegal        bool
	ModernLegal          bool
	LegacyLegal          bool
	PauperLegal          bool
	VintageLegal         bool
	PennyLegal           bool
	CommanderLegal       bool
	BrawlLegal           bool
	HistoricBrawlLegal   bool
	AlchemyLegal         bool
	PauperCommanderLegal bool
	DuelLegal            bool
	OldschoolLegal       bool
	PremodernLegal       bool
}
