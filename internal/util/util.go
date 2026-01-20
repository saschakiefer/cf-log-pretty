/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package util

import (
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var (
	terminalWidth     = 80 // fallback default
	lastCheck         time.Time
	checkInterval     = 2 * time.Second
	terminalWidthLock sync.RWMutex
)

// GetTerminalWidth retrieves the current width of the terminal in characters, updating it periodically if needed.
func GetTerminalWidth() int {
	terminalWidthLock.RLock()
	if time.Since(lastCheck) < checkInterval {
		defer terminalWidthLock.RUnlock()
		return terminalWidth
	}
	terminalWidthLock.RUnlock()

	// Update width
	terminalWidthLock.Lock()
	defer terminalWidthLock.Unlock()
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && width > 0 {
		terminalWidth = width
	}
	lastCheck = time.Now()
	return terminalWidth
}
