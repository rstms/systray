//go:build !windows

package icon

import (
	_ "embed"
)

//go:embed icon.png
var Data []byte
