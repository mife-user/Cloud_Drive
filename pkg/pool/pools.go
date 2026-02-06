package pool

//定义协程池
type Pool struct {
	Size  int         // 协程池大小
	Tasks chan func() // 任务队列
}

//初始化协程池
func NewPool(size int) *Pool {
	return &Pool{
		Size:  size,
		Tasks: make(chan func(), size),
	}
}

//传入任务并启动协程池
func (p *Pool) Start() {
	for i := 0; i < p.Size; i++ {
		go func() {
			for task := range p.Tasks {
				task()
			}
		}()
	}
}

//传入任务到协程池
func (p *Pool) Submit(task func()) {
	p.Tasks <- task
}

//关闭协程池
func (p *Pool) Stop() {
	close(p.Tasks)
}
