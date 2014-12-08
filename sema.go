// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"time"
)

//A binary semaphore supports an operation to gain and release control of
// exclusive execution used to protect one or more critical sections.
//The orginal terminology P (aka. signal/acquire) and V (aka. wait/release)
// are used given by they Semaphore inventor Edsger Dijkstra in 1965.
type Semaphore interface {
	P() bool
	V() bool
}

//A CountingSema is a semaphores which allows an arbitrary resource count.
// Wait/Signal are used for the mult-P / multi-V operations.
type CountingSema interface {
	Semaphore
	Wait(units uint) bool
	Signal(units uint) bool
}

//A TimeoutCountingSema further extends a semaphore to provide the ability to
// timeout while waiting to gain access to a critical section.
type TimeoutCountingSema interface {
	CountingSema
	PTO(timeout time.Duration) bool
	VTO(timeout time.Duration) bool
	WaitTO(units uint, timeout time.Duration) bool
	SignalTO(units uint, timeout time.Duration) bool
}

/* Provided default implementations based on benchmarking:
BenchmarkP_condSema			  100000000  18.6 ns/op  0 B/op  0 allocs/op
BenchmarkP_chanSema			  100000000  22.3 ns/op  0 B/op	 0 allocs/op
BenchmarkV_condSema			  100000000  18.5 ns/op  0 B/op	 0 allocs/op
BenchmarkV_chanSema			  100000000  22.8 ns/op  0 B/op	 0 allocs/op
BenchmarkWake_condSema		  100000000  18.4 ns/op  0 B/op	 0 allocs/op
BenchmarkWake_chanSema		  100000000  23.1 ns/op  0 B/op	 0 allocs/op
BenchmarkWakeTimeout_condSema 100000000  25.9 ns/op  0 B/op	 0 allocs/op
BenchmarkWakeTimeout_chanSema 50000000   36.1 ns/op  0 B/op	 0 allocs/op */
var (
	NewSemaphore    = NewCondSema
	NewCountingSema = NewCondSemaCount
	NewTimeoutSema  = NewCondSemaTimeout
)

//UseChannelBasedImpl sets the default exported Semaphore variant
// implementations to be channel based
func UseChannelBasedImpl() {
	NewSemaphore = NewChanSema
	NewCountingSema = NewChanSemaCount
	NewTimeoutSema = NewChanSemaTimeout
}

//UseChannelBasedImpl sets the default exported Semaphore variant
// implementations to be channel based
func UseConditionBasedImpl() {
	NewSemaphore = NewCondSema
	NewCountingSema = NewCondSemaCount
	NewTimeoutSema = NewCondSemaTimeout
}
