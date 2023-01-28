package platform

import (
	"github.com/high-creek-software/ansel"
	"github.com/high-creek-software/goscryfall"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/notifier"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/symbol"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/sync"
)

type Registry struct {
	Manager       *storage.Manager
	SymbolRepo    symbol.SymbolRepo
	Client        *goscryfall.Client
	ImportManager *sync.ImportManager
	Notifier      notifier.Notifier

	SetIconLoader       *ansel.Ansel[string]
	CardThumbnailLoader *ansel.Ansel[string]
	CardFullLoader      *ansel.Ansel[string]
}

func NewRegistry(manager *storage.Manager, symbolRepo symbol.SymbolRepo, client *goscryfall.Client, importManager *sync.ImportManager, not notifier.Notifier) *Registry {
	reg := &Registry{
		Manager:       manager,
		SymbolRepo:    symbolRepo,
		Client:        client,
		ImportManager: importManager,
		Notifier:      not,
	}

	return reg
}
