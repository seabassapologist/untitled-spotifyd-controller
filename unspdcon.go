package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)

type Playing struct {
	Title       string
	Album       string
	ArtUrl      string
	Status      string
	Artist      []string
	AlbumArtist []string
}

func GetPlaying(bus *dbus.Conn) Playing {

	var info map[string]dbus.Variant
	sp := bus.Object(
		"org.mpris.MediaPlayer2.spotifyd",
		dbus.ObjectPath("/org/mpris/MediaPlayer2"),
	).Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.mpris.MediaPlayer2.Player").Store(&info)

	if sp != nil {
		fmt.Println(sp)
		fmt.Println("Check if spotifyd is running!")
		os.Exit(1)
	}

	var meta map[string]dbus.Variant
	info["Metadata"].Store(&meta)

	var ti, al, au, st string = "", "", "", ""
	var ar, aa []string = []string{}, []string{}

	info["PlaybackStatus"].Store(&st)

	if len(meta) == 0 {
		if st == "Playing" || st == "Paused" {
			// There's a track playing, but no Metadata available
			ti = "?"
			al = "?"
			ar = append(ar, "?")
			aa = append(aa, "?")
		} else {
			ti = "Not Playing"
			st = "Stopped"
		}

	} else {
		meta["xesam:title"].Store(&ti)
		meta["xesam:album"].Store(&al)
		meta["mpris:artUrl"].Store(&au)
		meta["xesam:albumArtist"].Store(&aa)
		meta["xesam:artist"].Store(&ar)
	}

	return Playing{ti, al, au, st, ar, aa}

}

func Output(p Playing) {
	fmt.Printf(
		"Title:     %v\nAlbum:     %v\nArtist(s): %v\nPlayback:  %v\n",
		p.Title, p.Album, strings.Join(p.AlbumArtist[:], ", "), p.Status,
	)
}

func OutputWaybar(p Playing) {
	var text, tooltip string
	if p.Status == "Stopped" {
		text = "..."
		tooltip = "Not Playing"
	} else {

		text = fmt.Sprintf("%v • %v", p.Title, p.AlbumArtist[0])
		tooltip = fmt.Sprintf(
			"Title:     %v\\nAlbum:     %v\\nArtist(s): %v\\nPlayback:  %v",
			p.Title, p.Album, strings.Join(p.Artist[:], ", "), p.Status,
		)

		if p.Status == "Playing" {
			text += "   "
		} else if p.Status == "Paused" {
			text += "   "
		}
	}

	fmt.Printf("{\"text\": \"%v\", \"tooltip\": \"%v\", \"class\": \"$class\"}\n", text, tooltip)
}

func main() {

	// Get dbus connection to Session Bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Println("Failed to connect to Session Bus ", err)
		os.Exit(1)
	}

	pl := GetPlaying(conn)

	OutputWaybar(pl)
}
