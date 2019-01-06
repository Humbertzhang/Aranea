package core

import (
	"fmt"
	"github.com/humbertzhang/Aranea/core/utils"
	"net/http"
	"testing"
)

func TestJobQueue_PushPopJob(t *testing.T) {
	jobQueue := &JobQueue{}
	if jobQueue.JobCounter != 0 {
		t.Error("Job queue counter error")
	}
	N := 10
	for i := 0; i < N; i++ {
		j := &Job{
			Crawler: &Crawler{
				URL:         "http://www.baidu.com",
				Method:      http.MethodGet,
				ContentType: "",
				Cookie:      "",
				Payload:     "",
				Headers:     nil,
			},
			Id:    utils.StringTimeStampNanoSecond(),
			Name:  utils.StringTimeStampNanoSecond(),
			Delay: 0,
		}
		jobQueue.PushJob(j)
		if jobQueue.JobCounter != i+1 {
			t.Error("Job queue counter error")
		}
	}

	for i := 0; i < N; i++ {
		job, err := jobQueue.PopJob()
		if err != nil {
			t.Error("pop error")
		}
		if job.Crawler.Method != http.MethodGet {
			t.Error("job poped error")
		}
		if jobQueue.JobCounter != N-i-1 {
			fmt.Println("job counter:", jobQueue.JobCounter, " n-i:", N-i)
			t.Error("job number pop error")
		}
	}
}



