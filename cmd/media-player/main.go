package main

import (
	"path/filepath"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"

	"github.com/saltosystems/winrt-go/windows/foundation"
	"github.com/saltosystems/winrt-go/windows/media/core"
	"github.com/saltosystems/winrt-go/windows/media/playback"
)

func main() {
	_ = ole.CoInitialize(0)
	p, err := playback.NewMediaPlayer()
	if err != nil {
		panic(err)
	}

	_, path, _, _ := runtime.Caller(0)
	uri, err := foundation.UriCreateUri("file:///" + filepath.Dir(path) + "/a.flac")
	if err != nil {
		panic(err)
	}
	s, err := core.MediaSourceCreateFromUri(uri)
	if err != nil {
		panic(err)
	}
	err = p.SetSource((*playback.IMediaPlaybackSource)(unsafe.Pointer(s)))
	if err != nil {
		panic(err)
	}
	go func() {
		err = p.Play()
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Hour)
}
