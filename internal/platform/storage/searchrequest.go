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

func (s SearchRequest) IsEmpty() bool {
	if s.Name != "" {
		return false
	}

	if s.TypeLine != "" {
		return false
	}

	if s.OracleText != "" {
		return false
	}

	if s.White {
		return false
	}

	if s.Blue {
		return false
	}

	if s.Red {
		return false
	}

	if s.Green {
		return false
	}

	if s.Black {
		return false
	}

	if s.StandardLegal {
		return false
	}

	if s.FutureLegal {
		return false
	}

	if s.HistoricLegal {
		return false
	}

	if s.GladiatorLegal {
		return false
	}

	if s.PioneerLegal {
		return false
	}

	if s.ExplorerLegal {
		return false
	}

	if s.ModernLegal {
		return false
	}

	if s.LegacyLegal {
		return false
	}

	if s.PauperLegal {
		return false
	}

	if s.VintageLegal {
		return false
	}

	if s.PennyLegal {
		return false
	}

	if s.CommanderLegal {
		return false
	}

	if s.BrawlLegal {
		return false
	}

	if s.HistoricBrawlLegal {
		return false
	}

	if s.AlchemyLegal {
		return false
	}

	if s.PauperCommanderLegal {
		return false
	}

	if s.DuelLegal {
		return false
	}

	if s.OldschoolLegal {
		return false
	}

	if s.PremodernLegal {
		return false
	}

	return true
}
