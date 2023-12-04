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

type card struct {
	id      int
	winning []int
	picks   []int
	matches []int
	points  int
	rec     bool
	p2      int
}

func main() {
	// Read file
	file, err := os.Open("./input.txt")
	if err != nil {
		slog.Error("cannot read file", "error", err)
		os.Exit(1)
	}

	defer file.Close()

	rxCardID := regexp.MustCompile(`([\d]{1,3}):`)
	rxNumbers := regexp.MustCompile(`([\d]{1,2})`)

	var cards []card
	var sumPoints int

	// Scan the file lines
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		match := rxCardID.FindStringSubmatchIndex(line)

		id := line[int(match[2]):int(match[3])]
		cardID, _ := strconv.Atoi(string(id))

		// split the winning from picks
		rest := line[int(match[3])+1:]
		spl := strings.Split(rest, "|")

		winning := rxNumbers.FindAllString(spl[0], -1)
		picks := rxNumbers.FindAllString(spl[1], -1)
		var ws, ps, matches []int

		// convert winning to int
		for _, w := range winning {
			r, _ := strconv.Atoi(w)
			ws = append(ws, r)
		}

		// convert picks to int
		for _, p := range picks {
			r, _ := strconv.Atoi(p)
			ps = append(ps, r)

			for _, w := range ws {
				if w == r {
					matches = append(matches, w)
				}
			}
		}

		// Calc points
		points := 0
		if len(matches) > 0 {
			points = 1
			for i := 1; i < len(matches); i++ {
				points *= 2
			}
		}

		cards = append(cards, card{
			id:      cardID,
			winning: ws,
			picks:   ps,
			matches: matches,
			points:  points,
		})

		sumPoints += points
	}

	fmt.Printf("\n\n\nSUM Points: %+v\n\n\n", sumPoints)

	// Part 2
	func (c *card) calcP2Points(cards []card) int {
		// exit recursion
		if c.rec == true {
			return c.p2
		}
	
		ret := 0
		for i := c.id; i < len(c.matches)+c.id; i++ {
			ret += cards[i].calcP2Points(cards)
		}
	
		c.rec = true
		c.p2 = ret + len(c.matches)
	
		return c.p2
	}

	
	p2Points := 0
	for _, c := range cards {
		p2Points += c.calcP2Points(cards)
	}

	fmt.Printf("\nSUM P2: %+v", p2Points+len(cards))

	fmt.Printf("\n\n\n")
	if err := scanner.Err(); err != nil {
		slog.Error("scanning file failed", "error", err)
		os.Exit(1)
	}
}
