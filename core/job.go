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
	FailedTimes int
}

type JobQueue struct {
	Queue 		chan *Job
	JobCounter	int
}

func (jobQueue *JobQueue) PushJob(job *Job) {
	jobQueue.Queue <- job
	jobQueue.JobCounter += 1
}

func (jobQueue *JobQueue) PopJob() (job *Job, err error) {
	if(jobQueue.JobCounter <= 0) {
		return nil, errors.New("Job Queue Empty.")
	}
	job = <-jobQueue.Queue
	jobQueue.JobCounter -= 1
	return job, nil
}