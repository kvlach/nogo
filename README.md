# nogo

Check if a hostname is blacklisted.

Wraps https://github.com/StevenBlack/hosts into an easy to use library.

## Example

```go
package main

import (
	"github.com/kvlach/nogo"
)

func main() {
	// Not turning any of the Fakenews, Gambling, Porn, Social flags on will
	// mean that it will only download the default adware + malware list
	ng, err := nogo.Init().Fakenews().Gambling().Porn().Social().Download()
	if err != nil {
		// ...
	}

	if !ng.Safe("google.com") {
		// ...
	}

	// to update the list just run download again
	ng, err = ng.Download()
	if err != nil {
		// ...
	}
}
```
