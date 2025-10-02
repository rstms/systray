//go:build windows

package menu

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestMenu(t *testing.T) {

	title := "Test Title"
	tooltip := "The Tooltray Test Program"

	menu := NewMenu(title, tooltip, []byte{})
	err := menu.Start()
	require.Nil(t, err)
	timer := time.NewTimer(3 + time.Second)
	for done := false; !done; {
		select {
		case item := <-menu.Clicked:
			log.Printf("received Clicked: %v\n", item)
			require.IsType(t, &MenuItem{}, item)
			done = true
		case exited := <-menu.Exited:
			log.Printf("received Exited: %v\n", exited)
			require.IsType(t, struct{}{}, exited)
			done = true

		case <-timer.C:
			log.Println("timeout")
			done = true
		}
	}
}
