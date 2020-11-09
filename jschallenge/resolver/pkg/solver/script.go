package solver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Script is a wrapper around the cloudflare javascript challenge script
type Script interface {
	IsLoaded() bool
	Contents() string
}

type scriptImpl struct {
	loaded   bool
	lookup   lookupTable
	contents string

	challengeId string
}

// FetchScript fetches and parses the cloudflare javascript challenge script directly
func FetchScript() (Script, error) {
	log.Println("[Script] Fetching fresh script")
	res, err := http.Get("https://cloudflare.com/cdn-cgi/challenge-platform/h/g/orchestrate/jsch/v1")
	if err != nil {
		return nil, fmt.Errorf("failed requesting script: %w", err)
	}
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading script body: %w", err)
	}

	return LoadScript(string(contents))
}

// LoadScript parses the cloudflare javascript challenge script
func LoadScript(scriptContents string) (Script, error) {
	script := &scriptImpl{
		contents: scriptContents,
	}

	log.Printf("[Script] Loading script of %d bytes", len([]byte(scriptContents)))
	before := time.Now()
	err := script.load()
	if err != nil {
		return nil, err
	}
	log.Printf("[Script] Loaded script in %v", time.Now().Sub(before))

	return script, nil
}

func (s scriptImpl) IsLoaded() bool {
	return s.loaded
}

func (s scriptImpl) Contents() string {
	return s.contents
}

func (s *scriptImpl) load() (err error) {
	// Step 1: Parse lookup table
	s.lookup, err = s.parseLookupTable()
	if err != nil {
		return err
	}

	// Step 2: Find Challenge ID
	// Either look for the location where it's used or use smart regexp, I went for the latter
	challengeId, chlIdIdx, found := s.lookup.findChallengeId()
	if !found {
		return fmt.Errorf("failed finding challenge id")
	}
	s.challengeId = challengeId
	log.Printf("[Script] Found ChallengeID at lookup idx %d: %s", chlIdIdx, challengeId)

	s.loaded = true
	return nil
}
