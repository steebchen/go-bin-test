package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/steebchen/go-binaries/fetch/bindata"
	"github.com/steebchen/go-binaries/fetch/platform"
)

// EngineVersion is a hardcoded version of the Prisma Engine.
// The versions can be found under https://github.com/prisma/prisma-engine/commits/master.
const EngineVersion = "4dbefe724ae6df3621eb660a5d0e0402e34189bf"

// EngineURL points to an S3 bucket URL where the Prisma engines are stored.
const EngineURL = "https://binaries.prisma.sh/master/%s/%s/%s.gz"

func DownloadEngineByPlatform(name string, binaryName string, toDir string) (file string, err error) {
	to := platform.CheckForExtensionByPlatform(binaryName, path.Join(toDir, fmt.Sprintf("prisma-%s-%s", name, binaryName)))

	url := platform.CheckForExtensionByPlatform(binaryName, fmt.Sprintf(EngineURL, EngineVersion, binaryName, name))

	if err := download(name, binaryName, url, to); err != nil {
		return "", fmt.Errorf("could not download %s to %s: %w", url, to, err)
	}

	return to, nil
}

func download(name, platform, url, to string) error {
	if err := os.MkdirAll(path.Dir(to), os.ModePerm); err != nil {
		return fmt.Errorf("could not run MkdirAll on path %s: %w", to, err)
	}

	// copy to temp file first
	dest := to + ".tmp"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not get %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received code %d from %s: %+v", resp.StatusCode, url, string(out))
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", dest, err)
	}
	defer out.Close()

	if err := os.Chmod(dest, os.ModePerm); err != nil {
		return fmt.Errorf("could not chmod +x %s: %w", url, err)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("could not copy %s: %w", url, err)
	}

	// temp file is ready, now copy to the original destination
	if err := copyFile(dest, to+".gz"); err != nil {
		return fmt.Errorf("copy temp file: %w", err)
	}

	if err := generateGoFile(name, platform, to); err != nil {
		return fmt.Errorf("generate go binary file: %w", err)
	}

	// clean up temp file
	_ = os.Remove(dest)

	return nil
}

func generateGoFile(name, platform, file string) error {
	f, err := os.Create(file + ".go")
	if err != nil {
		return fmt.Errorf("generate open go file: %w", err)
	}

	if err := bindata.WriteFile(f, name, platform, file+".gz"); err != nil {
		return fmt.Errorf("generate write go file: %w", err)
	}

	return nil
}

func copyFile(from string, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return fmt.Errorf("readfile: %w", err)
	}

	err = ioutil.WriteFile(to, input, os.ModePerm)
	if err != nil {
		return fmt.Errorf("writefile: %w", err)
	}

	return nil
}
