package solver

import "regexp"

var lookupTableExpr = regexp.MustCompile(`a='(.+?)'`)
var shuffleRoundsExpr = regexp.MustCompile(`}\(a,(\d+)\),`)
var lookupTableUsageExpr = regexp.MustCompile(`\[?b\('(0x[0-9a-f]+)'\)\]?`)

var (
	// / generation time : revision/id : sha265 hash /
	challengeIdExpr = regexp.MustCompile(`/\d+\.?\d*:\d+:[a-f0-9]{64}/`)
)
