package sync

import (
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
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
	//bs, err := i.client.ListBulk()
	//if err != nil {
	//	return nil, nil, err
	//}

	sets, err := i.client.ListSets()
	if err != nil {
		return nil, nil, err
	}

	i.manager.StoreSets(sets)
	resChan := make(chan StatusUpdate)
	doneChan := make(chan bool)

	total := float64(len(sets))

	go func() {
		for idx, set := range sets {
			resChan <- StatusUpdate{Percent: (float64(idx) / total) * 100, SetName: set.Name}

			cards, err := i.client.ListCards(set.Code, "")
			if err == nil {
				i.manager.Store(cards.Data)

				time.Sleep(250 * time.Millisecond)
				for cards.HasMore {
					cards, err = i.client.ListCards(set.Code, cards.NextPage)
					if err == nil {
						i.manager.Store(cards.Data)
					}
					time.Sleep(250 * time.Millisecond)
				}
			} else {
				log.Println(err)
			}
		}

		doneChan <- true
	}()

	return resChan, doneChan, nil
}

type StatusUpdate struct {
	SetName string
	Percent float64
}
