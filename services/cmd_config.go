package services

import (
	"github.com/geekymedic/x-lite-cli/templates"
	"github.com/geekymedic/x-lite/xerrors"
	"os"
	"path/filepath"
)

func CreateConfig(sysDir string) error {
	configDir := filepath.Join(sysDir, "config")
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return xerrors.By(err)
	}

	{
		file := filepath.Join(configDir, "config.yml")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return xerrors.By(err)
			}
			defer f.Close()
			err = templates.ConfigYmlTxt.Execute(f, nil)
			if err != nil {
				return xerrors.By(err)
			}
		}
	}

	return nil
}
