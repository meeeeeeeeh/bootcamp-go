/*
три взаимоисключающих (можно указать только один за раз,

иначе будет выведено сообщение об ошибке) флага для вашего кода:
-l для подсчета строк,
-m для подсчета символов
-w для подсчета слов - дефолтное поведение
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func getInfo() (uint8, []string, error) {
	info := os.Args[1:]
	if len(info) == 0 {
		return 0, nil, fmt.Errorf("invalid input: some argumends are missing")
	}

	var flag uint8
	var fileName []string
	var isFlag bool
	for _, val := range info {
		if val[0] == '-' && isFlag {
			return 0, nil, fmt.Errorf("invalid input: too many flags")
		}
		if val[0] == '-' {
			flag = val[len(val)-1]
			isFlag = true
		} else {
			fileName = append(fileName, val)
		}
	}

	if flag != 'l' && flag != 'm' && flag != 'w' && flag != 0 {
		return 0, nil, fmt.Errorf("invalid flag")
	}

	return flag, fileName, nil
}

func printInfo(flag uint8, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var counter int
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		switch flag {
		case 'l':
			counter++
		case 'w', 0:
			words := strings.Split(line, " ")
			counter += len(words)
		case 'm':
			for range line {
				counter++
			}
		}
		if err == io.EOF {
			break
		}
	}
	fmt.Printf("%d\t%s\n", counter, fileName)
	return nil
}

func main() {
	flag, files, err := getInfo()
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		go func(f string) {
			defer wg.Done()
			err = printInfo(flag, f)
			if err != nil {
				log.Fatalln(err)
			}
		}(file)

	}
	wg.Wait()
}
