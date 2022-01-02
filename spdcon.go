package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)

type Track struct {
	Title       string
	Album       string
	ArtUrl      string
	Status      string
	Artist      []string
	AlbumArtist []string
}

type Waybar struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

// Wrapper function to execute org.mpris.MediaPlayer2.Player Methods
func SendCommand(bus *dbus.Conn, mt string) {
	bus.Object(
		"org.mpris.MediaPlayer2.spotifyd",
		dbus.ObjectPath("/org/mpris/MediaPlayer2"),
	).Call("org.mpris.MediaPlayer2.Player."+mt, 0)
}

func GetPlaying(bus *dbus.Conn) Track {

	var info map[string]dbus.Variant
	sp := bus.Object(
		"org.mpris.MediaPlayer2.spotifyd",
		dbus.ObjectPath("/org/mpris/MediaPlayer2"),
	).Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.mpris.MediaPlayer2.Player").Store(&info)

	var ti, al, au, st string = "", "", "", ""
	var ar, aa []string = []string{}, []string{}

	if sp != nil {
		ti = "Not Playing"
		st = "Stopped"
	} else {
		var meta map[string]dbus.Variant
		info["Metadata"].Store(&meta)

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

	}

	return Track{ti, al, au, st, ar, aa}

}

// Output track information
func Output(t Track) {
	fmt.Printf(
		"Title:     %v\nAlbum:     %v\nArtist(s): %v\nPlayback:  %v\n",
		t.Title, t.Album, strings.Join(t.AlbumArtist[:], ", "), t.Status,
	)
}

// Output a JSON formatted string for use with Waybar's 'custom' module
func OutputWaybar(t Track) {
	var text, tooltip string
	if t.Status == "Stopped" {
		text = "Not Playing"
		tooltip = "It's quiet..."
	} else {

		text = fmt.Sprintf("%v • %v", t.Title, t.AlbumArtist[0])
		tooltip = fmt.Sprintf(
			"Title:     %v\nAlbum:     %v\nArtist(s): %v\nPlayback:  %v",
			t.Title, t.Album, strings.Join(t.Artist[:], ", "), t.Status,
		)

		if t.Status == "Playing" {
			text += "   "
		} else if t.Status == "Paused" {
			text += "   "
		}
	}

	wb, _ := json.Marshal(Waybar{
		Text:    text,
		Tooltip: tooltip,
		Class:   "$class",
	})

	fmt.Println(string(wb))

}

func main() {

	op := flag.String("o", "", "Formatting for output: 'Waybar', 'None' (default)")
	cmd := flag.String("c", "", "Commands: 'PlayPause', 'Stop', 'Next', 'Previous'")
	flag.Parse()

	// Connect to dbus Session Bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Println("Failed to connect to Session Bus ", err)
		os.Exit(1)
	}

	cm := map[string]bool{"playpause": true, "stop": true, "next": true, "previous": true}
	// check if cmd flag was passed
	if cm[strings.ToLower(*cmd)] {
		*cmd = strings.Title(strings.ToLower(*cmd)) // Ensure cmd is properly formatted
		if *cmd == "Playpause" {
			*cmd = "PlayPause"
		}
		SendCommand(conn, *cmd)
	} else { // if no cmd flags output track info
		pl := GetPlaying(conn)

		switch strings.ToLower(*op) {
		case "waybar":
			OutputWaybar(pl)
		case "none":
			Output(pl)
		default:
			Output(pl)
		}
	}
}
