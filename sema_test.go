// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

//Semaphore P/V test

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
		state++
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

func TestCountDown_chanSema(t *testing.T) {
	for i := 0; i < 10; i++ {
		testCountDown(NewChanSemaCount(countFrm), t)
	}
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
		for !sema.Wait(need) {
			time.Sleep(time.Millisecond / 2)
		}
		defer func() {
			for !sema.Signal(need + 1) {
				time.Sleep(time.Millisecond / 2)
			}
		}()
		results <- need
	}

	if !sema.Wait(countFrm) {
		t.Fatal("Initial claiming of resources failed")
	}
	go stateChanger(1)

	for i := uint(1); i < countFrm; i++ {
		if result := <-results; i != result {
			t.Errorf("Out of sequence count down, expected: %d, got: %d!", i, result)
			return
		}
	}
}

//Test TimeoutCountingSema timing out

const testTO = time.Millisecond * 10

func TestTimeout_chanSema(t *testing.T) {
	testTimeout(NewChanSemaTimeout(1, testTO), t)
	testTimeoutUnits(NewChanSemaTimeout(10, testTO), t)
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

func testTimeoutUnits(sema TimeoutCountingSema, t *testing.T) {
	doneCh := make(chan bool)
	sema.Wait(sema.Capacity())
	go func() {
		doneCh <- sema.Wait(sema.Capacity())
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

func TestContext_chanSema(t *testing.T) {
	for i := 0; i < 10; i++ {
		testContext(NewChanSemaTimeout(1, testTO), t)
	}
	for i := 0; i < 10; i++ {
		testContextUnits(NewChanSemaTimeout(10, testTO), t)
	}
}

func testContext(sema TimeoutCountingSema, t *testing.T) {
	testCtx, testCancel := context.WithTimeout(context.Background(), time.Second)
	defer testCancel()

	if !sema.AcquireCtx(testCtx) {
		t.Fatalf("Intial aquire failed")
	}
	var wg sync.WaitGroup
	wg.Add(2)

	acquireErr := fmt.Errorf("AcquireCtx never unblocked")
	releaseErr := fmt.Errorf("ReleaseCtx never unblocked")

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(testCtx, time.Microsecond)
		defer cancel()

		if sema.AcquireCtx(ctx) {
			acquireErr = fmt.Errorf("Inital acquire should have failed")
			return
		}

		time.Sleep(time.Millisecond * 2)
		ctx, cancel = context.WithTimeout(testCtx, time.Microsecond)
		defer cancel()
		if sema.ReleaseCtx(ctx) {
			releaseErr = fmt.Errorf("Followup release should have failed")
			return
		}

		ctx, cancel = context.WithTimeout(testCtx, time.Microsecond)
		defer cancel()
		if sema.AcquireCtx(ctx) {
			acquireErr = nil
		} else {
			acquireErr = fmt.Errorf("Final acquire failed")
		}
	}()

	go func() {
		defer wg.Done()

		time.Sleep(time.Millisecond)
		if sema.ReleaseCtx(testCtx) {
			releaseErr = nil
		}
	}()

	wg.Wait()
	switch {
	case acquireErr != nil:
		t.Errorf("AcquireCtx failure: %s", acquireErr)
	case releaseErr != nil:
		t.Errorf("ReleaseCtx failure: %s", releaseErr)
	}
}

func testContextUnits(sema TimeoutCountingSema, t *testing.T) {
	testCtx, testCancel := context.WithTimeout(context.Background(), time.Second)
	defer testCancel()

	success, units := sema.WaitCtx(testCtx, sema.Capacity())
	switch {
	case !success:
		t.Fatalf("Intial wait failed")
	case units != sema.Capacity():
		t.Fatalf("Final intial wait got %d units rather than %d", units, sema.Capacity())
	}

	var wg sync.WaitGroup
	wg.Add(2)

	signalErr := fmt.Errorf("SignalCtx never unblocked")
	waitErr := fmt.Errorf("WaitCtx never unblocked")

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(testCtx, time.Microsecond)
		defer cancel()

		if success, _ := sema.WaitCtx(ctx, sema.Capacity()); success {
			waitErr = fmt.Errorf("Inital wait should have failed")
			return
		}

		time.Sleep(time.Millisecond * 2)
		ctx, cancel = context.WithTimeout(testCtx, time.Microsecond)
		defer cancel()
		if success, units := sema.SignalCtx(ctx, sema.Capacity()); success {
			signalErr = fmt.Errorf("Followup signal should have failed but %d of %d units were released", units, sema.Capacity())
			return
		} else if units > 0 {
			signalErr = fmt.Errorf("Followup signal failed, but should have release no units and %d od %d units were released", units, sema.Capacity())
			return
		}

		ctx, cancel = context.WithTimeout(testCtx, time.Microsecond)
		defer cancel()
		success, units := sema.WaitCtx(ctx, sema.Capacity())
		switch {
		case !success:
			waitErr = fmt.Errorf("Final wait failed")
		case units != sema.Capacity():
			waitErr = fmt.Errorf("Final wait got %d units rather than %d", units, sema.Capacity())
		default:
			waitErr = nil
		}
	}()

	go func() {
		defer wg.Done()

		time.Sleep(time.Millisecond)
		if success, _ := sema.SignalCtx(testCtx, sema.Capacity()); success {
			signalErr = nil
		}
	}()

	wg.Wait()
	switch {
	case waitErr != nil:
		t.Errorf("WaitCtx failure: %s", waitErr)
	case signalErr != nil:
		t.Errorf("SignalCtx failure: %s", signalErr)
	}
}
