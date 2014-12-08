// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"sync"
	"testing"
	"time"
)

//Semaphore P/V test

func TestPV_condSema(t *testing.T) {
	testP(NewCondSemaCount(1), t)
}

func TestPV_chanSema(t *testing.T) {
	testP(NewChanSemaCount(1), t)
}

func testP(sema CountingSema, t *testing.T) {
	x, y := -1, -1
	var started, finished sync.WaitGroup
	started.Add(1)
	finished.Add(1)
	sema.P()
	go func() {
		started.Done()
		sema.P()
		y = x + 1
		sema.V()
		finished.Done()
	}()
	started.Wait()
	x = 5
	sema.V()
	finished.Wait()
	if !(x == 5 && y == 6) {
		t.Errorf("sema.P() failed to protect x,y; x = 5, y = 6 was expected, x = %d, y = %d was found.", x, y)
	}
}

//Another more demanding Semaphore P/V test

const countFrm = 100

func TestPVRecur_condSema(t *testing.T) {
	testPVRecur(NewCondSema(), t)
}

func TestPVRecur_chanSema(t *testing.T) {
	testPVRecur(NewChanSema(), t)
}

func testPVRecur(sema Semaphore, t *testing.T) {
	results := make(chan int, countFrm)
	state := 0
	var stateChanger func()
	stateChanger = func() {
		sema.P()
		defer sema.V()
		if state > countFrm {
			return
		}
		go stateChanger()
		results <- state
		state += 1
	}
	go stateChanger()
	for i := 0; i < countFrm; i++ {
		if result := <-results; i != result {
			t.Errorf("Out of sequence count down, expected: %d, got: %d!", i, result)
			return
		}
	}
}

//Test CountingSema counting down with Wait/Signal

func TestCountDown_condSema(t *testing.T) {
	testCountDown(NewCondSemaCount(countFrm), t)
}

func TestCountDown_chanSema(t *testing.T) {
	testCountDown(NewChanSemaCount(countFrm), t)
}

func testCountDown(sema CountingSema, t *testing.T) {
	results := make(chan uint, countFrm)
	var stateChanger func(need uint)
	stateChanger = func(need uint) {
		switch {
		case need == countFrm:
			sema.V() //Let the last/(this) goroutine finish
		case need > countFrm:
			return
		}
		go stateChanger(need + 1)
		sema.Wait(need)
		defer sema.Signal(need + 1)
		results <- need
	}
	sema.Wait(countFrm)
	go stateChanger(1)
	sema.V()
	for i := uint(1); i < countFrm; i++ {
		if result := <-results; i != result {
			t.Errorf("Out of sequence count down, expected: %d, got: %d!", i, result)
			return
		}
	}
}

//Test TimeoutCountingSema timing out

const testTO = time.Millisecond * 10

func TestTimeout_condSema(t *testing.T) {
	testTimeout(NewCondSemaTimeout(1, testTO), t)
}

func TestTimeout_chanSema(t *testing.T) {
	testTimeout(NewChanSemaTimeout(1, testTO), t)
}

func testTimeout(sema TimeoutCountingSema, t *testing.T) {
	testP(sema, t)
	doneCh := make(chan bool)
	sema.P()
	go func() {
		doneCh <- sema.P()
	}()
	select {
	case success := <-doneCh:
		if success {
			t.Errorf("P() returned: true, when is should have timed-out and returned: false.")
		}
	case <-time.After(testTO + time.Millisecond*10):
		t.Errorf("Test timed-out: Timeout sema did not timeout and return as expected")
	}

}
