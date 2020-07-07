package services

import (
	"github.com/geekymedic/x-lite-cli/templates"
	"github.com/geekymedic/x-lite/xerrors"
	"os"
	"path/filepath"
	"strings"
)

func CreateService(sysDir string, serviceName string) error {
	serviceBaseDir := filepath.Join(sysDir, "services")
	err := os.MkdirAll(serviceBaseDir, os.ModePerm)
	if err != nil {
		return xerrors.By(err)
	}

	//create main
	{
		file := filepath.Join(serviceBaseDir, "main.go")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return xerrors.By(err)
			}
			defer f.Close()
			err = templates.ServiceMainTxt.Execute(f, nil)
			if err != nil {
				return xerrors.By(err)
			}
		}
	}

	err = filepath.Walk(serviceBaseDir, func(path string, info os.FileInfo, err error) error {
		depth := strings.Count(path, string(filepath.Separator)) - strings.Count(serviceBaseDir, string(filepath.Separator))
		if depth > 1 {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == serviceName {
			return xerrors.Format("service already exists:%v", serviceName)
		}
		return nil
	})

	if err != nil {
		return err
	}

	serviceImplDir := filepath.Join(serviceBaseDir, serviceName)
	err = os.Mkdir(serviceImplDir, os.ModePerm)
	if err != nil {
		return xerrors.By(err)
	}

	{
		file := filepath.Join(serviceImplDir, "server.go")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.ServiceServerTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	{
		file := filepath.Join(serviceImplDir, "sum.go")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.ServiceImlTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	return nil
}
