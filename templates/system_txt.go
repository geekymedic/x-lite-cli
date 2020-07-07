package templates

const (
	configYmlTxt = `Name: {{.SystemName}}-bff
Address: ":8003"`

	gomodTxt = `module {{.SystemName}}

go 1.14`

	systemMainTxt = `package main

import (
	"github.com/geekymedic/x-lite/logger"
	_ "github.com/geekymedic/x-lite/plugin/xmertics"

	"{{.SystemName}}/bff"
	"{{.SystemName}}/services"
)

func main() {
	if err := services.Main(); err != nil {
		logger.Error(err)
	}
	if err := bff.Main(); err != nil {
		logger.Error(err)
	}
}`
)
