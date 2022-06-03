package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// PingHandler 测试handler
// @Summary 测试Summary
// @Description 测试Description
// @Accept application/json
// @Produce application/json
// @Router /ping [get]
func PingHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(200, map[string]string{
		"ping": "pong",
	})
}
