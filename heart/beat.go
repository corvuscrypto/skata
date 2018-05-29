package heart

import (
	"time"
)

// Beat is the type that represents a single actable moment
// in time. The way this works is that a structure is sent with
// a reference day and a beat number which represents a duration of your choosing.
type Beat struct {
	Day        time.Time
	BeatNumber int
	// Sure I could back-calculate the duration, but why?
	Duration time.Duration
}

// ToTime converts the beat into a full datetime
func (b Beat) ToTime() time.Time {
	return b.Day.Add(time.Duration(b.BeatNumber) * b.Duration)
}

// IsEqual determines if the Beats are the same
func (b Beat) IsEqual(other Beat) bool {
	return b.Day.Equal(other.Day) && b.BeatNumber == other.BeatNumber
}

// BeatFromTime takes a Time structure and returns a Beat-representation
func BeatFromTime(t time.Time, duration time.Duration) Beat {
	var beat Beat
	beat.Day = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	beat.BeatNumber = int((t.UnixNano() - beat.Day.UnixNano()) / int64(duration))
	beat.Duration = duration
	return beat
}

// BeatGenerator is any type that can expose methods to
// be used as a generator for daily beats. This is so
// anyone can override the way beats are generated or define their own.
type BeatGenerator interface {
	StartBeating()
	GetBeatChannel() <-chan Beat
	StopBeating()
}

// DefaultBeatGenerator is the default generator
type DefaultBeatGenerator struct {
	lastBeat     Beat
	ticker       *time.Ticker
	beatDuration time.Duration
	beatChannel  chan Beat
	stopChannel  chan bool
	started      bool
}

// NewBeatGenerator is the instantiation method. USE THIS INSTEAD OF new()
func NewBeatGenerator(beatDuration time.Duration) (gen *DefaultBeatGenerator) {
	gen = new(DefaultBeatGenerator)
	gen.beatChannel = make(chan Beat)
	gen.stopChannel = make(chan bool)
	gen.beatDuration = beatDuration
	gen.ticker = time.NewTicker(beatDuration)
	return
}

func (g *DefaultBeatGenerator) beatSubroutine() {
	for {
		select {
		case currentTime := <-g.ticker.C:
			currentTime = currentTime.UTC()
			currentBeat := BeatFromTime(currentTime, g.beatDuration)
			if currentBeat.IsEqual(g.lastBeat) {
				// Just in case the Ticker is early, go to sleep for a bit
				nextBeatTime := g.lastBeat.ToTime().Add(g.beatDuration)
				time.Sleep(nextBeatTime.Sub(currentTime))
				currentBeat = BeatFromTime(nextBeatTime, g.beatDuration)
			}
			// send the beat
			g.beatChannel <- currentBeat
			// set last beat as the current beat... yeah.
			g.lastBeat = currentBeat
		case <-g.stopChannel:
			return
		}
	}
}

// StartBeating starts the heartbeat
func (g *DefaultBeatGenerator) StartBeating() {
	if g.started {
		return
	}

	g.started = true
	g.ticker.Stop()
	g.ticker = time.NewTicker(g.beatDuration)
	go g.beatSubroutine()
}

// StopBeating stops the heartbeat
func (g *DefaultBeatGenerator) StopBeating() {
	if !g.started {
		return
	}

	g.started = false
	g.stopChannel <- true
	g.ticker.Stop()
	close(g.beatChannel)
	g.beatChannel = make(chan Beat)
}

// GetBeatChannel returns the internal channel so beats can be retrieved
func (g *DefaultBeatGenerator) GetBeatChannel() <-chan Beat {
	return g.beatChannel
}
