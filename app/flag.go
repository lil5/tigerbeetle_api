package app

import (
	"path"
	"strings"
)

func ReadFlag(fl string) (fpath, name, ext string) {
	fpath = path.Dir(fl)
	configFname := strings.TrimPrefix(fl, fpath)
	ext = strings.TrimPrefix(path.Ext(fl), ".")
	name = strings.TrimPrefix(strings.TrimSuffix(configFname, "."+ext), "/")

	return
}
