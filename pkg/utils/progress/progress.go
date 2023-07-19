package progress

import "math"

type Progress struct {
	// the maximum progress value
	// optional, default 100
	Max int

	// the Current progress value
	// optional, default 0
	Current int

	// defines the progress Steps, for example if this is set to 2, then
	// a progress that has max set to 100 will take 50 Steps
	//
	// used if the function progress.add() is used
	//
	// optional, default 1
	Steps int
}

// sets the current progress directly without keeping in mind the steps
func (p *Progress) Set(v int) {
	p.Current = v
}

// adds v steps to the current progress, if the parameter is 0 / not defined, then the
// it will Add one step to the progress
func (p *Progress) Add(v int) {
	if p.Steps == 0 {
		p.Steps = 1
	}

	if v == 0 {
		p.Current += p.Steps
	} else {
		p.Current += v * p.Steps
	}
}

// returns the progress in percent
func (p *Progress) Get() int {
	if p.Max == 0 {
		p.Max = 100
	}
	return int(math.Round(float64(p.Current) / float64(p.Max) * 100))
}
