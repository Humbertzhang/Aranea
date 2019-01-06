package core

import (
	"errors"
)

type Job struct {
	*Crawler					`json:"crawler"`
	Id			string			`json:"id"`
	Name		string			`json:"name"`
	// 每个Node在多久之内可以访问一次
	Delay		int				`json:"delay"`
}

type JobQueue struct {
	Queue 		[]*Job
	JobCounter	int
}

func (jobQueue *JobQueue) PushJob(job *Job) {
	jobQueue.Queue = append(jobQueue.Queue, job)
	jobQueue.JobCounter += 1
}

func (jobQueue *JobQueue) PopJob() (job *Job, err error) {
	if(jobQueue.JobCounter <= 0) {
		return nil, errors.New("Job Queue Empty.")
	}
	job = jobQueue.Queue[0]
	jobQueue.Queue = jobQueue.Queue[1:]
	jobQueue.JobCounter -= 1
	return job, nil
}