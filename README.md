sema: Semaphores for Go
====
### Author & Version

Author: [Tylor Arndt]

0.95 RC - API may change. (Bug reports and PRs are welcome)

### Features

[Semaphore](https://en.wikipedia.org/wiki/Semaphore_(programming)) variants provided are:

* Binary semaphores
* Counting semaphores
* Counting semaphores with timeout support
* [Context](https://pkg.go.dev/context) support was recently added
 
All written in pure Go.

### Implementations

Previously there was a [sync.Cond](http://golang.org/pkg/sync/#Cond) based implemenation that was removed with
Go runtime performance improvements rendered it overly complex for a small performance gain
over the channel-based implemenation.

### Getting Started

In [sema.go](https://github.com/tarndt/sema/blob/master/sema.go) you will find the three default constructors and related core interfaces.
```go
	func NewSemaphore() Semaphore {...}
	func NewCountingSema(count uint) CountingSema {...}
	func NewTimeoutSema(count uint, defaultTimeout time.Duration) TimeoutCountingSema {...}
```
The `Semaphore` interface is extended from being binary to counting by `CountingSema` which in turn is enhanced with time-out support in its `TimeoutCountingSema` variant.
