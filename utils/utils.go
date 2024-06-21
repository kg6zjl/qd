package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Int32Ptr(i int32) *int32 { return &i }

func UniqName(base string) string {
	// generate a unique name
	rand.NewSource(time.Now().UnixNano())
	randInt := rand.Int()
	uniqueId := strconv.Itoa(randInt)

	// verify the unique bit is at least 12 chars or use the whatever we've got
	if len(uniqueId) > 8 {
		uniqueId = uniqueId[:8]
	}

	// format the name
	uniqueName := fmt.Sprintf("%s-qd-%s", base, uniqueId)

	return uniqueName
}

func KubeCpMakeTar(srcPath string, writer *io.PipeWriter) error {
	// Create a new gzip writer
	gw := gzip.NewWriter(writer)
	defer gw.Close()

	// Create a new tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Walk through every file in the folder
	err := filepath.Walk(srcPath, func(file string, fi os.FileInfo, err error) error {
		// Return on any error
		if err != nil {
			return err
		}

		// Create a new dir/file header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// Write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If it's not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
