package repo

import (
	"fmt"
)

type Job struct {
	Movie Movie
	Path  string
}

type ThumbNailGenWorker struct {
	repo     *MovieRepo
	jobQueue chan Job
	done     chan bool
}

func NewThumbNailGenWorker(repo *MovieRepo) *ThumbNailGenWorker {
	return &ThumbNailGenWorker{
		repo:     repo,
		jobQueue: make(chan Job),
		done:     make(chan bool),
	}
}

func (t *ThumbNailGenWorker) Start() {
	go func() {
		for {
			select {
			case job, ok := <-t.jobQueue:
				if !ok {
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
