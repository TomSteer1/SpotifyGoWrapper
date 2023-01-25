package spotify

type ConfigFile struct {
	client_id string
	client_secret string
	token string
	filePath string
	refresh_token string
	expiry int
	redirect_uri string
	auth_type string
	scope string
	enable_logs bool
}

type CurrentTrack struct {
	Item Track `json:"item"`
	Progress_ms int `json:"progress_ms"`
	Timestamp int `json:"timestamp"`
	Is_playing bool `json:"is_playing"`
}

type Track struct {
	Album Album `json:"album"`
	Artists []Artist `json:"artists"`
	Disc_number int `json:"disc_number"`
	Duration_ms int `json:"duration_ms"`
	Explicit bool `json:"explicit"`
	External_ids struct {
		Isrc string `json:"isrc"`
	} `json:"external_ids"`
	External_urls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href string `json:"href"`
	Id string `json:"id"`
	Is_local bool `json:"is_local"`
	Name string `json:"name"`
	Popularity int `json:"popularity"`
	Preview_url string `json:"preview_url"`
	Track_number int `json:"track_number"`
	Type string `json:"type"`
	Uri string `json:"uri"`
}

type Artist struct {
	External_urls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href string `json:"href"`
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Uri string `json:"uri"`
}

type Auth_Response struct {
	Access_token string `json:"access_token"`
	Token_type string `json:"token_type"`
	Expires_in int `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
}

type Playlist struct {
	Href string `json:"href"`
	Id string `json:"id"`
	Name string `json:"name"`
	Owner struct {
		Display_name string `json:"display_name"`
		External_urls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		Id string `json:"id"`
		Type string `json:"type"`
		Uri string `json:"uri"`
	} `json:"owner"`
	Public bool `json:"public"`
	Collaborative bool `json:"collaborative"`
	Followers struct {
		Href interface{} `json:"href"`
		Total int `json:"total"`
	} `json:"followers"`
	Images []Image `json:"images"`
	Tracks struct {
		Href string `json:"href"`
		Total int `json:"total"`
		Item []struct {
			Added_at string `json:"added_at"`
			Added_by struct {
				External_urls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				Id string `json:"id"`
				Type string `json:"type"`
				Uri string `json:"uri"`
			} `json:"added_by"`
			Is_local bool `json:"is_local"`
			Track Track `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
	Type string `json:"type"`
	Uri string `json:"uri"`
}

type Album struct {
	Album_type string `json:"album_type"`
	Artists []Artist `json:"artists"`
	External_urls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href string `json:"href"`
	Id string `json:"id"`
	Images []Image `json:"images"`
	Name string `json:"name"`
	Release_date string `json:"release_date"`
	Release_date_precision string `json:"release_date_precision"`
	Total_tracks int `json:"total_tracks"`
	Type string `json:"type"`
	Uri string `json:"uri"`
}

type Image struct {
	Height int `json:"height"`
	Url string `json:"url"`
	Width int `json:"width"`
}
