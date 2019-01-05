package core

import (
	"errors"
)

type Job struct {
	*Crawler					`json:"crawler"`
	Id 			string			`json:"id"`
	Name 		string			`json:"name"`		
	// 每个Node在多久之内可以访问一次
	Delay 		int				`json:"delay"`
}

type JobQueue struct {
	Queue 		[]*Job
	JobCounter	int
}

func (jobqueue *JobQueue) PushJob(job *Job) {
	jobqueue.Queue = append(jobqueue.Queue, job)
	jobqueue.JobCounter += 1
}

func (jobqueue *JobQueue) PopJob() (job *Job, err error) {
	if(jobqueue.JobCounter <= 0) {
		return nil, errors.New("Job Queue Empty.")
	}
	job = jobqueue.Queue[0]
	jobqueue.Queue = jobqueue.Queue[1:]
	return job, nil
}