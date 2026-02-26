package bootstrap

import "drive/pkg/res"

// initRedis 初始化Redis
func (a *Application) initRedis() error {
	if err := res.Init(); err != nil {
		return err
	}
	a.Redis = res.GetRD()
	return nil
}
