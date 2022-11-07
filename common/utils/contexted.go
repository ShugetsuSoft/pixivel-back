package utils

import (
	"context"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
)

type configCtx struct {
	context.Context
	val *models.Config
}

func ConfigWrapper(ctx context.Context, config *models.Config) *configCtx {
	return &configCtx{
		Context: ctx,
		val:     config,
	}
}

func (c *configCtx) Value(key interface{}) interface{} {
	switch s := key.(type) {
	case string:
		if s == "config" {
			return c.val
		}
	}
	return c.Context.Value(key)
}
