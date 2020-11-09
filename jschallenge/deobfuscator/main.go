package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var inputFile = flag.String("input", "jschal.js", "The js challenge javascript file")
var outputFile = flag.String("output", "jschal.deobfuscated.js", "The deobfuscated javascript file location")
var lookupTableExpr = regexp.MustCompile(`a = '(.+?)'`)
var shuffleRoundsExpr = regexp.MustCompile(`}\(a, (\d+)\),`)
var lookupTableUsageExpr = regexp.MustCompile(`\[?b\('(0x[0-9a-f]+)'\)\]?`)
var expressionFuncExpr = regexp.MustCompile(`(?m)(\w\[b\('0x[0-9a-f]+'\)\]) = function\(\w, \w\) {[\r\n][\s\p{Zs}]+return \w ([\S]+) \w[\r\n][\s\p{Zs}]+}`)

func parseLookupTable(contents string) ([]string, error) {
	log.Printf("Parsing lookup table")

	matches := lookupTableExpr.FindStringSubmatch(contents)
	if len(matches) != 2 {
		return nil, fmt.Errorf("Failed to find lookup table (variable a)")
	}

	lookupTable := strings.Split(matches[1], ",")
	log.Printf("Parsed lookup table, contains %d values", len(lookupTable))

	// Shuffle the table with the b function
	shuffleRoundsMatches := shuffleRoundsExpr.FindStringSubmatch(contents)
	if len(shuffleRoundsMatches) != 2 {
		return nil, fmt.Errorf("Failed to parse shuffle rounds (function after a)")
	}

	shuffleRounds, err := strconv.Atoi(shuffleRoundsMatches[1])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse shuffle rounds (function after a): %w", err)
	}
	log.Printf("Parsed shuffle rounds: %d", shuffleRounds)

	for i := shuffleRounds; i > 0; i-- {
		first := lookupTable[0]
		lookupTable = lookupTable[1:]
		lookupTable = append(lookupTable, first)
	}

	log.Printf("Applied shuffle rounds")

	return lookupTable, nil
}

func applyLookupTable(contents string, lookupTable []string) string {
	log.Println("Applying lookup table")
	replacementStats := 0
	newContents := lookupTableUsageExpr.ReplaceAllStringFunc(contents, func(match string) string {
		parts := lookupTableUsageExpr.FindStringSubmatch(match)
		lookupTableIdx, err := strconv.ParseInt(parts[1], 0, 0)
		if err != nil {
			log.Printf("WARN: Failed converting %s to number: %v", parts[1], err)
			return match
		}

		lookupTableValue := lookupTable[lookupTableIdx]

		var result string
		runes := []rune(match)
		// If encapsulated with both a [ and ], replace with .{lookupTableValue}
		if runes[0] == '[' && runes[len(runes)-1] == ']' {
			result = fmt.Sprintf(".%s", lookupTableValue)
		} else {
			result = fmt.Sprintf("'%s'", lookupTableValue)

			// Item can be part of an array, I'm lazy so deal with this fix here ;)
			if runes[0] == '[' {
				result = "[" + result
			}
			if runes[0] == ']' {
				result = result + "]"
			}
		}

		log.Printf("Replacing %s, with lookup value %s", match, result)
		replacementStats++

		return result
	})

	log.Printf("Applied lookup table, replaced %d references", replacementStats)

	return newContents
}

type expressionStats struct {
	funcs []string
	count int
}

func replaceExpressionFunctions(contents string) string {

	expressionFuncStatsMap := make(map[string]*expressionStats)

	expressionFuncMatches := expressionFuncExpr.FindAllStringSubmatch(strings.TrimSpace(contents), -1)
	for _, expressionFuncParts := range expressionFuncMatches {
		funcName, funcOperator := expressionFuncParts[1], expressionFuncParts[2]

		if stats := expressionFuncStatsMap[funcOperator]; stats != nil {
			stats.count++
			stats.funcs = append(stats.funcs, funcName)
		} else {
			expressionFuncStatsMap[funcOperator] = &expressionStats{
				funcs: []string{funcName},
				count: 1,
			}
		}

		// log.Printf("Found expression function: %s, operator: %s", funcName, funcOperator)
	}

	for operatorName, stats := range expressionFuncStatsMap {
		log.Printf("Found %d %s operator functions: %v", stats.count, operatorName, stats.funcs)
	}

	// TODO: Actually replace them, CF used a nasty trick to copy the object containing the expression functions to another object so it's harder to find

	return contents
}

func main() {
	log.Println("Cloudflare JSChallenge Deobfuscator 1.0.0 by Paul Hobbel")
	flag.Parse()

	log.Printf("Opening %s", *inputFile)
	rawContents, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Opened file, read %d bytes", len(rawContents))

	contents := string(rawContents)
	lookupTable, err := parseLookupTable(contents)
	if err != nil {
		log.Fatal(err)
	}

	contents = replaceExpressionFunctions(contents)

	contents = applyLookupTable(contents, lookupTable)

	log.Printf("Done, writing result to %s", *outputFile)
	err = ioutil.WriteFile(*outputFile, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
