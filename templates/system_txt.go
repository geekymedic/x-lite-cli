package templates

const (
	configYmlTxt = `Name: demo-system-bff
Address: ":8003"`

	gomodTxt = `module demo-system

go 1.14`

	systemMainTxt = `package main

import (
	"github.com/geekymedic/x-lite/logger"
	_ "github.com/geekymedic/x-lite/plugin/xmertics"

	"demo-system/bff"
	"demo-system/services"
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
