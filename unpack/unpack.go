package unpack

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

// TODO check checksum after expanding file

//noinspection GoUnusedExportedFunction
func Unpack(data []byte) {
	name := "prisma-query-engine-" + runtime.GOOS

	// TODO check if dev env/dev binary in ~/.prisma
	// TODO check if engine in local dir OR env var

	start := time.Now()
	dir := path.Join(".", name)
	// dir := path.Join(os.TempDir(), name)
	unzipped, err := unzip(data)
	if err != nil {
		panic(fmt.Errorf("unpack unzip: %w", err))
	}
	if err := ioutil.WriteFile(dir, unzipped, os.ModePerm); err != nil {
		panic(fmt.Errorf("unpack write file: %w", err))
	}
	log.Printf("unpacked query engine at %s in %s", dir, time.Since(start))
}

func unzip(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	var r io.Reader
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	var res bytes.Buffer
	if _, err = res.ReadFrom(r); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
