//go:build !windows

package menu

import (
	_ "embed"
)

//go:embed icon.png
var DefaultIconData []byte
