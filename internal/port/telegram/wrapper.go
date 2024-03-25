package telegram

import "github.com/NicoNex/echotron/v3"

func wrapStateFn[S any](fn func(*echotron.Update, *S) stateFn, state *S) stateFn {
	return func(u *echotron.Update) stateFn {
		return fn(u, state)
	}
}
