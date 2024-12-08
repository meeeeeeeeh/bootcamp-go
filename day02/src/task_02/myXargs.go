package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("invalid input: some arguments are missing")
	}

	command := os.Args[1]
	args := os.Args[2:]

	var arg string
	for {
		_, err := fmt.Scan(&arg)
		if err != nil {
			break
		}
		args = append(args, arg)
	}

	fmt.Println(command, args)

	cmd := exec.Command(command, args...)
	// Перенаправляем стандартный вывод команды на стандартный вывод текущего процесса
	//По умолчанию, команда exec.Command не наследует стандартный вывод родительского процесса,
	//поэтому вывод команды необходимо явно перенаправить.

	cmd.Stdout = os.Stdout
	// Перенаправляем стандартный поток ошибок команды на стандартный поток ошибок текущего процесса
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalln("command execution failed")
	}
}
