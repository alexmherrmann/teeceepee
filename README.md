
# A tcp test server (like httpbin)

This is being built to test a microcontroller networking library so don't expect much.

The server just waits for connections, reads a "command" and then responds.
Commands must be terminated with a newline.


### Commands

* `drip bytes times delay`: drip feed a few bytes a number of times with a delay
    * it goes `delay bytes delay bytes delay bytes`
