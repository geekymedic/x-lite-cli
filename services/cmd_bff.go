package services

import (
	"github.com/geekymedic/x-lite-cli/templates"
	"github.com/geekymedic/x-lite/xerrors"
	"os"
	"path/filepath"
	"strings"
)

func CreateBff(sysDir string, bffName string) error {
	bffBaseDir := filepath.Join(sysDir, "bff")
	err := os.MkdirAll(bffBaseDir, os.ModePerm)
	if err != nil {
		return xerrors.By(err)
	}

	//create main
	{
		file := filepath.Join(bffBaseDir, "main.go")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return xerrors.By(err)
			}
			defer f.Close()
			err = templates.BffMainTxt.Execute(f, nil)
			if err != nil {
				return xerrors.By(err)
			}
		}
	}

	err = filepath.Walk(bffBaseDir, func(path string, info os.FileInfo, err error) error {
		depth := strings.Count(path, string(filepath.Separator)) - strings.Count(bffBaseDir, string(filepath.Separator))
		if depth > 1 {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == bffName {
			return xerrors.Format("bff already exists:%v", bffName)
		}
		return nil
	})

	if err != nil {
		return err
	}

	bffBaseDir = filepath.Join(bffBaseDir, bffName)
	err = os.Mkdir(bffBaseDir, os.ModePerm)
	if err != nil {
		return xerrors.By(err)
	}

	//create router
	{
		routerDir := filepath.Join(bffBaseDir, "router")
		err := os.Mkdir(routerDir, os.ModePerm)
		if err != nil {
			return xerrors.By(err)
		}
		file := filepath.Join(routerDir, "router.go")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.BffRouterTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	//create impls
	{
		routerDir := filepath.Join(bffBaseDir, "impls", "demo")
		err := os.MkdirAll(routerDir, os.ModePerm)
		if err != nil {
			return xerrors.By(err)
		}
		file := filepath.Join(routerDir, "sum.go")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.BffImplTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}

	//create codes
	{
		codesDir := filepath.Join(bffBaseDir, "codes")
		err := os.MkdirAll(codesDir, os.ModePerm)
		if err != nil {
			return xerrors.By(err)
		}
		file := filepath.Join(codesDir, "error_code.go")
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return xerrors.By(err)
		}
		defer f.Close()
		err = templates.ErrCodeTxt.Execute(f, nil)
		if err != nil {
			return xerrors.By(err)
		}
	}
	return nil
}
