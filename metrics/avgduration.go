package metrics

import (
	"fmt"
	"time"
)

type AverageDuration struct {
	RollingAvg   int64
	Measurements int64
}

func (p *AverageDuration) Add(duration time.Duration) {
	p.Measurements++

	prevAvg := p.RollingAvg
	p.RollingAvg = prevAvg + (duration.Nanoseconds()-prevAvg)/p.Measurements
}

func (p *AverageDuration) String() string {
	d, _ := time.ParseDuration(fmt.Sprintf("%dns", p.RollingAvg))
	return d.Round(time.Millisecond).String()
}

func (p *AverageDuration) Raw() int64 {
	d, _ := time.ParseDuration(fmt.Sprintf("%dns", p.RollingAvg))
	return d.Round(time.Millisecond).Milliseconds()
}

func (p *AverageDuration) Reset() {
	p.RollingAvg = 0
	p.Measurements = 0
}
