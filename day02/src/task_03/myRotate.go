package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func archiveFile(filePath, archiveDir string) error {
	info, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("cannot get file info, error: %v", err)
	}

	timestamp := info.ModTime().Unix()
	archiveFileName := fmt.Sprintf("%s_%d.tar.gz", filePath, timestamp)
	if archiveDir != "" {
		archiveFileName = filepath.Join(archiveDir, archiveFileName)
	}

	archiveFile, err := os.Create(archiveFileName)
	if err != nil {
		return fmt.Errorf("cannot create archive file, error: %v", err)
	}
	defer archiveFile.Close()

	gzipWriter := gzip.NewWriter(archiveFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file, error: %v", err)
	}
	defer file.Close()

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("error creating tar header for file %s, error: %v", filePath, err)
	}
	header.Name = filepath.Base(filePath)

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("error writing tar header for file %s, error: %v", filePath, err)
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return fmt.Errorf("error writing file %s to tar, error: %v", filePath, err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error removing log file %s, error: %v", filePath, err)
	}
	return nil
}

func main() {
	var archiveDir string
	flag.StringVar(&archiveDir, "a", "", "directory to store archives")
	flag.Parse()
	files := flag.Args()

	if len(files) == 0 {
		log.Fatalln("file name is missing")
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	// Запускаем архивацию файлов в горутинах
	for _, file := range files {
		go func(f string) {
			defer wg.Done()
			err := archiveFile(f, archiveDir)
			if err != nil {
				log.Fatalf("cannot archive file: %v", err)
			}
		}(file)
	}
	wg.Wait()
}
