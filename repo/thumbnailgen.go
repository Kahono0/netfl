package repo

import (
	"fmt"
)

var ThumbNailGenWorkerInstance *ThumbNailGenWorker

type Job struct {
	Movie    Movie
	HostAddr string
	Path     string
}

type ThumbNailGenWorker struct {
	jobQueue chan Job
	done     chan bool
}

func NewThumbNailGenWorker() {
	ThumbNailGenWorkerInstance = &ThumbNailGenWorker{
		jobQueue: make(chan Job),
		done:     make(chan bool),
	}
}

func (t *ThumbNailGenWorker) Start() {
	go func() {
		for {
			select {
			case job := <-t.jobQueue:
				fmt.Printf("Creating thumbnail for %s\n", job.Path)
				job.Movie.CreateThumbnail(job.HostAddr, job.Path)
				Repo.UpdateMovie(job.Movie)
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
