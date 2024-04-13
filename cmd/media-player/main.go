package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go"
	"github.com/saltosystems/winrt-go/windows/foundation"
	"github.com/saltosystems/winrt-go/windows/media/core"
	"github.com/saltosystems/winrt-go/windows/media/playback"
)

func must[T any](d T, err error) T {
	must0(err)
	return d
}

func must0(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	must0(ole.CoInitialize(0))
	p := must(playback.NewMediaPlayer())

	playbackSession := must(p.GetPlaybackSession())
	fmt.Println(playbackSession)

	_, path, _, _ := runtime.Caller(0)
	uri := must(foundation.UriCreateUri("file:///" + filepath.Dir(path) + "/a.flac"))

	s := must(core.MediaSourceCreateFromUri(uri))

	must0(p.SetSource((*playback.IMediaPlaybackSource)(unsafe.Pointer(s))))
	//must0(p.SetVolume(0.8))

	must0(p.SetAudioCategory(playback.MediaPlayerAudioCategoryMedia))

	eventReceivedGuid := winrt.ParameterizedInstanceGUID(
		foundation.GUIDTypedEventHandler,
		playback.SignatureMediaPlaybackSession,
		"cinterface(IInspectable)",
	)
	t1 := must(playbackSession.AddPlaybackStateChanged(
		foundation.NewTypedEventHandler(
			ole.NewGUID(eventReceivedGuid), func(_ *foundation.TypedEventHandler, sender, args unsafe.Pointer) {
				fmt.Println("playback state changed", sender, args)
				fmt.Println(playbackSession.GetPlaybackState())
			})))
	fmt.Println(t1)

	eventReceivedGuid2 := winrt.ParameterizedInstanceGUID(
		foundation.GUIDTypedEventHandler,
		playback.SignatureMediaPlayer,
		"cinterface(IInspectable)",
	)
	t2 := must(p.AddCurrentStateChanged(
		foundation.NewTypedEventHandler(
			ole.NewGUID(eventReceivedGuid2), func(_ *foundation.TypedEventHandler, sender, args unsafe.Pointer) {
				fmt.Println("current state changed", sender, args)
				fmt.Println(((*playback.MediaPlaybackSession)(sender)).GetPlaybackState())
			})))
	fmt.Println(t2)

	t3 := must(p.AddMediaEnded(
		foundation.NewTypedEventHandler(
			ole.NewGUID(eventReceivedGuid2), func(_ *foundation.TypedEventHandler, sender, args unsafe.Pointer) {
				fmt.Println("media ended", sender, args)
				fmt.Println(((*playback.MediaPlaybackSession)(sender)).GetPlaybackState())
			})))
	fmt.Println(t3)

	go func() {
		must0(p.Play())
		fmt.Println(must(p.GetPlaybackSession()))
	}()

	c := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			<-time.After(time.Second)
			st, _ := playbackSession.GetPlaybackState()
			if i >= 20 || st == playback.MediaPlaybackStatePaused {
				c <- struct{}{}
			}
			fmt.Println(st)
		}
	}()

	//go func() {
	//	<-time.After(time.Second * 20)
	//	//playbackSession.GetPlaybackState()
	//	err = p.Pause()
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	//time.Sleep(time.Hour)
	<-c
}
