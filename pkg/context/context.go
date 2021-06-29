package context

import (
	gocontext "context"

	"prottohw/pkg/log"
)

type Context struct {
	gocontext.Context
	*log.Logger
}

func Background() Context {
	return Context{
		Context: gocontext.Background(),
		Logger:  log.New(),
	}
}
