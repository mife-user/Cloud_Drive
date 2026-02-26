package bootstrap

import (
	"drive/internal/api/routes"
	"fmt"
)

// initRouter 初始化路由
func (a *Application) initRouter() error {
	router := routes.GetRouter()
	if !router.NewRouter(a.Database, a.Redis, a.Config) {
		return fmt.Errorf("创建路由失败")
	}
	a.Engine = router.Setup()
	return nil
}
