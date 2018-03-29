// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fsnotify implements filesystem notification.
package fsnotify

import (
	"fmt"
	"log"
	"sync"
)

const (
	FSN_CREATE = 1
	FSN_MODIFY = 2
	FSN_DELETE = 4
	FSN_RENAME = 8

	FSN_ALL = FSN_MODIFY | FSN_DELETE | FSN_RENAME | FSN_CREATE
)

var (
	FsWatcher *TFsWatcher
)

type (
	TFsWatcher struct {
		*Watcher
		NewsMu sync.Mutex // Map access
		News   map[string]*FileEvent
	}
)

func init() { /*
		log.Println("FsWatcher:Open")
		//创建Watcher
		watcher, err := NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		//创建文件系统Watcher
		FsWatcher = &TFsWatcher{
			Watcher: watcher,
			News:    make(map[string]*FileEvent),
		}
		log.Println("go:", err)
		// Process events
		go func() {
			for {
				select {
				case ev := <-FsWatcher.Event:
					if ev == nil {
						log.Println("Event:", ev)
					} else {
						log.Println("Event:", ev)
						//FsWatcher.News[ev.Name] = ev
					}
					//log.Println("Event:", ev)
				case err := <-FsWatcher.Error:
					if err == nil {

					} else {
						log.Println("Error:", err)
					}
				default:

				}
			}
		}()
		err = FsWatcher.Watch("/modules/")
		log.Println("goout:", err)
		//err = FsWatcher.Watch("/")
		log.Println("goout1:", err)

		log.Println("Watch:", err)

		FsWatcher.Close()
	*/
}

func (self *TFsWatcher) Watch(path string) error {
	return self.watch(path)
}

func (self *TFsWatcher) IsNew(path string) bool {
	log.Println("IsNew:", path)
	self.NewsMu.Lock()
	_, ok := self.News[path]
	if ok {
		self.NewsMu.Unlock()
		return false
	}
	return true
}

func (self *TFsWatcher) Del(path string) {
	log.Println("IsNew:", path)
	self.NewsMu.Lock()
	delete(self.News, path)
	self.NewsMu.Unlock()
}

// Purge events from interal chan to external chan if passes filter
func (w *Watcher) purgeEvents() {
	for ev := range w.internalEvent {
		sendEvent := false
		w.fsnmut.Lock()
		fsnFlags := w.fsnFlags[ev.Name]
		w.fsnmut.Unlock()

		if (fsnFlags&FSN_CREATE == FSN_CREATE) && ev.IsCreate() {
			sendEvent = true
		}

		if (fsnFlags&FSN_MODIFY == FSN_MODIFY) && ev.IsModify() {
			sendEvent = true
		}

		if (fsnFlags&FSN_DELETE == FSN_DELETE) && ev.IsDelete() {
			sendEvent = true
		}

		if (fsnFlags&FSN_RENAME == FSN_RENAME) && ev.IsRename() {
			//w.RemoveWatch(ev.Name)
			sendEvent = true
		}

		if sendEvent {
			w.Event <- ev
		}
	}

	close(w.Event)
}

// Watch a given file path
func (w *Watcher) Watch(path string) error {
	w.fsnmut.Lock()
	w.fsnFlags[path] = FSN_ALL
	w.fsnmut.Unlock()
	return w.watch(path)
}

// Watch a given file path for a particular set of notifications (FSN_MODIFY etc.)
func (w *Watcher) WatchFlags(path string, flags uint32) error {
	w.fsnmut.Lock()
	w.fsnFlags[path] = flags
	w.fsnmut.Unlock()
	return w.watch(path)
}

// Remove a watch on a file
func (w *Watcher) RemoveWatch(path string) error {
	w.fsnmut.Lock()
	delete(w.fsnFlags, path)
	w.fsnmut.Unlock()
	return w.removeWatch(path)
}

// String formats the event e in the form
// "filename: DELETE|MODIFY|..."
func (e *FileEvent) String() string {
	var events string = ""

	if e.IsCreate() {
		events += "|" + "CREATE"
	}

	if e.IsDelete() {
		events += "|" + "DELETE"
	}

	if e.IsModify() {
		events += "|" + "MODIFY"
	}

	if e.IsRename() {
		events += "|" + "RENAME"
	}

	if len(events) > 0 {
		events = events[1:]
	}

	return fmt.Sprintf("%q: %s", e.Name, events)
}
