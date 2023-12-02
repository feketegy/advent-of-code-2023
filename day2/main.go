package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	maxCubes := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	// Read file
	file, err := os.Open("./input.txt")
	if err != nil {
		slog.Error("cannot read file", "error", err)
		os.Exit(1)
	}

	defer file.Close()

	// Compile the regex that will get the game ID and the rest of the line
	rx, err := regexp.Compile(`^Game\s([\d]{1,3}):\s(.{1,})`)
	if err != nil {
		slog.Error("cannot compile main regex", "error", err)
		os.Exit(1)
	}

	rxCubes, err := regexp.Compile(`([\d]{1,3})\s(red|green|blue)`)
	if err != nil {
		slog.Error("cannot compile regex rxThrows", "error", err)
		os.Exit(1)
	}

	gameIDSum := 0
	sumOfPowers := 0

	// Scan the file lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		matches := rx.FindStringSubmatch(scanner.Text())

		// Store the game ID
		gameID, _ := strconv.Atoi(matches[1])
		fmt.Printf("\n====== Game ID: %+v", gameID)

		goodPlay := true

		// Split the game groups on each line by ;
		throws := strings.Split(matches[2], ";")

		// Part 2
		minCubes := map[string]int{
			"red":   0,
			"green": 0,
			"blue":  0,
		}

		// Loop over each throw and split the cubes
		for _, g := range throws {
			fmt.Printf("\n  - throw: %+v", g)

			cubes := rxCubes.FindAllStringSubmatch(g, -1)
			for _, t := range cubes {
				fmt.Printf("\n    - nr: %+v, color: %+v", t[1], t[2])

				nrCubes, _ := strconv.Atoi(t[1])
				color := strings.ToLower(t[2])

				// Part 2
				if minCubes[color] < nrCubes {
					minCubes[color] = nrCubes
				}

				// Part 1
				if nrCubes > maxCubes[color] {
					goodPlay = false
				}
			}
		}

		fmt.Printf("\n  - min cubes: red: %+v, green: %+v, blue: %+v", minCubes["red"], minCubes["green"], minCubes["blue"])

		powerCubes := minCubes["red"] * minCubes["green"] * minCubes["blue"]
		sumOfPowers += powerCubes

		fmt.Printf("\n  - power: %+v\n", powerCubes)

		if goodPlay == true {
			gameIDSum += gameID
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanning file failed", "error", err)
		os.Exit(1)
	}

	fmt.Printf("\n\nGood Game IDs Sum: %+v\nSum of the powers: %+v\n\n", gameIDSum, sumOfPowers)
}
