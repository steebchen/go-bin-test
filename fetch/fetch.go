package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	engines := []string{
		"query-engine",
		"migration-engine",
		"introspection-engine",
		"prisma-fmt",
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

	for _, engine := range engines {
		log.Printf("downloading %s", engine)

		for _, binary := range binaries {
			log.Printf("  %s", binary)

			if _, err := DownloadEngineByPlatform(engine, binary, binary, path.Join(wd, "binaries", engine)); err != nil {
				panic(fmt.Errorf("download engine failed: %s", err))
			}
		}
	}
}
