// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

//chanSema is channel-based counting semaphore implementation
type chanSema chan struct{}

//NewChanSema constructs a new binary semaphore
func NewChanSema() Semaphore {
	newChanSema := make(chanSema, 1)
	newChanSema.V()
	return newChanSema
}

//NewChanSemaCount constructs a new counting semaphore
func NewChanSemaCount(count uint) CountingSema {
	newChanSema := make(chanSema, count)
	newChanSema.Signal(count)
	return newChanSema
}

func (s chanSema) Capacity() uint {
	return uint(cap(s))
}

func (s chanSema) P() bool {
	<-s
	return true
}

func (s chanSema) Acquire() bool {
	return s.P()
}

func (s chanSema) Wait(units uint) bool {
	total := units

	for ; units > 0; units-- {
		select {
		case <-s:
		default:
			if actual := total - units; actual > 0 {
				s.Signal(actual)
			}
			return false
		}
	}
	return true
}

func (s chanSema) V() bool {
	s <- struct{}{}
	return true
}

func (s chanSema) Release() bool {
	return s.V()
}

func (s chanSema) Signal(units uint) bool {
	total := units

	for ; units > 0; units-- {
		select {
		case s <- struct{}{}:
		default:
			if actual := total - units; actual > 0 {
				s.Wait(actual)
			}
			return false
		}
	}
	return true
}
