// Package watcher implements FlexUI's watch mode without third-party
// dependencies. It polls file state and reports which files changed; a periodic
// rescan picks up files that were created after the build started.
package watcher

import (
	"os"
	"sort"
	"time"
)

// FileState is the polled identity of a file.
type FileState struct {
	ModTime time.Time
	Size    int64
	Missing bool
}

// Snapshot maps absolute paths to their state.
type Snapshot map[string]FileState

// Take records the state of every file.
func Take(files []string) Snapshot {
	snap := make(Snapshot, len(files))
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			snap[f] = FileState{Missing: true}
			continue
		}
		snap[f] = FileState{ModTime: info.ModTime(), Size: info.Size()}
	}
	return snap
}

// Diff returns the sorted list of paths that differ between two snapshots,
// including files that were added or removed.
func Diff(prev, next Snapshot) []string {
	var changed []string
	for path, state := range next {
		old, ok := prev[path]
		if !ok || old != state {
			changed = append(changed, path)
		}
	}
	for path := range prev {
		if _, ok := next[path]; !ok {
			changed = append(changed, path)
		}
	}
	sort.Strings(changed)
	return changed
}
