package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	binaries := []string{
		"darwin",
		"linux-musl",
		"windows",
		"debian-openssl-1.0.x",
		"debian-openssl-1.1.x",
		"rhel-openssl-1.0.x",
		"rhel-openssl-1.1.x",
	}

	for _, binary := range binaries {
		if _, err := DownloadEngineByPlatform("query-engine", binary, binary, path.Join(wd, "binaries", "query-engine")); err != nil {
			panic(fmt.Errorf("download engine failed: %s", err))
		}
	}
}
