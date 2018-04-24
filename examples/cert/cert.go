package cert

import (
	"io/ioutil"
	"path"
	"runtime"
)

func getFileContent(fixtureFile string) (content []byte, err error) {
	_, curFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}
	content, err = ioutil.ReadFile(path.Dir(curFile) + "/" + fixtureFile)
	return
}

func Plain(filename string) ([]byte, error) {
	return getFileContent(filename)
}
