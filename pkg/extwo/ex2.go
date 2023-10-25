package extwo

import (
	"sync"

	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

/*
	Product wants to start scraping websistes for informations about job roles, and skills requestsd to build a better test library.
 	Engineering is asked to come up with a first implementeation that can receive urls and store the resutls in a database for furure analysis.
  	The system should be able to work in parallel, handling errors from the scraping process gracefully.
*/

type Scrapper interface {
	Scrape(url string) (ScrapedData, error)
	ScrapeMultiple(urls []string) ([]ScrapedData, error)
}

type ScrapedData struct {
	URL    string
	Title  string
	Body   string
	Links  []string
	Extras map[string]string
}

type CollyScrapper struct {
	collector *colly.Collector
}

func NewCollyScrapper() *CollyScrapper {
	return &CollyScrapper{
		collector: colly.NewCollector(),
	}
}

func (cs *CollyScrapper) Scrape(url string) (ScrapedData, error) {
	var data ScrapedData
	data.URL = url

	cs.collector.OnHTML("title", func(e *colly.HTMLElement) {
		data.Title = e.Text
	})

	cs.collector.OnHTML("body", func(e *colly.HTMLElement) {
		data.Body = e.Text
	})

	cs.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		data.Links = append(data.Links, link)
	})

	err := cs.collector.Visit(url)
	if err != nil {
		return ScrapedData{}, err
	}

	return data, nil
}

func (cs *CollyScrapper) ScrapeMultiple(urls []string) ([]ScrapedData, error) {
	var g errgroup.Group
	var dataLock sync.Mutex

	results := make([]ScrapedData, len(urls))

	for i, url := range urls {
		currentIndex, currentURL := i, url

		g.Go(func() error {
			data, err := cs.Scrape(currentURL)
			if err != nil {
				return err
			}

			dataLock.Lock()
			results[currentIndex] = data
			dataLock.Unlock()

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
