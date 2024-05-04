package app

import (
	"github.com/distributed-calendar/calendar-server/internal/adapter/cache"
	"github.com/distributed-calendar/calendar-server/internal/adapter/timepad"
)

func (a *App) initAdapters() {
	a.initCacheAdapter()
	a.initTimepadAdapter()
}

func (a *App) initCacheAdapter() {
	var err error
	a.cacheAdapter, err = cache.NewAdapter(a.cfg.Redis.Addrs, a.cfg.Redis.Password, a.cfg.Redis.CertPath)

	if err != nil {
		panic(err)
	}
}

func (a *App) initTimepadAdapter() {
	a.timepadAdapter = timepad.NewAdapter(a.cfg.Timepad.URL, a.cfg.Timepad.Token)
}
