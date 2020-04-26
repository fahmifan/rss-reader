package worker

import (
	"log"

	"github.com/miun173/rss-reader/model"
	"github.com/miun173/rss-reader/repository"
)

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
		limit  int64 = 2
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
	var rssItems []model.RSSItem
	for _, src := range sources {
		log.Println("fetch rss items from ", src.Name)
		items, err := w.rssItemRepo.FetchFromSource(src.URL)
		if err != nil {
			break
		}

		rssItems = append(rssItems, items...)
	}

	err := w.rssItemRepo.SaveMany(rssItems)
	if err != nil {
		log.Println(err)
	}
}
