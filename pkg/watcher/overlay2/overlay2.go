package overlay2

import (
	"ContainerWatcher/pkg/file"
	"ContainerWatcher/pkg/log"
	"ContainerWatcher/pkg/watcher/docker"
	"os"
	"path/filepath"
	"strings"

	fsevents "github.com/tywkeene/go-fsevents"
)

// Watcher Watcher on overlay2
type Watcher struct {
	paths []string
	Files chan file.File
}

// NewWatcher create a overlay2 watcher
func NewWatcher(digest string) *Watcher {
	info, _, err := docker.Client.ImageInspectWithRaw(docker.Ctx, digest)
	if err != nil {
		return &Watcher{paths: []string{}}
	}

	upperLayersPath := info.GraphDriver.Data["UpperDir"]
	lowerLayersPath := info.GraphDriver.Data["LowerDir"]

	upperLayersPathArray := strings.Split(upperLayersPath, ":")
	lowerLayersPathArray := strings.Split(lowerLayersPath, ":")

	layersPath := []string{}
	layersPath = append(layersPath, upperLayersPathArray...)
	layersPath = append(layersPath, lowerLayersPathArray...)

	// remove "" string
	deduplicatedLayersPath := []string{}
	for i := len(layersPath) - 1; i >= 0; i-- {
		if layersPath[i] != "" {
			deduplicatedLayersPath = append(deduplicatedLayersPath, layersPath[i])
		}
	}

	log.Logger.Info("Get path from image: ", deduplicatedLayersPath)

	return &Watcher{paths: deduplicatedLayersPath, Files: make(chan file.File, 200)}
}

// Watch start watch
func (o *Watcher) Watch() chan file.File {
	var mask uint32 = fsevents.Accessed | fsevents.Open | fsevents.CloseRead
	w, err := fsevents.NewWatcher()
	if err != nil {
		log.Logger.Fatal(err)
	}

	for _, p := range o.paths {
		fAndD := traverPath(p)
		for _, item := range fAndD {
			d, err := w.AddDescriptor(item, mask)
			if err != nil {
				log.Logger.Warn(err)
			}
			if err := d.Start(); err != nil {
				log.Logger.Warn(err)
			}
		}
	}

	go o.handleEvents(w)

	return o.Files
}

func (o *Watcher) handleEvents(w *fsevents.Watcher) {
	go w.Watch()

	for {
		select {
		case event := <-w.Events:
			fi, err := os.Lstat(event.Path)
			if err != nil {
				log.Logger.Warn(err)
			}
			if fi.Mode().IsRegular() {
				o.Files <- file.File{Name: event.Path, Size: fi.Size(), Times: 0}
			}
		case err := <-w.Errors:
			log.Logger.Warn(err)
		default:
		}
	}
}

func traverPath(root string) []string {
	fileAndDir := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Mode().IsRegular() {
			fileAndDir = append(fileAndDir, path)
		}
		return nil
	})
	if err != nil {
		log.Logger.Warn(err)
	}

	return fileAndDir
}
