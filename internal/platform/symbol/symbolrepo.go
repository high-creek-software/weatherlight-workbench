package symbol

import (
	"fyne.io/fyne/v2"
	"github.com/high-creek-software/goscryfall"
	"github.com/high-creek-software/goscryfall/symbols"
	"log"
	"sync"
)

type SymbolRepo interface {
	Get(code string) *symbols.Symbol
	Image(code string) fyne.Resource
}

type SymbolRepoImpl struct {
	client *goscryfall.Client

	cache  map[string]symbols.Symbol
	images map[string]fyne.Resource
	loader func(string) ([]byte, error)

	workerChan chan *symbols.Symbol
	locker     sync.RWMutex
}

func NewSymbolRepo(client *goscryfall.Client, loader func(string) ([]byte, error)) SymbolRepo {
	r := &SymbolRepoImpl{client: client, cache: make(map[string]symbols.Symbol), loader: loader, images: make(map[string]fyne.Resource), workerChan: make(chan *symbols.Symbol, 50)}
	ss, err := client.ListSymbols()
	if err != nil {
		log.Println(err)
		return r
	}

	workers := 8
	for idx := 0; idx < workers; idx++ {
		go r.loadAsset()
	}

	for _, s := range ss.Data {
		r.cache[s.Symbol] = s
		func(symb symbols.Symbol) {
			r.workerChan <- &symb
		}(s)
	}

	for idx := 0; idx < workers; idx++ {
		r.workerChan <- nil
	}

	return r
}

func (r *SymbolRepoImpl) loadAsset() {
	for {
		select {
		case s := <-r.workerChan:
			if s == nil {
				return
			}
			data, err := r.loader(s.SvgUri)
			if err != nil {
				log.Println("error loading image:", s.SvgUri, err)
			} else {
				r.locker.Lock()
				r.images[s.Symbol] = fyne.NewStaticResource(s.Symbol, data)
				r.locker.Unlock()
			}
		}
	}

}

func (r *SymbolRepoImpl) Get(code string) *symbols.Symbol {
	return nil
}

func (r *SymbolRepoImpl) Image(code string) fyne.Resource {
	return r.images[code]
}
