package heart

import (
	"testing"
	"time"
)

func TestDefaultBeatGenerator(T *testing.T) {
	today := time.Now().UTC().Truncate(time.Hour * 24)
	generator := NewBeatGenerator(time.Second)
	testChan := generator.GetBeatChannel()
	generator.StartBeating()
	beat := <-testChan
	// Test the day
	if !beat.Day.Equal(today) {
		T.Errorf("expected %s, got %s", today, beat.Day)
	}
	// Test that the beats are in fact one apart
	beat2 := <-testChan
	if difference := (beat2.BeatNumber - beat.BeatNumber); difference != 1 {
		T.Errorf("beat difference is off. Difference: %d", difference)
	}
	generator.StopBeating()

}
