// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"time"
)

//condSema is channel-based counting semaphore implementation with timeout
// support
type chanSemaTO struct {
	chanSema
	defTimeout time.Duration
}

func NewChanSemaTimeout(count uint, defaultTimeout time.Duration) TimeoutCountingSema {
	return &chanSemaTO{
		chanSema:   NewChanSemaCount(count).(chanSema),
		defTimeout: defaultTimeout,
	}
}

func (s *chanSemaTO) P() bool {
	return s.PTO(s.defTimeout)
}

func (s *chanSemaTO) PTO(timeout time.Duration) bool {
	select {
	case <-s.chanSema:
		return true
	default:
		if timeout < 1 {
			return false
		}
		select {
		case <-s.chanSema:
			return true
		case <-time.After(timeout):
			return false
		}
	}
}

func (s *chanSemaTO) Wait(units uint) bool {
	return s.WaitTO(units, s.defTimeout)
}

func (s *chanSemaTO) WaitTO(units uint, timeout time.Duration) bool {
	var opTimedOut <-chan time.Time //We avoid allocating this if possible
	for ; units > 0; units-- {
		select {
		case <-s.chanSema:
			continue
		default:
			if timeout < 1 {
				return false
			}
			if opTimedOut == nil {
				opTimedOut = time.After(timeout)
			}
			select {
			case <-s.chanSema:
				continue
			case <-opTimedOut:
				return false
			}
		}
	}
	return true
}

func (s *chanSemaTO) V() bool {
	return s.VTO(s.defTimeout)
}

func (s *chanSemaTO) VTO(timeout time.Duration) bool {
	select {
	case s.chanSema <- struct{}{}:
		return true
	default:
		if timeout < 1 {
			return false
		}
		select {
		case s.chanSema <- struct{}{}:
			return true
		case <-time.After(timeout):
			return false
		}
	}
}

func (s *chanSemaTO) Signal(units uint) bool {
	return s.SignalTO(units, s.defTimeout)
}

func (s *chanSemaTO) SignalTO(units uint, timeout time.Duration) bool {
	var opTimedOut <-chan time.Time //We avoid allocating this if possible
	for ; units > 0; units-- {
		select {
		case s.chanSema <- struct{}{}:
			continue
		default:
			if timeout < 1 {
				return false
			}
			if opTimedOut == nil {
				opTimedOut = time.After(timeout)
			}
			select {
			case s.chanSema <- struct{}{}:
				continue
			case <-opTimedOut:
				return false
			}
		}
	}
	return true
}
