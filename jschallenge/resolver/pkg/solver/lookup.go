package solver

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type lookupTable []string

func (s *scriptImpl) parseLookupTable() (lookupTable, error) {
	log.Printf("[Script] Parsing lookup table")

	matches := lookupTableExpr.FindStringSubmatch(s.contents)
	if len(matches) != 2 {
		return nil, fmt.Errorf("failed to find lookup table (variable a)")
	}

	table := lookupTable(strings.Split(matches[1], ","))
	log.Printf("[Script] Parsed lookup table, contains %d values", len(table))

	// Shuffle the table with the b function
	shuffleRoundsMatches := shuffleRoundsExpr.FindStringSubmatch(s.contents)
	if len(shuffleRoundsMatches) != 2 {
		return nil, fmt.Errorf("failed to parse shuffle rounds (function after a)")
	}

	shuffleRounds, err := strconv.Atoi(shuffleRoundsMatches[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse shuffle rounds (function after a): %w", err)
	}
	log.Printf("[Script] Parsed shuffle rounds: %d", shuffleRounds)

	for i := shuffleRounds; i > 0; i-- {
		first := table[0]
		table = table[1:]
		table = append(table, first)
	}

	log.Printf("[Script] Applied shuffle rounds")

	return table, nil
}

// findChallengeId tries to find the challenge id based on a regexp
func (t lookupTable) findChallengeId() (string, int, bool) {
	for idx, value := range t {
		if challengeIdExpr.MatchString(value) {
			return value, idx, true
		}
	}

	return "", -1, false
}
