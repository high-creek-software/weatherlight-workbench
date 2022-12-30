package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/api"
	"gitlab.com/kendellfab/mtgstudio/internal/resources"
	"gitlab.com/kendellfab/mtgstudio/internal/set"
	"log"
	"os"
	"sync"
)

const (
	scryfallBaseURL = "https://api.scryfall.com"
)

type MtgStudio struct {
	app    fyne.App
	window fyne.Window

	setList     *widget.List
	sets        []set.Set
	setIcons    map[string]*fyne.StaticResource
	setLocker   sync.RWMutex
	pendingLoad map[string]*widget.Icon

	manager *resources.Manager

	endpoint *api.Endpoint
	setRepo  set.SetRepo
}

func NewMtgStudio() *MtgStudio {
	os.Setenv("FYNE_THEME", "light")
	mtgs := &MtgStudio{app: app.New(), setIcons: make(map[string]*fyne.StaticResource), pendingLoad: make(map[string]*widget.Icon)}
	mtgs.window = mtgs.app.NewWindow("MTG Studio")
	mtgs.window.Resize(fyne.NewSize(1200, 700))
	mtgs.endpoint = api.NewEndpoint(scryfallBaseURL)
	mtgs.setRepo = set.NewRestSetRepo(mtgs.endpoint)
	mtgs.manager = resources.NewManager()

	mtgs.setupBody()

	return mtgs
}

func (m *MtgStudio) setupBody() {
	m.setList = widget.NewList(m.setCount, m.createSetTemplate, m.updateSetTemplate)
	m.window.SetContent(m.setList)
}

func (m *MtgStudio) setCount() int {
	if m.sets == nil {
		return 0
	}
	return len(m.sets)
}

func (m *MtgStudio) createSetTemplate() fyne.CanvasObject {
	icon := widget.NewIcon(nil)
	lbl := widget.NewLabel("template")
	//return container.NewGridWithColumns(2, widget.NewIcon(nil), widget.NewLabel("template"))
	return container.New(layout.NewFormLayout(), icon, lbl)
}

func (m *MtgStudio) updateSetTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	set := m.sets[id]
	icon := co.(*fyne.Container).Objects[0].(*widget.Icon)
	lbl := co.(*fyne.Container).Objects[1].(*widget.Label)

	go m.loadIcon(set, icon)

	lbl.SetText(set.Name)
}

func (m *MtgStudio) loadIcon(set set.Set, icon *widget.Icon) {

	// Check the cache for if this image has already been loaded.
	m.setLocker.RLock()
	if resource, ok := m.setIcons[set.Name]; ok {
		icon.SetResource(resource)
		m.setLocker.RUnlock()
		return
	}
	m.setLocker.RUnlock()

	// The callback to load data could happen multiple times for a set item, and a load could be pending.
	// So we check here, if the pendingLoad has the set name, we update the icon to load to and return.
	m.setLocker.Lock()
	if _, ok := m.pendingLoad[set.Name]; ok {
		m.pendingLoad[set.Name] = icon
		m.setLocker.Unlock()
		return
	}
	// If not, we create the pending load record and move on to loading.
	m.pendingLoad[set.Name] = icon
	m.setLocker.Unlock()

	log.Println("loading image for:", set.Name, set.Id, set.IconSvgUri)
	// We make the actual request right here.
	if data, err := api.RequestSetResource(set.IconSvgUri); err != nil {
		log.Println("error loading images", err)
	} else {
		// Here we load the byte array to a static resource
		// Add it to the cache
		// Set the image to the icon that was last added to the pendingLoad
		// Then remove it from the pending load
		m.setLocker.Lock()
		res := fyne.NewStaticResource(set.IconSvgUri, data)
		m.setIcons[set.Name] = res
		if icn, ok := m.pendingLoad[set.Name]; ok {
			icn.SetResource(res)
		}
		delete(m.pendingLoad, set.Name)
		m.setLocker.Unlock()
	}
}

func (m *MtgStudio) Start() {
	m.loadSets()
	m.window.ShowAndRun()
}

func (m *MtgStudio) loadSets() {
	sets, err := m.setRepo.ListSets()
	if err != nil {
		dialog.NewInformation("Error loading sets", err.Error(), m.window).Show()
		return
	}

	m.sets = sets
}
