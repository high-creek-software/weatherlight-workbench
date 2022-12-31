package symbol

import (
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/high-creek-software/goscryfall/symbols"
	"log"
)

type SymbolRepo interface {
	Get(code string) *symbols.Symbol
}

type SymbolRepoImpl struct {
	client *goscryfall.Client

	cache map[string]symbols.Symbol
}

func NewSymbolRepo(client *goscryfall.Client) SymbolRepo {
	r := &SymbolRepoImpl{client: client, cache: make(map[string]symbols.Symbol)}
	symbols, err := client.ListSymbols()
	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range symbols.Data {
		r.cache[s.Symbol] = s
	}

	return r
}

func (r *SymbolRepoImpl) Get(code string) *symbols.Symbol {
	return nil
}
