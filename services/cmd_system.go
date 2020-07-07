package services

import (
	"github.com/geekymedic/x-lite-cli/templates"
	"github.com/geekymedic/x-lite/xerrors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"
)

const (
	systemNameSuffix = "-system"
)

func CreateSystem(systemName string, dirOut string) error {
	dirOut, err := homedir.Expand(dirOut)
	if err != nil {
		return xerrors.By(err)
	}
	dirOut, err = filepath.Abs(dirOut)
	if err != nil {
		return xerrors.By(err)
	}

	systemName += systemNameSuffix

	err = filepath.Walk(dirOut, func(path string, info os.FileInfo, err error) error {
		depth := strings.Count(path, string(filepath.Separator)) - strings.Count(dirOut, string(filepath.Separator))
		if depth > 1 {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == systemName {
			return xerrors.Format("system already exists:%v", systemName)
		}
		return nil
	})

	if err != nil {
		return err
	}

	sysBaseDir := filepath.Join(dirOut, systemName)
	//create bff
	{
		CreateBff(sysBaseDir, "admin")
	}

	//create services
	{
		CreateService(sysBaseDir, "cal")
	}

	//create config
	{
		CreateConfig(sysBaseDir)
	}

	//create gomod
	{
		file := filepath.Join(sysBaseDir, "go.mod")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.GomodTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	//main
	{
		file := filepath.Join(sysBaseDir, "main.go")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.SystemMainTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	//makefile
	{
		file := filepath.Join(sysBaseDir, "Makefile")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.MakeFileTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	return nil
}
