package services

import (
	"github.com/geekymedic/x-lite-cli/util"
	"github.com/geekymedic/x-lite/xerrors"
	"os"
	"path/filepath"
	"strings"
)

func CreateMd(sysDir string, bffName string, implName string) error {
	bffAndImpls := map[string][]string{}

	bffBaseDir := filepath.Join(sysDir, "bff")
	filepath.Walk(bffBaseDir, func(path string, info os.FileInfo, err error) error {
		depth := strings.Count(path, string(filepath.Separator)) - strings.Count(bffBaseDir, string(filepath.Separator))
		if depth > 1 {
			return filepath.SkipDir
		} else if depth == 0 {
			return nil
		}
		if bffName == "" && info.IsDir() {
			bffAndImpls[info.Name()] = []string{}
		} else if bffName == info.Name() && info.IsDir() {
			bffAndImpls[info.Name()] = []string{}
		}
		return nil
	})

	for bffName := range bffAndImpls {
		implsDir := filepath.Join(bffBaseDir, bffName, "impls")
		filepath.Walk(implsDir, func(path string, info os.FileInfo, err error) error {
			depth := strings.Count(path, string(filepath.Separator)) - strings.Count(implsDir, string(filepath.Separator))
			if depth > 1 {
				return filepath.SkipDir
			} else if depth == 0 {
				return nil
			}
			if implName == "" && info.IsDir() {
				bffAndImpls[bffName] = append(bffAndImpls[bffName], info.Name())
			} else if implName == info.Name() && info.IsDir() {
				bffAndImpls[bffName] = append(bffAndImpls[bffName], info.Name())
			}
			return nil
		})
	}

	//generate
	for bffName, impls := range bffAndImpls {
		err := util.CreateMd(sysDir, bffName, impls)
		return xerrors.By(err)
	}

	return nil
}
