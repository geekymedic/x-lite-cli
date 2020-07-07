package templates

const (
	bffImplTxt = `package demo

import (
	"github.com/geekymedic/x-lite/framework/bff"
	"github.com/geekymedic/x-lite/framework/bff/codes"
)

// @type: b.i.rt
// @interface: SumHandler
// @describe:
type SumRequest struct {
	N int
}

// @type: b.i.re
// @interface: SumHandler
// @describe:
type SumResponse struct {
	Ret int
}

// @type: b.i
// @name: 获取设备信息
// @login: Y
// @version: 0.0.1
// @page:
// @uri: /api/admin/v1/{{.BffName}}/demo/sum
// @describe:
func SumHandler(state *bff.State) {
	var (
		ask = &SumRequest{}
		ack = &SumResponse{}
	)
	if err := state.ShouldBindJSON(ask); err != nil {
		state.Error(codes.CodeRequestBody, err)
		return
	}
	state.Success(ack)
}`

	errCodeTxt = `package codes

import "github.com/geekymedic/neon/bff"

const (
	CodeLess = 8000
)

var (
	_codes = bff.Codes{
		CodeLess: "小于",
	}
)

func GetMessage(code int) string {
	return _codes[code]
}

func init() {
	bff.MergeCodes(_codes)
}
`

	bffRouterTxt = `package router

import (
	"github.com/geekymedic/x-lite/framework/bff"

	"{{.SystemName}}/bff/{{.BffName}}/impls/demo"
)

func init() {
	group := bff.Engine().Group("/admin/v1")
	{
		group := group.Group("/{{.BffName}}/demo")
		group.POST("/sum", bff.HttpHandler(demo.SumHandler))
	}
}`

	bffMainTxt = `package bff

import (
	"github.com/geekymedic/x-lite/framework/bff"
	_ "github.com/geekymedic/x-lite/plugin/xvalidate"
	_ "{{.SystemName}}/bff/admin/router"
)

func Main() error {
	return bff.Main()
}`
)
