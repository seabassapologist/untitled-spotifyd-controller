# Untitled Spotifyd Controller
A simple Go utility to display track information from, and send commands to, [spotifyd](https://github.com/Spotifyd/spotifyd) from Tiling Window Managers like Sway and i3

![spdcon running in Waybar in Sway](docs/spdcon_waybar.png)

```
Usage of spdcon:
  -c string
        Commands: 'PlayPause', 'Stop', 'Next', 'Previous'
  -o string
        Formatting for output: 'Waybar', 'None' (default)
```

Example Configuration for Waybar:
```json
"custom/spotifyd": {
    "format": " {}",
    "return-type": "json",
    "exec": "spdcon -o waybar",
    "on-click": "spdcon -c PlayPause",
    "on-click-right": "spdcon -c Next",
    "on-click-middle": "spdcon -c Previous",
    "interval": 2
}
```
### Known Issues:
* Podcast metadata won't display in output (See [spotifyd#1000](https://github.com/Spotifyd/spotifyd/issues/1000))
