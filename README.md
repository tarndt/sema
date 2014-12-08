sema: Semaphores for Go
====
### Author & Version

Author: [Tylor Arndt]

0.9 Beta - API may change. (Bug reports and PRs are welcome)

###Features

Semaphore variants provided for:
* Binary semaphores
* Counting semaphores
* Counting semaphores with timeout support
 
All written in pure Go.

###Implementations

Each variant has two implementaions:
* Channel (`struct{}`) based
* Condition ([sync.Cond]) based
 
All implementations are available at run-time and the defaults provided by the package can be toggled at runtime.

### Other Notes

The Condition-based implemenation is newer, lower-level, (aka. more bug prone), but seems to be faster. As always benchmark on your own hardware to confirm.

Basic unit-tests and benches are provided.

### Getting Started

In [sema.go] you will find the three default constructors and core interfaces.
```go
	func NewSemaphore() Semaphore {...}
	func NewCountingSema(count uint) CountingSema {...}
	func NewTimeoutSema(count uint, defaultTimeout time.Duration) TimeoutCountingSema {...}
```
The `Semaphore` interface is extended from being binary to counting by `CountingSema` which in turn is enhanced with time-out support in its `TimeoutCountingSema` variant.

[Tylor Arndt]:https://plus.google.com/u/0/+TylorArndt/posts
[sync.Cond]:http://golang.org/pkg/sync/#Cond
[sema.go]:https://github.com/tarndt/sema/blob/master/sema.go]




