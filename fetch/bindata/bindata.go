package bindata

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TODO go fmt files after creation

// WriteFile writes the release code file.
func WriteFile(w io.Writer, name, platform, file string) error {
	if err := writeHeader(w, name, platform); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	if err := writeAsset(w, file); err != nil {
		return fmt.Errorf("write asset: %w", err)
	}

	return nil
}

// writeHeader writes output file headers.
// This targets release builds.
func writeHeader(w io.Writer, name, platform string) error {
	var excludes string
	if platform == "linux-musl" {
		platform = "linux"
		// TODO dynamically construct these
		// TODO only include these for engines, not for the CLi :D
		excludes = `
// +build !debian_openssl_1_0_x
// +build !debian_openssl_1_1_x
// +build !rhel_openssl_1_0_x
// +build !rhel_openssl_1_1_x
		`
	}

	platform = strings.Replace(platform, "-", "_", -1)
	platform = strings.Replace(platform, ".", "_", -1)

	_, err := fmt.Fprintf(w, `
// +build %s
// +build !prisma_ignore

%s

package binaries

import (
	"github.com/steebchen/go-binaries/unpack"
)

func init() {
	unpack.Unpack(data, "%s-%s")
}

`, platform, excludes, name, platform)
	return err
}

// writeAsset write a release entry for the given asset.
// A release entry is a function which embeds and returns
// the file's byte content.
func writeAsset(w io.Writer, file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	h := sha256.New()
	tr := io.TeeReader(fd, h)
	if err := uncompressedMemcopy(w, tr); err != nil {
		return err
	}
	var digest [sha256.Size]byte
	copy(digest[:], h.Sum(nil))
	return assetMetadata(w, file, digest)
}

var (
	backquote = []byte("`")
	bom       = []byte("\xEF\xBB\xBF")
)

// sanitize prepares a valid UTF-8 string as a raw string constant.
// Based on https://code.google.com/p/go/source/browse/godoc/static/makestatic.go?repo=tools
func sanitize(b []byte) []byte {
	var chunks [][]byte
	for i, b := range bytes.Split(b, backquote) {
		if i > 0 {
			chunks = append(chunks, backquote)
		}
		for j, c := range bytes.Split(b, bom) {
			if j > 0 {
				chunks = append(chunks, bom)
			}
			if len(c) > 0 {
				chunks = append(chunks, c)
			}
		}
	}

	var buf bytes.Buffer
	sanitizeChunks(&buf, chunks)
	return buf.Bytes()
}

func sanitizeChunks(buf *bytes.Buffer, chunks [][]byte) {
	n := len(chunks)
	if n >= 2 {
		buf.WriteString("(")
		sanitizeChunks(buf, chunks[:n/2])
		buf.WriteString(" + ")
		sanitizeChunks(buf, chunks[n/2:])
		buf.WriteString(")")
		return
	}
	b := chunks[0]
	if bytes.Equal(b, backquote) {
		buf.WriteString("\"`\"")
		return
	}
	if bytes.Equal(b, bom) {
		buf.WriteString(`"\xEF\xBB\xBF"`)
		return
	}
	buf.WriteString("`")
	buf.Write(b)
	buf.WriteString("`")
}

func uncompressedMemcopy(w io.Writer, r io.Reader) error {
	if _, err := fmt.Fprintf(w, `var data = []byte(`); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if len(b) > 0 && validSanitizedUtf8(b) {
		w.Write(sanitize(b))
	} else {
		fmt.Fprintf(w, "%+q", b)
	}

	if _, err := fmt.Fprintf(w, `)`); err != nil {
		return err
	}
	return nil
}

func assetMetadata(w io.Writer, file string, digest [sha256.Size]byte) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}

	mode := uint(fi.Mode())
	modTime := fi.ModTime().Unix()
	size := fi.Size()
	_, err = fmt.Fprintf(w, `
//func Data() (*asset, error) {
//	info := bindataFileInfo{name: "bytes", size: %d, mode: os.FileMode(%#o), modTime: time.Unix(%d, 0)}
//	a := &asset{bytes: data, info: info, digest: %#v}
//	return a, nil
//}

`, size, mode, modTime, digest)
	return err
}

func validSanitizedUtf8(b []byte) bool {
	if !utf8.Valid(b) {
		return false
	}
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == 0 {
			return false
		}
		if unicode.In(r, unicode.Cf) {
			// staticcheck doesn't like these; fallback to slow decoder
			return false
		}
		b = b[size:]
	}
	return true
}
