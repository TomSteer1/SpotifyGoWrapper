# SpotifyGoWrapper
## Installation and Usage
``` go 
import "github.com/tomsteer1/SpotifyGoWrapper"
```
The LoadConfig function will initalise the wrapper

## spotify.conf
```
client_id: <Insert Client ID>
client_secret: <Insert Client Secret>
token:
refresh_token: 
expiry:
redirect_uri: http://localhost:8080/callback
auth_type: <User or Client>
scope: <Insert Scopes seperated with %20>
enable_logs: <true/false>
```

## func LoadConfig
``` go 
func LoadConfig()
```
Loads the local file `spotify.conf`

## func RefreshToken
``` go 
func RefreshToken()
```
Refreshes the token if using the user auth type 

## func GetCurrentSong 
``` go 
func GetCurrentSong() CurrentTrack
```
Will return the currently playing track

## func GetPlaylist
``` go 
func GetPlaylist(playlist_id string) Playlist 
```
Gets the information about the playlist 

## func AddTrackToPlaylist
``` go 
func AddTrackToPlaylist(playlist_id string, track_id string) bool
```
Adds a single track of `track_id` to the playlist of `playlist_id`

## func AddTracksToPlaylist
``` go 
func AddTracksToPlaylist(playlist_id string, track_ids []string) bool
```
Adds multiple tracks to the playlist of `playlist_id`
