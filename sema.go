// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"context"
	"time"
)

//Semaphore is a binary semaphore supports operations to gain and release
// control of exclusive execution used to protect one or more critical sections.
//The original terminology P (aka. signal/acquire) and V (aka. wait/release)
// are used given by they Semaphore inventor Edsger Dijkstra in 1965.
type Semaphore interface {
	P() bool
	Acquire() bool

	V() bool
	Release() bool
}

//A CountingSema is a semaphores which allows an arbitrary resource count.
// Wait/Signal are used for the mult-P / multi-V operations.
type CountingSema interface {
	Semaphore
	//These are HARD to use correctly and may be removed prior to v1
	Wait(units uint) bool
	Signal(units uint) bool
}

//A TimeoutCountingSema further extends a semaphore to provide the ability to
// timeout while waiting to gain access to a critical section.
type TimeoutCountingSema interface {
	CountingSema

	//Semaphore P (aka. Acquire)
	PTO(timeout time.Duration) bool
	AcquireTO(timeout time.Duration) bool
	AcquireCtx(ctx context.Context) bool

	//Semaphore V (aka. Release)
	VTO(timeout time.Duration) bool
	ReleaseTO(timeout time.Duration) bool
	ReleaseCtx(ctx context.Context) bool

	//Counting Semaphore multi-unit P (aka. Wait)
	// returns the if all units were acquired followed by the number of units
	// actually acquired.
	WaitTO(units uint, timeout time.Duration) (bool, uint)
	WaitCtx(ctx context.Context, units uint) (bool, uint)

	//Counting Semaphore multi-unit V (aka. Signal)
	// returns the if all units were released followed by the number of units
	// actually released.
	SignalTO(units uint, timeout time.Duration) (bool, uint)
	SignalCtx(ctx context.Context, units uint) (bool, uint)
}

//Exported constructors, this is configurable to allow multiple implementation
// stratigies. In the past we provided a condition based implementation that is
// no longer faster enough to justify the complexity and had some likely bugs.
var (
	NewSemaphore    = NewChanSema
	NewCountingSema = NewChanSemaCount
	NewTimeoutSema  = NewChanSemaTimeout
)
