// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package sema

//chanSema is channel-based counting semaphore implementation
type chanSema chan struct{}

func NewChanSema() Semaphore {
	newChanSema := make(chanSema, 1)
	newChanSema.V()
	return newChanSema
}

func NewChanSemaCount(count uint) CountingSema {
	newChanSema := make(chanSema, count)
	newChanSema.Signal(count)
	return newChanSema
}

func (s chanSema) P() bool {
	<-s
	return true
}

func (s chanSema) Wait(units uint) bool {
	for ; units > 0; units-- {
		<-s
	}
	return true
}

func (s chanSema) V() bool {
	s <- struct{}{}
	return true
}

func (s chanSema) Signal(units uint) bool {
	for ; units > 0; units-- {
		s <- struct{}{}
	}
	return true
}
