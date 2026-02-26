package cron

import "github.com/robfig/cron/v3"

// AddFunc 添加定时任务
func (c *Cron) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return c.cron.AddFunc(spec, cmd)
}
