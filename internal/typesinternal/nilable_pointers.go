// Copyright authors of this Go fork
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typesinternal

import "strings"

// ParseNilablePointersFromMod returns the nilable_pointers mode from go.mod content.
func ParseNilablePointersFromMod(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "nilable_pointers ") {
			f := strings.Fields(line)
			if len(f) == 2 && (f[1] == "enable" || f[1] == "disable" || f[1] == "warnings") {
				return f[1]
			}
		}
	}
	return ""
}
