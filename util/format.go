package util

import (
	"github.com/geekymedic/x-lite/xerrors"
	"github.com/kyokomi/emoji"
	"os"
	"path/filepath"
	"regexp"
)

func StdoutExit(exitCode int, format string, args ...interface{}) {
	if exitCode == 0 {
		emoji.Fprintf(os.Stdout, ":heavy_check_mark:"+format, args...)
	} else {
		emoji.Fprintf(os.Stdout, ":heavy_multiplication_x:"+format, args...)
	}
	os.Exit(exitCode)
}

func SystemBaseDir() (sysDir string, systemName string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	sysDir = regexp.MustCompile(".*-system").FindString(dir)
	if sysDir == "" {
		err = xerrors.Format("invalid system directory:%v", dir)
	}

	_, systemName = filepath.Split(sysDir)
	return
}
