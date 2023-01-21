package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/bento"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/bookmarked"
	"gitlab.com/kendellfab/mtgstudio/internal/browse"
	"gitlab.com/kendellfab/mtgstudio/internal/deck"
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/icons"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/sync"
	"gitlab.com/kendellfab/mtgstudio/internal/search"
	"golang.org/x/image/colornames"
	"strings"
	"time"
)

const (
	lastSyncKey = "last_sync_key"
	syncFormat  = "2006-01-02 15:04:05"
)

type MtgStudio struct {
	app              fyne.App
	window           fyne.Window
	browseLayout     *browse.BrowseLayout
	searchLayout     *search.SearchLayout
	bookmarkedLayout *bookmarked.BookmarkedLayout
	deckLayout       *deck.DeckLayout

	registry *platform.Registry

	syncBtn      *widget.Button
	syncLastLbl  *widget.Label
	syncProgress *widget.ProgressBar
	syncSetLbl   *widget.Label
	settingsBtn  *widget.Button

	bentoBox *bento.Box
}

func NewMtgStudio() *MtgStudio {
	//os.Setenv("FYNE_THEME", "dark")
	mtgs := &MtgStudio{app: app.NewWithID("gitlab.com/kendellfab/mtgstudio")}
	mtgs.app.SetIcon(icons.AppIconResource)
	mtgs.window = mtgs.app.NewWindow("MTG Studio")
	mtgs.window.SetMaster()
	mtgs.window.Resize(fyne.NewSize(1920, 1080))
	client := goscryfall.NewClient()
	manager := storage.NewManager(client)
	importManager := sync.NewImportManager(client, manager)
	symbolRepo := symbol.NewSymbolRepo(client, manager.LoadSymbolImage)

	mtgs.registry = platform.NewRegistry(manager, symbolRepo, client, importManager, mtgs)
	mtgs.registry.SymbolRepo = symbolRepo

	mtgs.registry.SetIconLoader = ansel.NewAnsel[string](200, ansel.SetLoadedCallback[string](mtgs.updateSetIcon), ansel.SetLoader[string](mtgs.registry.Manager.LoadSetIcon))
	mtgs.registry.CardThumbnailLoader = ansel.NewAnsel[string](800, ansel.SetLoader[string](mtgs.registry.Manager.LoadCardImage), ansel.SetWorkerCount[string](20), ansel.SetLoadingImage[string](icons.CardLoadingResource), ansel.SetFailedImage[string](icons.CardFailedResource))
	mtgs.registry.CardFullLoader = ansel.NewAnsel[string](200, ansel.SetLoader[string](mtgs.registry.Manager.LoadCardImage), ansel.SetLoadingImage[string](icons.FullCardLoadingResource), ansel.SetFailedImage[string](icons.FullCardFailedResource))

	mtgs.app.Lifecycle().SetOnStarted(mtgs.appStartedCallback)

	mtgs.setupBody()

	return mtgs
}

func (m *MtgStudio) setupBody() {
	m.browseLayout = browse.NewBrowseLayout(m.registry, m.updateSetIcon, m.resizeCardArt)
	m.searchLayout = search.NewSearchLayout(m.registry)
	m.bookmarkedLayout = bookmarked.NewBookmarkedLayout(m.registry)
	m.deckLayout = deck.NewDeckLayout(m.registry, m.showImport)
	appTabs := container.NewAppTabs(container.NewTabItem("Browse", m.browseLayout.Split),
		container.NewTabItem("Search", m.searchLayout.Split),
		container.NewTabItem("Bookmarked", m.bookmarkedLayout.Split),
		container.NewTabItem("Decks", m.deckLayout.Container),
	)

	m.syncBtn = widget.NewButton("Sync", m.syncBtnTouched)
	m.syncLastLbl = widget.NewLabel("Last Synced")
	m.syncProgress = widget.NewProgressBar()
	m.syncProgress.SetValue(0)
	m.syncProgress.Max = 100
	m.syncProgress.Hide()
	m.syncSetLbl = widget.NewLabel("Set Name")
	m.syncSetLbl.Hide()
	m.settingsBtn = widget.NewButtonWithIcon("", theme.SettingsIcon(), m.settingsBtnTouched)

	// container.NewBorder(nil, nil, m.syncSetLbl, nil, m.syncProgress)
	syncBorder := container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), m.settingsBtn, layout.NewSpacer()), nil, container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), m.syncBtn, layout.NewSpacer()), nil, container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), m.syncLastLbl, layout.NewSpacer()), nil, container.NewVBox(m.syncSetLbl, m.syncProgress))))
	mainBody := container.NewBorder(nil, syncBorder, nil, nil, appTabs)

	m.bentoBox = bento.NewBox()
	m.bentoBox.UpdateBottomOffset(45)

	m.window.SetContent(container.NewMax(mainBody, m.bentoBox))

	appTabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Bookmarked" {
			m.bookmarkedLayout.LoadBookmarked()
		}
	}
}

func (m *MtgStudio) showImport() {
	window := m.app.NewWindow("Import Deck")
	window.Resize(fyne.NewSize(800, 400))

	entry := widget.NewEntry()
	entry.PlaceHolder = "Deck Name"

	data := widget.NewEntry()
	data.MultiLine = true

	//deckTypes := maps.Keys(cards.Legality())
	//deckType := widget.NewSelect()

	save := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		err := m.registry.Manager.ImportDeck(entry.Text, data.Text, "Unknown")
		window.Hide()
		window = nil
		if err != nil {
			m.ShowError(err)
			return
		}
		m.deckLayout.LoadDecks()
	})

	border := container.NewBorder(entry, save, nil, nil, data)
	window.SetContent(border)
	window.Show()
}

func (m *MtgStudio) updateSetIcon(bs []byte) []byte {
	if m.app.Settings().ThemeVariant() == theme.VariantDark {
		strData := string(bs)
		if strings.Contains(strData, `fill="#000"`) {
			strData = strings.Replace(strData, `fill="#000"`, `fill="#999"`, -1)
		} else {
			strData = strings.Replace(strData, "<path d=", `<path style="fill:#999999" d=`, -1)
		}
		return []byte(strData)
	}
	return bs
}

func (m *MtgStudio) resizeCardArt(bs []byte) []byte {
	return bs
	//buff := bytes.NewBuffer(bs)
	//img, err := jpeg.Decode(buff)
	//if err != nil {
	//	log.Println("error parsing image:", err)
	//	return bs
	//}
	//
	//r := resize.Resize(150, 0, img, resize.Lanczos3)
	//
	//var out bytes.Buffer
	//jpeg.Encode(&out, r, nil)
	//return out.Bytes()
}

func (m *MtgStudio) Start() {
	m.window.ShowAndRun()
}

func (m *MtgStudio) settingsBtnTouched() {
	dialog.ShowInformation("Settings", "Does nothing yet", m.window)
}

func (m *MtgStudio) syncBtnTouched() {
	m.runSync()
}

func (m *MtgStudio) appStartedCallback() {
	// TODO: Figure out how to determine if an import is needed.
	m.showLastSyncedAt()
	setCount := m.registry.Manager.SetCount()
	if setCount == 0 {
		m.runSync()
	} else {
		m.browseLayout.LoadSets()
	}
}

func (m *MtgStudio) runSync() {
	m.syncSetLbl.Show()
	m.syncProgress.Show()
	startCount := m.registry.Manager.SetCount()
	resChan, doneChan, err := m.registry.ImportManager.Import()
	if err != nil {
		m.ShowError(err)
	} else {
		go func() {
		F:
			for {
				select {
				case status := <-resChan:
					m.syncSetLbl.SetText(status.SetName)
					m.syncProgress.SetValue(status.Percent)
					if startCount == 0 {
						m.browseLayout.LoadSets()
					}
				case <-doneChan:
					m.syncSetLbl.Hide()
					m.syncProgress.Hide()
					m.browseLayout.LoadSets()
					break F
				}
			}
			synced := time.Now()
			m.app.Preferences().SetString(lastSyncKey, synced.Format(time.RFC3339))
			m.showLastSyncedAt()
		}()
	}
}

func (m *MtgStudio) showLastSyncedAt() {
	syncedStr := m.app.Preferences().StringWithFallback(lastSyncKey, "--")
	if syncedStr == "--" {
		m.syncLastLbl.SetText(syncedStr)
		return
	}

	if synced, err := time.Parse(time.RFC3339, syncedStr); err == nil {
		m.syncLastLbl.SetText(synced.Format(syncFormat))
	}

}

func (m *MtgStudio) ShowDialog(title, message string) {
	item := bento.NewItemWithMessage(message, bento.LengthIndefinite)
	item.SetBackgroundColor(colornames.Darkslateblue)
	m.bentoBox.AddItem(item)
}

func (m *MtgStudio) ShowError(err error) {
	item := bento.NewItemWithMessage(err.Error(), bento.LengthIndefinite)
	item.SetBackgroundColor(colornames.Rosybrown)
	m.bentoBox.AddItem(item)
}
