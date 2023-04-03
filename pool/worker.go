package pool

import (
	"math/rand"
)

type Worker struct {
	id       int
	jobQueue chan Job
	workers  chan chan Job
	quit     chan bool
}

func NewWorker(workers chan chan Job) *Worker {
	return &Worker{
		id:       rand.Int(),
		jobQueue: make(chan Job),
		workers:  workers,
		quit:     make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// 将自己注册到 worker 池中
			w.workers <- w.jobQueue

			select {
			case job := <-w.jobQueue: // 收到任务后执行，并通知任务完成
				job.Do()
			case <-w.quit: // 收到停止信号后，退出循环并关闭该 worker
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.quit <- true
}
