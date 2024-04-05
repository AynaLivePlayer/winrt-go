package main

import (
	"github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go/windows/foundation"
	"github.com/saltosystems/winrt-go/windows/media/core"
	"github.com/saltosystems/winrt-go/windows/media/playback"
	"time"
	"unsafe"
)

func main() {
	ole.CoInitialize(0)
	p, err := playback.NewMediaPlayer()
	if err != nil {
		panic(err)
	}
	uri, err := foundation.UriCreateUri("http://m801.music.126.net/20240406002338/cb697f571b59a1f6f635b53dcbfffd48/jdyyaac/obj/w5rDlsOJwrLDjj7CmsOj/20161325709/cf66/e916/8bf0/4e906e7aec4cbc4a504136e0473b4706.m4a")
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
