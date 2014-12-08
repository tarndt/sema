sema: Semaphores for Go
====
### Version

0.9 Beta - API may change.

###Features

Semaphore variants provided for:
* Binary semaphores
* Counting semaphores
* Counting semaphores with timeout support
 
All written in pure Go.

###Implementations

Each varient has two implementaions:
* Channel (struct{}) based
* Condition (sync.Condition) based
 
All implementations are available at run-time and the defaults provided by the package can be toggled at runtime.

### Other Notes

The Condition-based implemenation is newer, lower-level, (aka. more bug prone), but seems to be faster. As always benchmark on your own hardware to confirm.

Basic unit-tests and benches are provided.

