package services

import (
	"errors"
	"github.com/apavlov1992/golang_course/internal/xkcd"
	"github.com/schollz/progressbar/v3"
	"sync"
)

type ComicsWorker struct {
	client  xkcd.Client
	wgCount int
}

type resultWithError struct {
	Comic xkcd.ComicsInfo
	Err   error
}

func NewComicsWorker(client xkcd.Client, wgCount int) *ComicsWorker {
	return &ComicsWorker{
		client:  client,
		wgCount: wgCount,
	}
}

func (w ComicsWorker) handleComic(wg *sync.WaitGroup, inCh <-chan int, resCh chan resultWithError) {
	defer wg.Done()
	for i := range inCh {
		if !xkcd.IDinDB(i) {
			comic, err := w.client.GetComics(i)
			resCh <- resultWithError{
				Comic: comic,
				Err:   err,
			}
		}
	}
}

func (w ComicsWorker) HandleComics() ([]xkcd.ComicsInfo, error) {
	bar := progressbar.NewOptions(-1, progressbar.OptionShowCount(), progressbar.OptionSetDescription("Fetching comics..."))
	var res []xkcd.ComicsInfo
	inputCh := make(chan int)
	wg := &sync.WaitGroup{}
	resCh := make(chan resultWithError)
	quit := make(chan int)
	errCount := 0
	go func() {
		defer close(inputCh)

		for i := 1; ; i++ {
			select {
			case <-quit:
				return
			default:
				bar.Add(1)
				inputCh <- i
			}
		}
	}()

	go func() {
		for i := 1; i < w.wgCount; i++ {
			wg.Add(1)
			go w.handleComic(wg, inputCh, resCh)
		}
		wg.Wait()
		close(resCh)
	}()

	for r := range resCh {
		if r.Err != nil && errors.Is(r.Err, xkcd.NotFound) {
			errCount++
			if errCount == 2 {
				close(quit)
			}
		}
		res = append(res, r.Comic)
	}
	return res, nil
}
