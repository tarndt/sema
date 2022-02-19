// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

import (
	"context"
	"time"
)

//condSema is channel-based counting semaphore implementation with timeout
// support
type chanSemaTO struct {
	chanSema
	defTimeout time.Duration
}

//NewChanSemaTimeout constructs a new counting semaphore that will construct
// a new counting semaphore with all P/V (acquire/release) opterations timing out
// after the providing default timeout unless overridden at call time (see PTO, VTO).
//Important: If you would like P/V (acquire/release) to be non-blocking you should
// provide a defaultTimeout of zero or use NewNonBlockChanSema.
func NewChanSemaTimeout(count uint, defaultTimeout time.Duration) TimeoutCountingSema {
	return &chanSemaTO{
		chanSema:   NewChanSemaCount(count).(chanSema),
		defTimeout: defaultTimeout,
	}
}

//NewNonBlockChanSema is like NewChanSemaTimeout with a defaultTimeout of zero.
func NewNonBlockChanSema(count uint) TimeoutCountingSema {
	return NewChanSemaTimeout(count, 0)
}

func (s *chanSemaTO) P() bool {
	return s.PTO(s.defTimeout)
}

func (s *chanSemaTO) Acquire() bool {
	return s.P()
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

func (s *chanSemaTO) AcquireTO(timeout time.Duration) bool {
	return s.PTO(timeout)
}

func (s *chanSemaTO) AcquireCtx(ctx context.Context) bool {
	select {
	case <-s.chanSema:
		return true
	default:
		select {
		case <-s.chanSema:
			return true
		case <-ctx.Done():
			return false
		}
	}
}

func (s *chanSemaTO) Wait(units uint) (success bool) {
	success, _ = s.WaitTO(units, s.defTimeout)
	return
}

func (s *chanSemaTO) WaitTO(units uint, timeout time.Duration) (bool, uint) {
	var opTimedOut <-chan time.Time //We avoid allocating this if possible
	total := units

	for ; units > 0; units-- {
		select {
		case <-s.chanSema:
			continue
		default:
			if timeout < 1 {
				return units == 0, total - units
			}
			if opTimedOut == nil {
				opTimedOut = time.After(timeout)
			}
			select {
			case <-s.chanSema:
				continue
			case <-opTimedOut:
				return units == 0, total - units
			}
		}
	}
	return true, 0
}

func (s *chanSemaTO) WaitCtx(ctx context.Context, units uint) (bool, uint) {
	total, doneCh := units, ctx.Done()

	for ; units > 0; units-- {
		select {
		case <-s.chanSema:
			continue
		default:
			select {
			case <-s.chanSema:
				continue
			case <-doneCh:
				return units == 0, total - units
			}
		}
	}
	return true, 0
}

func (s *chanSemaTO) V() bool {
	return s.VTO(s.defTimeout)
}

func (s *chanSemaTO) Release() bool {
	return s.V()
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

func (s *chanSemaTO) ReleaseTO(timeout time.Duration) bool {
	return s.VTO(timeout)
}

func (s *chanSemaTO) ReleaseCtx(ctx context.Context) bool {
	select {
	case s.chanSema <- struct{}{}:
		return true
	default:
		select {
		case s.chanSema <- struct{}{}:
			return true
		case <-ctx.Done():
			return false
		}
	}
}

func (s *chanSemaTO) Signal(units uint) (success bool) {
	success, _ = s.SignalTO(units, s.defTimeout)
	return
}

func (s *chanSemaTO) SignalTO(units uint, timeout time.Duration) (bool, uint) {
	var opTimedOut <-chan time.Time //We avoid allocating this if possible
	total := units

	for ; units > 0; units-- {
		select {
		case s.chanSema <- struct{}{}:
			continue
		default:
			if timeout < 1 {
				return units == 0, total - units
			}
			if opTimedOut == nil {
				opTimedOut = time.After(timeout)
			}
			select {
			case s.chanSema <- struct{}{}:
				continue
			case <-opTimedOut:
				return units == 0, total - units
			}
		}
	}
	return true, 0
}

func (s *chanSemaTO) SignalCtx(ctx context.Context, units uint) (bool, uint) {
	total, doneCh := units, ctx.Done()

	for ; units > 0; units-- {
		select {
		case s.chanSema <- struct{}{}:
			continue
		default:
			select {
			case s.chanSema <- struct{}{}:
				continue
			case <-doneCh:
				return units == 0, total - units
			}
		}
	}
	return true, 0
}
