package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type Playlist struct {
	Url    string
	Videos []Videos
}

type Videos struct {
	Id  string
	Url string
}

func unique(strs []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strs {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (u *Playlist) parseUrl() (string, error) {
	r, err := url.Parse(u.Url)
	if err != nil {
		log.Fatal(err)
	}
	m, _ := url.ParseQuery(r.RawQuery)
	return m["list"][0], err
}

func GetPlaylist(p *Playlist) (Playlist, int) {
	var result []string
	var videos []Videos
	id, err := p.parseUrl()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get("https://www.youtube.com/playlist?list=" + id)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var re = regexp.MustCompile(`(?mi)\"(?:videoId)\"\:\"(?P<id>.*?)\"`)
	for _, m := range re.FindAllStringSubmatch(string(body), -1) {
		//fmt.Println(m[1], i)
		result = append(result, m[1])
	}
	for _, r := range unique(result) {
		videos = append(videos, Videos{Url: fmt.Sprintf("https://www.youtube.com/watch?v=%s", r), Id: r})
	}
	return Playlist{Url: p.Url, Videos: videos}, len(unique(result))
}
