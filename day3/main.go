package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type part struct {
	PartNo   int
	PosStart int
	PosEnd   int
	Valid    bool
}

type gear struct {
	PosStart int
	PosEnd   int
	Valid    bool
}

type ps struct {
	Line  string
	Parts []part
	Gears []gear
}

func main() {
	schematic := make(map[int]ps)

	// Read file
	file, err := os.Open("./input.txt")
	if err != nil {
		slog.Error("cannot read file", "error", err)
		os.Exit(1)
	}

	defer file.Close()

	rxNo, err := regexp.Compile(`[\d]{1,}`)
	if err != nil {
		slog.Error("cannot compile rxNo regex", "error", err)
		os.Exit(1)
	}

	rxSymbol, err := regexp.Compile(`[^.\d\n]`)
	if err != nil {
		slog.Error("cannot compile rxSymbol regex", "error", err)
		os.Exit(1)
	}

	rxGear, err := regexp.Compile(`[\*]`)
	if err != nil {
		slog.Error("cannot compile rxGear regex", "error", err)
		os.Exit(1)
	}

	// Scan the file lines
	scanner := bufio.NewScanner(file)
	lineCounter := 0
	gearRatioSum := 0

	for scanner.Scan() {
		line := scanner.Text()

		schematic[lineCounter] = ps{
			Line: line,
		}

		// Find parts
		matches := rxNo.FindAllStringIndex(line, -1)
		for _, m := range matches {
			partNo, _ := strconv.Atoi(line[m[0]:m[1]])

			sch := schematic[lineCounter]
			sch.Parts = append(sch.Parts, part{
				PartNo:   partNo,
				PosStart: m[0],
				PosEnd:   m[1],
			})

			schematic[lineCounter] = sch
		}

		// Find gears
		gearMatches := rxGear.FindAllStringIndex(line, -1)
		for _, m := range gearMatches {
			sch := schematic[lineCounter]
			sch.Gears = append(sch.Gears, gear{
				PosStart: m[0],
				PosEnd:   m[1],
			})

			schematic[lineCounter] = sch
		}

		lineCounter += 1
	}

	for _, v := range schematic {
		fmt.Printf("\n%+v", v.Line)
	}

	fmt.Printf("\n\n\n")

	for lineNr, v := range schematic {
		top, hasTop := schematic[lineNr-1]
		if hasTop {
			fmt.Printf("\ntop:     %+v", top.Line)
		} else {
			fmt.Printf("\ntop:")
		}

		fmt.Printf("\ncurrent: %+v", v.Line)

		bottom, hasBottom := schematic[lineNr+1]
		if hasBottom {
			fmt.Printf("\nbottom:  %+v", bottom.Line)
		} else {
			fmt.Printf("\nbottom:")
		}

		for ip, p := range v.Parts {
			var adjacentSymbols strings.Builder

			fmt.Printf("\n  - %+v", p.PartNo)

			// Check left
			sIdx := p.PosStart - 1
			eIdx := p.PosStart

			if sIdx >= 0 && eIdx <= len(v.Line)-1 {
				adjacentSymbols.WriteString(v.Line[sIdx:eIdx])
			}

			// Check right
			sIdx = p.PosEnd
			eIdx = p.PosEnd + 1

			if sIdx >= 0 && eIdx <= len(v.Line)-1 {
				adjacentSymbols.WriteString(v.Line[sIdx:eIdx])
			}

			// Check top
			top, hasTop := schematic[lineNr-1]
			if hasTop == true {
				sIdx := p.PosStart - 1
				eIdx := p.PosEnd + 1

				if sIdx < 0 {
					sIdx = 0
				}
				if eIdx > len(top.Line)-1 {
					eIdx = len(top.Line) - 1
				}

				adjacentSymbols.WriteString(top.Line[sIdx:eIdx])
			}

			// Check bottom
			bottom, hasBottom := schematic[lineNr+1]
			if hasBottom == true {
				sIdx := p.PosStart - 1
				eIdx := p.PosEnd + 1

				if sIdx < 0 {
					sIdx = 0
				}
				if eIdx > len(bottom.Line)-1 {
					eIdx = len(bottom.Line) - 1
				}

				adjacentSymbols.WriteString(bottom.Line[sIdx:eIdx])
			}

			schematic[lineNr].Parts[ip].Valid = rxSymbol.MatchString(adjacentSymbols.String())
			fmt.Printf("\n    %+v", adjacentSymbols.String())
			fmt.Printf("\n    - valid: %+v", schematic[lineNr].Parts[ip].Valid)
		}

		fmt.Printf("\n  - Gears")

		for _, g := range v.Gears {
			var overlappingParts []int

			sIdx := g.PosStart
			eIdx := g.PosEnd

			if sIdx < 0 {
				sIdx = 0
			}

			if eIdx > len(v.Line)-1 {
				eIdx = len(v.Line) - 1
			}

			// Check current line
			overlappingPart := getAdjacentPart(v.Parts, sIdx, eIdx)
			if len(overlappingPart) > 0 {
				overlappingParts = append(overlappingParts, overlappingPart...)
			}

			// Check top
			top, hasTop := schematic[lineNr-1]
			if hasTop == true {
				overlappingPart := getAdjacentPart(top.Parts, sIdx, eIdx)
				if len(overlappingPart) > 0 {
					overlappingParts = append(overlappingParts, overlappingPart...)
				}
			}

			// Check Bottom
			bottom, hasBottom := schematic[lineNr+1]
			if hasBottom == true {
				overlappingPart := getAdjacentPart(bottom.Parts, sIdx, eIdx)
				if len(overlappingPart) > 0 {
					overlappingParts = append(overlappingParts, overlappingPart...)
				}
			}

			fmt.Printf("\n    - overlap parts: %+v", overlappingParts)

			// Correct gear
			if len(overlappingParts) == 2 {
				gearRatioSum += overlappingParts[0] * overlappingParts[1]
			}

		}

		fmt.Printf("\n")
	}

	var validParts []int
	for _, v := range schematic {
		for _, p := range v.Parts {
			if p.Valid == false {
				continue
			}

			validParts = append(validParts, p.PartNo)
		}
	}

	sum := 0
	for _, partNo := range validParts {
		sum += partNo
	}

	fmt.Printf("\n\n\n Valid parts sum: %+v", sum)
	fmt.Printf("\n\n Gear ratio sum: %+v\n\n\n", gearRatioSum)

	if err := scanner.Err(); err != nil {
		slog.Error("scanning file failed", "error", err)
		os.Exit(1)
	}
}

func getAdjacentPart(data []part, sIdx, eIdx int) (ret []int) {
	for _, p := range data {
		overlap := math.Max(0, math.Min(float64(eIdx), float64(p.PosEnd))-math.Max(float64(sIdx), float64(p.PosStart))+1)
		if overlap > 0 {
			ret = append(ret, p.PartNo)
		}
	}

	return
}
