package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
)

func main() {
	byteMap := map[string][]byte{"one": []byte("1"), "two": []byte("2"), "three": []byte("3"), "four": []byte("4"), "five": []byte("5"), "six": []byte("6"), "seven": []byte("7"), "eight": []byte("8"), "nine": []byte("9")}

	// Read file
	file, err := os.Open("./input.txt")
	if err != nil {
		slog.Error("cannot read file", "error", err)
		os.Exit(1)
	}

	defer file.Close()

	rx, err := regexp.Compile(`(\d|one|two|three|four|five|six|seven|eight|nine)`)
	if err != nil {
		slog.Error("cannot compile regex", "error", err)
		os.Exit(1)
	}

	sum := 0

	// Scan the file lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		b := scanner.Bytes()
		idx := 0
		var matches [][]byte

		// Sneaky, need to detect value like oneight or sevenine and Go doesn't have overlapping support in regexp
		for {
			m := rx.Find(b[idx:])
			if m == nil {
				break
			}

			matches = append(matches, m)
			idx += 1
		}

		fd := matches[0]
		if len(fd) > 1 {
			fd = byteMap[string(fd)]
		}

		ed := matches[len(matches)-1]
		if len(ed) > 1 {
			ed = byteMap[string(ed)]
		}

		fd = append(fd, ed...)
		val, _ := strconv.Atoi(string(fd))

		sum += val
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanning file failed", "error", err)
		os.Exit(1)
	}

	fmt.Printf("\nSum: %+v\n\n", sum)
}
