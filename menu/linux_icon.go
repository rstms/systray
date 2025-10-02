//go:build linux

package menu

import (
	_ "embed"
	"fmt"
)

//go:embed icon.png
var DefaultIconData []byte
