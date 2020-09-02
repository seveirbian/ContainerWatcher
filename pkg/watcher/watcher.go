package watcher

import (
	"cw/pkg/file"
	"cw/pkg/watcher/docker"
	"cw/pkg/watcher/overlay2"
)

// Watcher a watcher for watching container's actions
type Watcher interface {
	Watch() chan file.File
}

// NewWatcher create a new watcher based on docker graphdriver
func NewWatcher(digest string) Watcher {
	info, err := docker.Client.Info(docker.Ctx)
	if err != nil {
		return nil
	}

	switch info.Driver {
	case "overlay2":
		w := overlay2.NewWatcher(digest)
		return w
	default:
	}

	return nil
}
