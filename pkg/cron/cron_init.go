package cron

import (
	"github.com/robfig/cron/v3"
)

type Cron struct {
	cron *cron.Cron
}

// NewCron 初始化定时任务
func NewCron() *Cron {
	return &Cron{
		cron: cron.New(),
	}
}
