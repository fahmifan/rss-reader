package worker

import (
	"log"
	"sync"

	"github.com/miun173/rss-reader/model"
	"github.com/miun173/rss-reader/repository"
)

type rss struct {
	SourceID int64
	Items    []model.RSSItem
}

// Worker :nodoc:
type Worker struct {
	sourceRepo  *repository.SourceRepository
	rssItemRepo *repository.RSSItemRepository
}

// NewWorker :nodoc:
func NewWorker(sr *repository.SourceRepository,
	rir *repository.RSSItemRepository) *Worker {
	return &Worker{
		sourceRepo:  sr,
		rssItemRepo: rir,
	}
}

// FetchRSS :nodoc:
func (w *Worker) FetchRSS() {
	var (
		limit  int64 = 4
		offset int64 = 0
	)

	for {
		log.Println("find sources")
		sources, err := w.sourceRepo.FindAll(limit, offset)
		if err != nil {
			break
		}

		if len(sources) == 0 {
			log.Println("finished")
			break
		}

		w.fetchItems(sources)
		offset += limit
	}
}

// fetchItems fetch rss items & saved it to db
func (w *Worker) fetchItems(sources []model.Source) {
	rssCh := make(chan rss, len(sources))
	wg := sync.WaitGroup{}
	for _, src := range sources {
		wg.Add(1)
		log.Println("fetch rss items from ", src.Name)

		go func(sourceID int64, url string) {
			defer wg.Done()

			items, err := w.rssItemRepo.FetchFromSource(url)
			if err != nil {
				log.Println("error: ", err)
				return
			}

			rssCh <- rss{
				SourceID: sourceID,
				Items:    items,
			}
		}(src.ID, src.URL)
	}

	wg.Wait()
	close(rssCh)

	for rss := range rssCh {
		err := w.rssItemRepo.SaveMany(rss.SourceID, rss.Items)
		if err != nil {
			log.Println(err)
		}
	}
}
