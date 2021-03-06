// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/goh-chunlin/go-onedrive/onedrive"
)

type Album struct {
	Title        string
	MusicItems   []MusicItem
	ThumbnailURL string
}

type MusicItem struct {
	Id                string
	Title             string
	Description       string
	Album             string
	AlbumArtist       string
	AlbumThumbnailURL string
	Duration          int
	DurationDisplay   string
	DownloadURL       string
}

type Thumbnail struct {
	Id          string
	Description string
	DownloadURL string
}

func showMusicListPage(context *gin.Context) {
	var isLoggedIn = false

	var titleOutput = "Welcome"
	token := getAccessAndRefreshTokenFromCookie(context)

	if token != nil {
		tc := oauthConfig.Client(context, token)

		client := onedrive.NewClient(tc)

		defaultDrive, err := client.Drives.Default(context)
		if err == nil && defaultDrive.Id != "" {
			isLoggedIn = true

			titleOutput = "Welcome back, " + defaultDrive.Owner.User.DisplayName

			oneDriveMusicDirectory := getItemsInMusicDirectory(context, client)
			if oneDriveMusicDirectory != nil {
				albumThumbnails := processThumbnails(context, client, oneDriveMusicDirectory.DriveItems)

				musicItems := processMusicItem(context, client, oneDriveMusicDirectory.DriveItems, albumThumbnails)

				albums := createAlbums(musicItems)

				context.HTML(http.StatusOK, "music-list.tmpl.html", gin.H{
					"isLoggedIn": isLoggedIn,
					"title":      titleOutput,
					"albums":     albums,
				})

				return
			}
		}
	}

	context.HTML(http.StatusOK, "index.tmpl.html", gin.H{
		"isLoggedIn": isLoggedIn,
		"title":      titleOutput,
	})
}

func createAlbums(musicItems []MusicItem) []*Album {
	sort.Slice(musicItems[:], func(i, j int) bool {
		return musicItems[i].Album < musicItems[j].Album
	})

	currentAlbum := ""
	var albums []*Album
	var album *Album
	for _, musicItem := range musicItems {

		var albumMusicItems []MusicItem

		if currentAlbum != musicItem.Album {
			if album != nil {
				sort.Slice(album.MusicItems[:], func(i, j int) bool {
					return album.MusicItems[i].Title < album.MusicItems[j].Title
				})
			}

			currentAlbum = musicItem.Album

			albumMusicItems = append(albumMusicItems, musicItem)

			album = &Album{
				Title:        musicItem.Album,
				ThumbnailURL: musicItem.AlbumThumbnailURL,
				MusicItems:   albumMusicItems,
			}

			albums = append(albums, album)
		} else {
			album.MusicItems = append(album.MusicItems, musicItem)
		}
	}

	return albums
}
