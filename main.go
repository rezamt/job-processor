package main

import (
	"time"
	"sync"
	"log"
)

type Job struct {
	ID int
	value int
}

type Result struct {
	workerId int
	job Job
	value int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)

func processor (value int) int {

	time.Sleep(1 * time.Second)
	return value * time.Now().Second()
}

func worker (id int, wg *sync.WaitGroup) {
	log.Printf("Worker #%v stated to pool jobs", id )
	for job := range jobs {
		result := Result {id, job, processor(job.value)}
		results <- result
	}

	wg.Done()
}

func createWorkerPool(poolSize int) {
	var wg sync.WaitGroup

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	close(results)
}

func produceJobs(numberOfJobs int) {
	for i:=0; i < numberOfJobs; i++ {
		jobs <- Job {i, i}
	}

	close(jobs)
}

func displayResults(terminate chan bool) {
	for result := range results {
		log.Printf("Worker #%v processed Job #%v  result is %v", result.workerId, result.job.ID, result.value)
	}

	terminate <- true
}


func main() {

	start := time.Now()

	go createWorkerPool(10)

	terminate := make(chan bool)

	go produceJobs(10)

	go displayResults(terminate)

	 <- terminate

	 end := time.Now()

	 log.Println("Total Processing Time: ", end.Sub(start).Seconds())
}
