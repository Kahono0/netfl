package movies

import (
	"fmt"
	"sync"
)

const MAX_PROCESS = 5

type Job struct {
	Movie Movie
	Path  string
}

type ThumbNailGenWorker struct {
	repo     *MovieRepo
	jobQueue chan Job
	done     chan bool
	wg       *sync.WaitGroup
}

func NewThumbNailGenWorker(repo *MovieRepo) *ThumbNailGenWorker {
	worker := &ThumbNailGenWorker{
		repo:     repo,
		jobQueue: make(chan Job),
		done:     make(chan bool),
		wg:       &sync.WaitGroup{},
	}

	for i := 0; i < MAX_PROCESS; i++ {
		worker.wg.Add(1)
		go worker.Start()
	}

	return worker
}

func (t *ThumbNailGenWorker) Start() {
	go func() {
		for {
			select {
			case job, ok := <-t.jobQueue:
				if !ok {
					t.done <- true
					return
				}
				fmt.Printf("Creating thumbnail for %s\n", job.Path)
				job.Movie.CreateThumbnail(t.repo.HostAddr, job.Path, t.repo.Dir)
				t.repo.UpdateMovie(job.Movie)
				fmt.Printf("Thumbnail created for %s\n", job.Path)
			case <-t.done:
				return
			}
		}
	}()
}

func (t *ThumbNailGenWorker) Stop() {
	t.done <- true
}

func (t *ThumbNailGenWorker) AddJob(job Job) {
	t.jobQueue <- job
}

func (t *ThumbNailGenWorker) Close() {
	close(t.jobQueue)
}
