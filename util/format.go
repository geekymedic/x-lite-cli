package util

import (
	"github.com/kyokomi/emoji"
	"github.com/geekymedic/x-lite/xerrors"
	"os"
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

func SystemBaseDir() (sysDir string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	sysDir = regexp.MustCompile(".*-system").FindString(dir)
	if sysDir == "" {
		err = xerrors.Format("invalid system directory:%v", dir)
	}
	return
}
