/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package config

// Config holds the application configuration flags
type Config struct {
	Level        string
	Exclude      []string
	TruncateRaw  bool
	RemovePrefix string
}
