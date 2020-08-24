package gobat

import (
	"errors"
	"sort"
	"time"
)

type CommonBatConfig struct {
	// The time to start the batch start time
	StartTime time.Time
	// Interval to go to check batch start time
	CheckInterval time.Duration
	// Interval from start time to start of next batch execution
	ActionInterval time.Duration
}

type OneWayBatConfig struct {
	Common *CommonBatConfig
}

type ParallelBatConfig struct {
	Common     *CommonBatConfig
	Dependency *Dependency
}

type Dependency struct {
	bats []*Batch
}

type Batch struct {
	priority int
	function []func()
}

func SetCommonBatConfig(st time.Time, ci, ai time.Duration) *CommonBatConfig {
	return &CommonBatConfig{
		StartTime:      st,
		CheckInterval:  ci,
		ActionInterval: ai,
	}
}

func SetOneWayBatConfig(common *CommonBatConfig) *OneWayBatConfig {
	return &OneWayBatConfig{
		Common: common,
	}
}

func SetParallelBatConfig(common *CommonBatConfig, dependency *Dependency) *ParallelBatConfig {
	return &ParallelBatConfig{
		Common:     common,
		Dependency: dependency,
	}
}

func SetPriority(p int, fc ...func()) *Batch {
	return &Batch{
		priority: p,
		function: fc,
	}
}

func GenerateDependency(ps ...*Batch) (*Dependency, error) {
	dependency := Dependency{
		bats: ps,
	}
	sort.Sort(dependency)

	// check duplicate
	var prePriority = -1
	for _, p := range dependency.bats {
		if prePriority == p.priority {
			return nil, errors.New("duplicate priority")
		}
		prePriority = p.priority
	}
	return &dependency, nil
}

func (c *CommonBatConfig) NextSchedule() {
	c.StartTime = c.StartTime.Add(c.ActionInterval)
}

// Parallel batch processing
func (b *ParallelBatConfig) ParallelBatRun() error {
	if b.Dependency != nil {
		cs := b.Common
		for {
			if time.Now().Equal(cs.StartTime) || time.Now().After(cs.StartTime) {
				for _, bat := range b.Dependency.bats {
					var running = make(chan int, len(bat.function))
					for _, fc := range bat.function {
						go func(fc func()) {
							fc()
							running <- bat.priority
						}(fc)
					}
					for {
						if len(running) == len(bat.function) {
							<-running
							break
						}
					}
				}
				break
			}
			time.Sleep(cs.CheckInterval)
		}
		return nil
	}
	return errors.New("dependency is not set")
}

// One-way batch processing
func (b *OneWayBatConfig) OneWayBatRun(f ...func()) error {
	var running = make(chan int, 1)
	cs := b.Common
	for {
		if time.Now().Equal(cs.StartTime) || time.Now().After(cs.StartTime) {
			for _, fc := range f {
				go func(fc func()) {
					fc()
					running <- 1
				}(fc)
				<-running
			}
			break
		}
		time.Sleep(cs.CheckInterval)
	}
	return nil
}

func (t Dependency) Len() int {
	return len(t.bats)
}

func (t Dependency) Swap(i, j int) {
	t.bats[i], t.bats[j] = t.bats[j], t.bats[i]
}

func (t Dependency) Less(i, j int) bool {
	return t.bats[i].priority < t.bats[j].priority
}
