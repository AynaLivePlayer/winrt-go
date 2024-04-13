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
	"github.com/saltosystems/winrt-go/windows/media"
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
	must0(ole.RoInitialize(1))

	p := must(playback.NewMediaPlayer())
	defer p.Release()

	playbackSession := must(p.GetPlaybackSession())
	defer playbackSession.Release()
	fmt.Println(playbackSession)

	_, path, _, _ := runtime.Caller(0)
	uri := must(foundation.UriCreateUri("file:///" + filepath.Dir(path) + "/a.flac"))
	defer uri.Release()

	s := must(core.MediaSourceCreateFromUri(uri))
	defer s.Release()

	must0(p.SetSource((*playback.IMediaPlaybackSource)(unsafe.Pointer(s))))
	//must0(p.SetVolume(0.8))
	must0(p.SetAudioCategory(playback.MediaPlayerAudioCategoryMedia))

	eventReceivedGuid := winrt.ParameterizedInstanceGUID(
		foundation.GUIDTypedEventHandler,
		playback.SignatureMediaPlaybackSession,
		"cinterface(IInspectable)",
	)
	h1 := foundation.NewTypedEventHandler(
		ole.NewGUID(eventReceivedGuid), func(_ *foundation.TypedEventHandler, sender, _ unsafe.Pointer) {
			session := (*playback.MediaPlaybackSession)(sender)
			fmt.Println("playback state changed", must(session.GetPlaybackState()))
		})
	defer h1.Release()
	t1 := must(playbackSession.AddPlaybackStateChanged(h1))
	fmt.Println(t1)

	eventReceivedGuid2 := winrt.ParameterizedInstanceGUID(
		foundation.GUIDTypedEventHandler,
		playback.SignatureMediaPlayer,
		"cinterface(IInspectable)",
	)
	h2 := foundation.NewTypedEventHandler(
		ole.NewGUID(eventReceivedGuid2), func(_ *foundation.TypedEventHandler, sender, _ unsafe.Pointer) {
			player := (*playback.MediaPlayer)(sender)
			fmt.Println("current state changed", must(player.GetCurrentState()))
		})
	defer h2.Release()
	t2 := must(p.AddCurrentStateChanged(h2))
	fmt.Println(t2)

	smtc := must(p.GetSystemMediaTransportControls())
	defer smtc.Release()
	must0(smtc.SetIsEnabled(true))
	must0(smtc.SetIsPauseEnabled(true))
	must0(smtc.SetIsPlayEnabled(true))
	eventReceivedGuid3 := winrt.ParameterizedInstanceGUID(
		foundation.GUIDTypedEventHandler,
		media.SignatureSystemMediaTransportControls,
		media.SignatureSystemMediaTransportControlsButtonPressedEventArgs,
	)
	h4 := foundation.NewTypedEventHandler(
		ole.NewGUID(eventReceivedGuid3), func(_ *foundation.TypedEventHandler, sender, args unsafe.Pointer) {
			ctrl := (*media.SystemMediaTransportControls)(sender)
			arg := (*media.SystemMediaTransportControlsButtonPressedEventArgs)(args)
			fmt.Println("button pressed", ctrl, arg)
		})
	defer h4.Release()
	t4 := must(smtc.AddButtonPressed(h4))
	fmt.Println(t4)

	go p.Play()

	c := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			<-time.After(time.Second)
			st, _ := playbackSession.GetPlaybackState()
			//if i >= 5 {
			//	must0(p.SetVolume(1.0))
			//}
			if st == playback.MediaPlaybackStatePaused {
				c <- struct{}{}
			}
			fmt.Println(st)
		}
	}()

	<-c
}
