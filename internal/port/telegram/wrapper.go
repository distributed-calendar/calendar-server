package telegram

import "github.com/NicoNex/echotron/v3"

type handlerOpt func(fn stateFn) stateFn

func withOpts(fn stateFn, opts ...handlerOpt) stateFn {
	if len(opts) == 0 {
		return fn
	}

	return withOpts(opts[0](fn), opts[1:]...)
}

func withState[S any](fn func(*echotron.Update, *S) stateFn, state *S) stateFn {
	return func(u *echotron.Update) stateFn {
		return fn(u, state)
	}
}

func (b *botAPI) withCancel(fn stateFn) stateFn {
	return func(u *echotron.Update) stateFn {
		if u.Message != nil && u.Message.Text == commandCancel {

			return b.handleDefault
		}

		return fn(u)
	}
}
