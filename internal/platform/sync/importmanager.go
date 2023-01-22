package sync

import (
	"gitlab.com/high-creek-software/goscryfall"
	cards2 "gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
	"log"
	"time"
)

type ImportManager struct {
	client  *goscryfall.Client
	manager *storage.Manager
}

func NewImportManager(client *goscryfall.Client, manager *storage.Manager) *ImportManager {
	return &ImportManager{client: client, manager: manager}
}

func (i *ImportManager) Import() (chan StatusUpdate, chan bool, error) {

	sets, err := i.client.ListSets()
	if err != nil {
		return nil, nil, err
	}

	i.manager.StoreSets(sets)
	resChan := make(chan StatusUpdate)
	doneChan := make(chan bool)

	total := float64(len(sets))

	go func() {
		start := time.Now()
		for idx, set := range sets {
			resChan <- StatusUpdate{Percent: (float64(idx) / total) * 100, SetName: set.Name}

			cards, err := i.client.ListCards(set.Code, cards2.UniquePrints, "")
			if err == nil {
				i.manager.Store(cards.Data)

				time.Sleep(150 * time.Millisecond)
				for cards.HasMore {
					cards, err = i.client.ListCards(set.Code, cards2.UniquePrints, cards.NextPage)
					if err == nil {
						i.manager.Store(cards.Data)
					}
					time.Sleep(100 * time.Millisecond)
				}
			} else {
				log.Println(err)
			}
		}
		log.Println(time.Now().Sub(start))
		doneChan <- true
		close(resChan)
		close(doneChan)
	}()

	return resChan, doneChan, nil
}

type StatusUpdate struct {
	SetName string
	Percent float64
}
