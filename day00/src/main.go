package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

//flags: -mode -median -mean -sd
// если были введены флаги то выводятся только определенные значения
// по умолчанию - все

// среднее значение
func mean(val []float64) float64 {
	var sum float64
	for _, val := range val {
		sum += val
	}
	return sum / float64(len(val))
}

// медиана
func median(val []float64) float64 {
	middle := len(val) / 2

	sort.Slice(val, func(i, j int) bool {
		return val[i] < val[j]
	})

	if len(val)%2 != 0 {
		return val[middle]
	}
	return (val[middle] + val[middle-1]) / 2

}

// мода
func mode(val []float64) float64 {
	m := make(map[float64]int)
	var res float64
	var amount int

	for _, v := range val {
		_, ok := m[v]
		if !ok {
			m[v] = 1
		} else {
			m[v] += 1
		}
	}

	for k, v := range m {
		if v > amount {
			res = k
			amount = v
		} else if v == amount {
			if k < res {
				res = k
			}
		}
	}

	return res
}

// среднеквадратичное отклонение
func sd(val []float64) float64 {
	mean := mean(val)
	var devSum, add float64
	for _, v := range val {
		add = v - mean
		if add < 0 {
			add = add * -1
		}
		devSum += add
	}
	return devSum / float64(len(val))
}

func getValues() []float64 {
	var values []float64
	var value string

	fmt.Println("Input all values and separate them by enter. To start the calculation write 'run'")

	for {
		fmt.Scan(&value)
		if value == "run" {
			break
		}
		convert, err := strconv.ParseFloat(value, 2)
		if err != nil {
			log.Panic()
		}
		values = append(values, convert)

	}
	return values
}

func flagCheck(flags []string) error {
	for _, val := range flags {
		if val == "-mode" || val == "-sd" || val == "-median" || val == "-mean" {
			continue
		}
		err := errors.New("invalid flag")
		return err
	}
	return nil
}

func main() {
	modules := os.Args[1:]
	err := flagCheck(modules)
	if err != nil {
		log.Panic()
	}

	values := getValues()

	if len(modules) != 0 {
		for _, val := range modules {
			if val == "-mean" {
				fmt.Printf("Mean: %f\n", mean(values))
			} else if val == "-median" {
				fmt.Printf("Median: %f\n", median(values))
			} else if val == "-mode" {
				fmt.Printf("Mode: %f\n", mode(values))
			} else if val == "-sd" {
				fmt.Printf("SD: %f\n", sd(values))
			}
		}
	} else {
		fmt.Printf("Mean: %f\n", mean(values))
		fmt.Printf("Median: %f\n", median(values))
		fmt.Printf("Mode: %f\n", mode(values))
		fmt.Printf("SD: %f\n", sd(values))
	}
}
