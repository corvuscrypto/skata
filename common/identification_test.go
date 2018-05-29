package common

import (
	"testing"
	"time"
)

// TestClock is a clock used to lock into a static
// time so that a time can be generated using a common
// method
type TestClock struct{}

func (c TestClock) Now() time.Time {
	t := time.Date(2018, 3, 2, 1, 0, 0, 0, time.UTC)
	return t
}

func TestIDMethods(T *testing.T) {
	IDClock = new(TestClock)
	id := GenerateID(HeartNode)
	// test the node ID
	if nodeType := id.GetNodeType(); nodeType != HeartNode {
		T.Errorf("Expected %s got %s", HeartNode, nodeType)
	}
	// test the node Timestamp
	expected := IDClock.Now()
	if timestamp := id.GetNodeCreationTime(); timestamp != expected {
		T.Errorf("Expected %s got %s", expected, timestamp)
	}
}
