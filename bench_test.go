// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"sync"
	"testing"
	"time"
)

//Benchmark most complex semaphore variant creation

func BenchmarkCreate_condSema(b *testing.B) {
	benchCreate(NewCondSemaTimeout, b)
}

func BenchmarkCreate_chanSema(b *testing.B) {
	benchCreate(NewChanSemaTimeout, b)
}

func benchCreate(contructor func(uint, time.Duration) TimeoutCountingSema, b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var sema Semaphore
	for i := 0; i < b.N; i++ {
		sema = contructor(8192, time.Hour)
	}
	if sema == nil { //noop
	}
}

//Benchmark P

func BenchmarkP_condSema(b *testing.B) {
	benchP(NewCondSemaCount(uint(b.N)), b)
}

func BenchmarkP_chanSema(b *testing.B) {
	benchP(NewChanSemaCount(uint(b.N)), b)
}

func benchP(sema CountingSema, b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sema.P()
	}
}

//Benchmark V

func BenchmarkV_condSema(b *testing.B) {
	benchV(NewCondSemaCount(uint(b.N)), b)
}

func BenchmarkV_chanSema(b *testing.B) {
	benchV(NewChanSemaCount(uint(b.N)), b)
}

func benchV(sema CountingSema, b *testing.B) {
	sema.Wait(uint(b.N))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sema.V()
	}
}

//Benchmark wake-ups for counting semaphores

func BenchmarkWake_condSema(b *testing.B) {
	benchV(NewCondSemaCount(uint(b.N)), b)
}

func BenchmarkWake_chanSema(b *testing.B) {
	benchV(NewChanSemaCount(uint(b.N)), b)
}

func benchWake(sema CountingSema, b *testing.B) {
	sema.Wait(uint(b.N))
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			sema.P()
			wg.Done()
		}()
	}
	b.ReportAllocs()
	b.ResetTimer()
	sema.Signal(uint(b.N))
	wg.Wait()
}

//Benchmark wake-ups for counting semaphores with timeout support

func BenchmarkWakeTimeout_condSema(b *testing.B) {
	benchV(NewCondSemaTimeout(uint(b.N), time.Minute), b)
}

func BenchmarkWakeTimeout_chanSema(b *testing.B) {
	benchV(NewChanSemaTimeout(uint(b.N), time.Minute), b)
}

func benchWakeTimeout(sema TimeoutCountingSema, b *testing.B) {
	sema.Wait(uint(b.N))
	var started, finished sync.WaitGroup
	started.Add(b.N)
	finished.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			started.Done()
			sema.P()
			finished.Done()
		}()
	}
	started.Wait()
	time.Sleep(time.Millisecond * 10)
	b.ReportAllocs()
	b.ResetTimer()
	sema.Signal(uint(b.N))
	finished.Wait()
}
