//go:build windows

package menu

import (
	_ "embed"
)

//go:embed icon.ico
var DefaultIconData []byte
