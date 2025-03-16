package timing

import "time"

type ExponentialTimer interface {
	Succeeded()
	Failed()
	Stop()
}

type expoTimer struct {
	timer *time.Timer

	currDur time.Duration
	minDur  time.Duration
	maxDur  time.Duration
}

func New(min, max time.Duration) ExponentialTimer {
	return &expoTimer{
		minDur:  min,
		maxDur:  max,
		currDur: min,
	}
}

// Succeeded resets the timer interval to the minimum and waits until timer expires.
// Should be called to indicate that the process utilizing the timer has succeeded.
func (t *expoTimer) Succeeded() {
	t.currDur = t.minDur
	t.timer = time.NewTimer(t.currDur)
	<-t.timer.C
}

// Failed doubles the timer interval (caps the interval to the maximum) and waits until timer expires.
// Should be called to indicate that the process utilizing the timer has failed and must back off.
func (t *expoTimer) Failed() {
	t.currDur *= 2
	if t.currDur > t.maxDur {
		t.currDur = t.maxDur
	}

	t.timer = time.NewTimer(t.currDur)
	<-t.timer.C
}

// Stop stops the timer
func (t *expoTimer) Stop() {
	t.timer.Stop()
}
