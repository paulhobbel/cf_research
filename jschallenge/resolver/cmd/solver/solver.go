package main

import (
	"github.com/paulhobbel/cloudflare-challenge-solver/pkg/solver"
	"log"
)

func main() {
	log.Println("[Main] Cloudflare JSChallenge Solver 1.0.0 by Paul Hobbel")
	log.Println("[Main] Caution: This software is intended for educational purposes only!")
	s := solver.NewSolver()

	err := s.Solve(solver.Options{
		CvId:    "1",
		CType:   "non-interactive",
		CNounce: "37224",
		CRay:    "5eeb63977ad09cab",
		CHash:   "04edc127e40c570",
		CFPWv:   "g",
	})
	if err != nil {
		log.Fatalf("[Main] Failed solving challenge: %v", err)
	}

	log.Fatalf("[Main] Failed solving challenge: failed decoding cipher")
	//log.Printf("[Main] Loaded script: %v", s.)
}
