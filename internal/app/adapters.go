package app

import "github.com/distributed-calendar/calendar-server/internal/adapter/cache"

func (a *App) initAdapters() {
	a.initCacheAdapter()
}

func (a *App) initCacheAdapter() {
	var err error
	a.cacheAdapter, err = cache.NewAdapter(a.cfg.Redis.Addrs, a.cfg.Redis.Password, a.cfg.Redis.CertPath)

	if err != nil {
		panic(err)
	}
}
