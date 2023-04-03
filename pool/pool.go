package pool

type WorkerPool struct {
	workers   chan chan Job
	jobs      chan Job
	Quit      chan bool
	maxWorker int
}

func NewWorkerPool(maxWorker int) *WorkerPool {
	pool := &WorkerPool{
		workers:   make(chan chan Job, maxWorker),
		jobs:      make(chan Job),
		Quit:      make(chan bool),
		maxWorker: maxWorker,
	}
	return pool
}

func (p *WorkerPool) Start() {
	// 初始化 worker
	for i := 0; i < p.maxWorker; i++ {
		worker := NewWorker(p.workers)
		worker.Start()
	}

	go p.dispatch()
}

func (p *WorkerPool) Stop() {
	close(p.Quit)
}

func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

func (p *WorkerPool) dispatch() {
	for {
		select {
		case job := <-p.jobs:
			// 获取任意一个空闲的 worker，并将任务分配给它
			workerJob := <-p.workers
			workerJob <- job
		case <-p.Quit:
			// 收到停止信号后，关闭所有 worker 并退出循环
			for i := 0; i < p.maxWorker; i++ {
				p.workers <- nil
			}
			return
		}
	}
}
