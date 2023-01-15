package platform

import (
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/sync"
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
