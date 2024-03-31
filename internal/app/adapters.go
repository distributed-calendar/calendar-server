package app

import "github.com/distributed-calendar/calendar-server/internal/adapter/cache"

func (a *App) initAdapters() {
	a.initCacheAdapter()
}

func (a *App) initCacheAdapter() {
	a.cacheAdapter = cache.NewAdapter(a.cfg.Redis.Addr, a.cfg.Redis.Password)
}
