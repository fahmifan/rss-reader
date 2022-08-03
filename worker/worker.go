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
		limit  int = 1
		offset int = 0
	)
	const nworker = 4

	sourcesChan := make(chan []model.Source, nworker)
	rssChan := make(chan rss, nworker)

	wgFetch := sync.WaitGroup{}
	wgFetch.Add(nworker)
	for i := 0; i < nworker; i++ {
		go w.workFetch(i, &wgFetch, sourcesChan, rssChan)
	}

	// single db writer
	wgRss := sync.WaitGroup{}
	wgRss.Add(1)
	go w.workSave(&wgRss, rssChan)

	log.Println("find sources")
	for {
		sources, err := w.sourceRepo.FindAll(limit, offset)
		if err != nil {
			break
		}

		if len(sources) == 0 {
			break
		}

		sourcesChan <- sources
		offset += limit
	}
	close(sourcesChan)

	wgFetch.Wait()
	close(rssChan)
	wgRss.Wait()

	log.Println("finished")
}

func (w *Worker) workFetch(id int, wg *sync.WaitGroup, srcsChan chan []model.Source, resultChan chan rss) {
	defer wg.Done()

	for srcs := range srcsChan {
		for _, src := range srcs {
			resultChan <- w.fetchItem(src)
			log.Printf("#%d success fetch rss items from %s\n", id, src.URL)
		}
	}
}

func (w *Worker) workSave(wg *sync.WaitGroup, rssChan chan rss) {
	defer wg.Done()
	for rss := range rssChan {
		err := w.rssItemRepo.SaveMany(rss.SourceID, rss.Items)
		if err != nil {
			log.Println(err)
		}
	}
}

func (w *Worker) fetchItem(src model.Source) rss {

	items, err := w.rssItemRepo.FetchFromSource(src.URL)
	if err != nil {
		log.Println("error: ", err)
		return rss{}
	}

	return rss{
		SourceID: src.ID,
		Items:    items,
	}
}
