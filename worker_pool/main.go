package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

//指定最大worker数
var maxWorker = 6

//任务队列长度
var maxJobQueue = 8

//任务队列
var JobQueue chan Job = make(chan Job, maxJobQueue)

//任务
type Job struct {
	Data interface{}
	Proc func(interface{})
}

//worker 工作者模型
type Worker struct {
	WorkerPool chan chan Job
	JobChan    chan Job
	QuitChan   chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChan:    make(chan Job),
		QuitChan:   make(chan bool),
	}
}

//worker处理任务
func (worker Worker) Start() {
	go func() {
		for {
			worker.WorkerPool <- worker.JobChan
			select {
			case job := <-worker.JobChan:
				job.Proc(job.Data)
			case <-worker.QuitChan:
				//收到退出信号
				return
			}
		}
	}()
}

//Dispatcher
type Dispatcher struct {
	MaxWorker  int
	WorkerPool chan chan Job
}

func NewDispatcher() Dispatcher {
	//初始化dispatcher
	workerPool := make(chan chan Job, maxWorker)
	return Dispatcher{
		MaxWorker:  maxWorker,
		WorkerPool: workerPool,
	}
}

func (dispatcher *Dispatcher) Run() {
	//启动worker
	for i := 0; i < dispatcher.MaxWorker; i++ {
		worker := NewWorker(dispatcher.WorkerPool)
		worker.Start()
	}
}

func (dispatcher *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			//从任务队列取
			go func(job Job) {
				//获取一个可用的worker
				jobChan := <-dispatcher.WorkerPool
				jobChan <- job
			}(job)
		}

		time.Sleep(time.Second)
	}
}

//启动任务测试协程
func generateJobGroutine() {
	jobID := 1
	for {
		JobQueue <- Job{
			Data: fmt.Sprintf("job_%d", jobID),
			Proc: func(d interface{}) {
				beego.Info("Job info:", d.(string))
			},
		}

		jobID++
	}
}

func main() {
	go generateJobGroutine()

	dispatcherInstance := NewDispatcher()
	dispatcherInstance.Run()
	dispatcherInstance.dispatch()
}
