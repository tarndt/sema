// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"sync"
)

//condSema is condition-based counting semaphore implementation
type condSema struct {
	count  int
	lock   *sync.Mutex
	wakeup *sync.Cond
}

func NewCondSema() Semaphore {
	lock := new(sync.Mutex)
	return &condSema{
		count:  1,
		lock:   lock,
		wakeup: sync.NewCond(lock),
	}
}

func NewCondSemaCount(count uint) CountingSema {
	lock := new(sync.Mutex)
	return &condSema{
		count:  int(count),
		lock:   lock,
		wakeup: sync.NewCond(lock),
	}
}

func (s *condSema) P() bool {
	s.Wait(1)
	return true
}

func (s *condSema) Wait(units uint) bool {
	s.lock.Lock()
	s.count -= int(units)
	for s.count < 0 {
		s.wakeup.Wait()
	}
	s.lock.Unlock() //Not using defer for performance
	return true
}

func (s *condSema) V() bool {
	s.Signal(1)
	return true
}

func (s *condSema) Signal(units uint) bool {
	s.lock.Lock()
	wakeOthers := s.count < 0
	s.count += int(units)
	if wakeOthers {
		s.wakeup.Signal()
	}
	s.lock.Unlock() //Not using defer for performance
	return true
}
