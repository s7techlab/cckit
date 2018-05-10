package cert

import (
	"errors"
	"io/ioutil"
	"path"
	"runtime"
)

func Content(fixtureFile string) ([]byte, error) {
	_, curFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New(`can't load file, error accessing runtime caller'`)
	}
	return ioutil.ReadFile(path.Dir(curFile) + "/" + fixtureFile)
}
