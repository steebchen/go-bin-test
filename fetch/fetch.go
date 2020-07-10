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

	fetchCLI(wd)
	fetchEngines(wd)
}

func fetchCLI(wd string) {
	dir := path.Join(wd, "binaries", "cli")

	platforms := []string{
		"darwin",
		"linux",
		"windows",
	}
	for _, platform := range platforms {
		name := platform
		if name == "windows" {
			name = "windows.exe"
		}
		if err := generateGoFile("cli", platform, path.Join(dir, fmt.Sprintf("prisma-cli-%s", name))); err != nil {
			panic(fmt.Errorf("generate go binary file: %w", err))
		}
	}
}

func fetchEngines(wd string) {
	dir := path.Join(wd, "binaries")

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

			if _, err := DownloadEngineByPlatform(engine, binary, path.Join(dir, engine)); err != nil {
				panic(fmt.Errorf("download engine failed: %s", err))
			}
		}
	}
}
