package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func findChanges(text string, file1, file2 *os.File) {
	tmpMap := make(map[string]struct{})
	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		tmpMap[scanner.Text()] = struct{}{}
	}

	scanner = bufio.NewScanner(file2)
	for scanner.Scan() {
		line := scanner.Text()
		_, ok := tmpMap[line]
		if !ok {
			fmt.Print(fmt.Sprintf("%s %s\n", text, line))
		}
	}

	file1.Seek(0, 0)
	file2.Seek(0, 0)
}

func getFileNames() (string, string) {
	oldFileName := flag.String("old", "", "old snapshot file name")
	newFileName := flag.String("new", "", "new snapshot file name")
	flag.Parse()
	return *oldFileName, *newFileName
}

func main() {
	oldFileName, newFileName := getFileNames()

	oldFile, err := os.Open(oldFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer oldFile.Close()

	newFile, err := os.Open(newFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	findChanges("ADDED", oldFile, newFile)
	findChanges("REMOVED", newFile, oldFile)
}
