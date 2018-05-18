package delaytask

import (
	"time"

	"github.com/rs/xid"
)

// DTask the delay task
type DTask struct {
	ID          string `json:"id"`
	Payload     string `json:"payload"`
	CreateTime  int64  `json:"create_time"`
	ExecuteTime int64  `json:"execute_time"`
}

// NewDTask create a task
func NewDTask(payload string, d time.Duration) *DTask {
	now := time.Now()
	return &DTask{
		ID:          xid.New().String(),
		Payload:     payload,
		CreateTime:  now.Unix(),
		ExecuteTime: now.Add(d).Unix(),
	}
}
