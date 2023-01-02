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

func (i *ImportManager) Import() (chan float64, chan bool, error) {
	//bs, err := i.client.ListBulk()
	//if err != nil {
	//	return nil, nil, err
	//}

	sets, err := i.client.ListSets()
	if err != nil {
		return nil, nil, err
	}

	i.manager.StoreSets(sets)
	resChan := make(chan float64)
	doneChan := make(chan bool)

	total := float64(len(sets))

	go func() {
		for idx, set := range sets {
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

			resChan <- (float64(idx) / total) * 100
		}

		doneChan <- true
	}()

	/*var cardBulk bulk.Bulk
	for _, b := range bs.Data {
		if b.Type == "default_cards" {
			cardBulk = b
			break
		}
	}

	resp, err := http.Get(cardBulk.DownloadUri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var cs []cards.Card
	err = json.NewDecoder(resp.Body).Decode(&cs)
	if err != nil {
		return err
	}

	for _, c := range cs {
		log.Println(c.Name, c.SetName, c.Set)
	}*/

	return resChan, doneChan, nil
}
