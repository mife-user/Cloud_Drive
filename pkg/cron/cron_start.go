package cron

// Start 启动定时任务
func (c *Cron) Start() {
	c.cron.Start()
}
