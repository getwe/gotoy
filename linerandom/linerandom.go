package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage : cat xxx | linerandom N")
		return
	}

	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("%s not an valid number\n", os.Args[1])
		return
	}

	result := make([]string, N)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	i := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		idx := 0
		if i < N {
			idx = i
		} else {
			idx = r.Intn(i)
		}

		if idx < N {
			result[idx] = line
		}
		i++
	}

	for _, line := range result {
		fmt.Println(line)
	}
}
