package solver

import "fmt"

type Options struct {
	CvId    string
	CType   string
	CNounce string
	CRay    string
	CHash   string
	CFPWv   string
}

type Solver interface {
	Solve(options Options) error
}

func NewSolver() Solver {
	return &solverImpl{}
}

func NewSolverFromScript(script Script) Solver {
	return &solverImpl{
		script: script,
	}
}

type solverImpl struct {
	script Script
}

func (s *solverImpl) Solve(options Options) (err error) {
	// TODO: Check if script is expired
	if s.script == nil {
		s.script, err = FetchScript()
		if err != nil {
			return fmt.Errorf("failed solving challenge: %w", err)
		}
	}

	return nil
}
