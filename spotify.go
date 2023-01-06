package spotify

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"net/url"
	"encoding/base64"
	"strings"
	"io"
	"encoding/json"
	"time"
	"github.com/pkg/browser"
	"io/ioutil"
)

var config ConfigFile

func LoadConfig() {
	cFile, err := os.Open("spotify.conf")
	if err != nil {
		if(os.IsNotExist(err)) {
			log.Println("Config file does not exist, creating new one")
			log.Println("Enter client id: ")
			fmt.Scanf("%s", config.client_id)
			log.Println("Enter client secret: ")
			fmt.Scanf("%s", config.client_secret)
			log.Println("Enter redirect uri: (leave blank for default)")
			var redirect_uri string
			fmt.Scanf("%s", redirect_uri)
			if(redirect_uri != "") {
				config.redirect_uri = redirect_uri
			} else {
				config.redirect_uri = "http://localhost:8080/callback"
			}
			log.Println("Enter auth type (client or user): ")
			fmt.Scanf("%s", config.auth_type)
			saveConfig()
		} else {
			log.Fatal(err)
		}
	} else {
		defer cFile.Close()
		fmt.Fscanf(cFile,"client_id: %s", &config.client_id)
		fmt.Fscanf(cFile,"client_secret: %s", &config.client_secret)
		fmt.Fscanf(cFile,"token: %s", &config.token)
		fmt.Fscanf(cFile,"refresh_token: %s", &config.refresh_token)
		fmt.Fscanf(cFile,"expiry: %d", &config.expiry)
		fmt.Fscanf(cFile,"redirect_uri: %s", &config.redirect_uri)
		fmt.Fscanf(cFile,"auth_type: %s", &config.auth_type)
		fmt.Fscanf(cFile,"scope: %s", &config.scope)
		fmt.Fscanf(cFile,"enable_logs: %t", &config.enable_logs)
	}
	if(config.enable_logs) {
		log.Println("Loaded config file")
	}else
	{
		log.SetOutput(ioutil.Discard)
	}
	if(config.client_id == "" || config.client_secret == "") {
		log.Fatal("Error: client_id or client_secret not set")
	}

	if(config.token == "")	{
		log.Println("No token found, getting new one")
		if(config.auth_type == "client") {
			log.Println("Using client credentials")
			clientAuthorize()
		} else if(config.auth_type == "user") {
			log.Println("Using user credentials")
			userAuthorize()
		} else {
			log.Fatal("No auth type specified")
		}
	}
	if(config.expiry < int(time.Now().Unix())) {
		log.Println("Token expired, refreshing")
		if(config.auth_type == "client") {
			clientAuthorize()
		} else if(config.auth_type == "user") {
			RefreshToken()
		}
	}
}

func saveConfig(){
	log.Println("Saving config")
	cFile, err := os.Create("spotify.conf")
	if err != nil {
		log.Println("Error creating config file")
		log.Fatal(err)
	}
	defer cFile.Close()
	fmt.Fprintf(cFile,"client_id: %s\nclient_secret: %s\ntoken: %s\nrefresh_token: %s\nexpiry: %d\nredirect_uri: %s\nauth_type: %s\nscope: %s\nenable_logs: %t\n", config.client_id, config.client_secret, config.token, config.refresh_token, config.expiry, config.redirect_uri,config.auth_type,config.scope,config.enable_logs)
}

func RefreshToken() {
	log.Println("Refreshing token")
	refresh_url := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", config.refresh_token)
	req, err := http.NewRequest("POST", refresh_url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(config.client_id + ":" + config.client_secret)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result Auth_Response
	json.Unmarshal(body, &result)
	config.token = result.Access_token
	config.expiry = int(time.Now().Unix()) + result.Expires_in
	saveConfig()
}

func clientAuthorize() {
	log.Println("Obtaining token")
	req_url := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", req_url, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(config.client_id + ":" + config.client_secret)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error authorizing")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if(resp.StatusCode == 200) {
		body, _ := io.ReadAll(resp.Body)
		type Auth_Response struct {
			Access_token string `json:"access_token"`
			Token_type string `json:"token_type"`
			Expires_in int `json:"expires_in"`
		}
		var result Auth_Response
		json.Unmarshal(body, &result)
		config.token = result.Access_token
		config.expiry = result.Expires_in + int(time.Now().Unix())
		saveConfig()
		log.Println("Authorized")
	} else {
		log.Println("Error authorizing")
	}
}

func userAuthorize() {
	var server http.Server
	log.Println("Creating authorization url")
	req_url := "https://accounts.spotify.com/authorize"
	req_url += "?client_id=" + config.client_id
	req_url += "&response_type=code"
	if(config.redirect_uri == "") {
		req_url += "&redirect_uri=localhost:8080"
	} else {
		req_url += "&redirect_uri=" + config.redirect_uri
	}
	req_url += "&scope=" + config.scope
	log.Println("The following url will open in your browser: " + req_url)
	browser.OpenURL(req_url)
	log.Println("Waiting for callback")
	server.Handler = http.DefaultServeMux
	server.Addr = ":8080"
	var code string
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		log.Println("Code Aquired")
		defer server.Close()
	})
	server.ListenAndServe()
	codeToToken(code)
}

func codeToToken(code string) {
	req_url := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", config.redirect_uri)
	req, err := http.NewRequest("POST", req_url, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(config.client_id + ":" + config.client_secret)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error authorizing")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if(resp.StatusCode == 200) {
		body, _ := io.ReadAll(resp.Body)
		var result Auth_Response
		json.Unmarshal(body, &result)
		config.token = result.Access_token
		config.expiry = result.Expires_in + int(time.Now().Unix())
		config.refresh_token = result.Refresh_token
		saveConfig()
		log.Println("Authorized")
	} else {
		log.Println("Error authorizing")
		body, _ := io.ReadAll(resp.Body)
		log.Fatal(string(body))
	}
}

func GetCurrentSong() CurrentTrack {
	log.Println("Getting current song")
	req_url := "https://api.spotify.com/v1/me/player/currently-playing"
	req, err := http.NewRequest("GET", req_url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer " + config.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error getting current song")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if(resp.StatusCode == 200) {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Got current song")
		var result CurrentTrack
		json.Unmarshal(body, &result)
		return result
	} else {
		log.Fatal("Error getting current song")
		return CurrentTrack{}
	}
}

func GetPlaylist(playlist_id string) Playlist {
	log.Println("Getting playlist")
	req_url := "https://api.spotify.com/v1/playlists/" + playlist_id
	req, err := http.NewRequest("GET", req_url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer " + config.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if(resp.StatusCode == 200) {
		body, _ := io.ReadAll(resp.Body)
		var result Playlist
		log.Println(string(body))
		json.Unmarshal(body, &result)
		return result
	} else {
		log.Fatal("Error getting playlist")
		return Playlist{}
	}
}

func AddTrackToPlaylist(playlist_id string, track_id string) bool{
	log.Println("Adding track to playlist")
	req_url := "https://api.spotify.com/v1/playlists/" + playlist_id + "/tracks"
	req_url += "?uris=" + track_id
	req, err := http.NewRequest("POST", req_url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer " + config.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if(resp.StatusCode == 201) {
		log.Println("Added track to playlist")
		return true
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Error adding track to playlist")
		log.Fatal(string(body))
		return false
	}
}

func AddTracksToPlaylist(playlist_id string, track_ids []string) bool{
	log.Println("Adding tracks to playlist")
	req_url := "https://api.spotify.com/v1/playlists/" + playlist_id + "/tracks"
	req_url += "?uris=" + strings.Join(track_ids[:], ",")
	req, err := http.NewRequest("POST", req_url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer " + config.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if(resp.StatusCode == 201) {
		log.Println("Added track to playlist")
		return true
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Error adding track to playlist")
		log.Fatal(string(body))
		return false
	}
}
