// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"time"
)

//condSema is condition-based counting semaphore implementation with timeout
// support
type condSemaTO struct {
	*condSema
	countMax                  int
	activeTimeout, defTimeout time.Duration
	timer                     *time.Timer
}

func NewCondSemaTimeout(count uint, defaultTimeout time.Duration) TimeoutCountingSema {
	s := &condSemaTO{
		condSema:   NewCondSemaCount(count).(*condSema),
		countMax:   int(count),
		defTimeout: defaultTimeout,
		timer:      time.NewTimer(time.Hour),
	}
	s.timer.Stop()
	go s.timeoutWorker()
	return s
}

func (s *condSemaTO) timeoutWorker() {
	for {
		<-s.timer.C
		s.wakeup.Broadcast()
	}
}

func (s *condSemaTO) P() bool {
	return s.PTO(s.defTimeout)
}

func (s *condSemaTO) PTO(timeout time.Duration) bool {
	return s.WaitTO(1, timeout)
}

func (s *condSemaTO) Wait(units uint) bool {
	return s.WaitTO(units, s.defTimeout)
}

func (s *condSemaTO) WaitTO(units uint, timeout time.Duration) bool {
	endOfLife := time.Now().Add(timeout)
	s.lock.Lock()
	switch {
	case s.activeTimeout > timeout:
		s.wakeup.Broadcast() //We don't know how much of the old timeout was used
		fallthrough
	case s.activeTimeout == 0:
		s.activeTimeout = timeout
		s.timer.Reset(timeout)
	}
	s.count -= int(units)
	for s.count < 0 {
		if time.Now().After(endOfLife) {
			s.lock.Unlock()
			return false
		}
		s.wakeup.Wait()
	}
	if s.count == s.countMax {
		s.timer.Stop()
		s.activeTimeout = 0
	}
	s.lock.Unlock() //Not using defer for performance
	return true
}

func (s *condSemaTO) V() bool {
	return s.VTO(s.defTimeout)
}

func (s *condSemaTO) VTO(timeout time.Duration) bool {
	return s.SignalTO(1, timeout)
}

func (s *condSemaTO) Signal(units uint) bool {
	return s.SignalTO(units, s.defTimeout)
}

func (s *condSemaTO) SignalTO(units uint, timeout time.Duration) bool {
	s.lock.Lock()
	wakeOthers := s.count < 0
	s.count += int(units)
	if wakeOthers {
		s.wakeup.Signal()
	}
	s.lock.Unlock() //Not using defer for performance
	return true
}
