// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goh-chunlin/go-onedrive/onedrive"
)

func getOneDriveOwnerUserId(context *gin.Context) (string, error) {
	var err error
	token := getAccessAndRefreshTokenFromCookie(context)

	if token != nil {
		tc := oauthConfig.Client(context, token)

		client := onedrive.NewClient(tc)

		defaultDrive, err := client.Drives.Default(context)
		if err == nil && defaultDrive.Id != "" {
			return defaultDrive.Owner.User.Id, nil
		}
	}

	return "", err
}

func getOneDriveItemDownloadUrlAndFileName(context *gin.Context, itemId string) (string, string, error) {
	var err error
	token := getAccessAndRefreshTokenFromCookie(context)

	if token != nil {
		tc := oauthConfig.Client(context, token)

		client := onedrive.NewClient(tc)

		driveItem, err := client.DriveItems.Get(context, itemId)
		if err == nil && driveItem.Id != "" {
			return driveItem.DownloadURL, driveItem.Name, nil
		}
	}

	return "", "", err
}

func getItemsInMusicDirectory(context *gin.Context, client *onedrive.Client) *onedrive.OneDriveDriveItemsResponse {
	defaultDriveRoot, err := client.DriveItems.List(context, "")
	if err == nil && defaultDriveRoot.DriveItems != nil {
		for _, item := range defaultDriveRoot.DriveItems {
			if item.Name == "Music" {
				oneDriveMusicDirectory, err := client.DriveItems.List(context, item.Id)
				if err == nil {
					return oneDriveMusicDirectory
				}
			}
		}
	}

	return nil
}

func processThumbnails(context *gin.Context, client *onedrive.Client, driveItems []*onedrive.DriveItem) []Thumbnail {
	var albumThumbnails []Thumbnail

	for _, musicItem := range driveItems {
		if musicItem.Name == "AlbumImages" {
			thumbnailDirectory, err := client.DriveItems.List(context, musicItem.Id)
			if err == nil && thumbnailDirectory.DriveItems != nil {
				for _, albumThumbnail := range thumbnailDirectory.DriveItems {
					if albumThumbnail.Image == nil {
						continue
					}

					thumbnail := Thumbnail{
						Id:          albumThumbnail.Id,
						Description: albumThumbnail.Description,
						DownloadURL: albumThumbnail.DownloadURL,
					}

					albumThumbnails = append(albumThumbnails, thumbnail)
				}
			}

			break
		}
	}

	return albumThumbnails
}

func processMusicItem(context *gin.Context, client *onedrive.Client, driveItems []*onedrive.DriveItem, albumThumbnails []Thumbnail) []MusicItem {
	var musicItems []MusicItem

	for _, musicItem := range driveItems {
		if musicItem.Audio != nil {
			var matchedAlbumThumbnailURL = ""
			for _, albumThumbnail := range albumThumbnails {
				if albumThumbnail.Description == "Default" {
					matchedAlbumThumbnailURL = albumThumbnail.DownloadURL
				} else if albumThumbnail.Description == musicItem.Audio.Album {
					matchedAlbumThumbnailURL = albumThumbnail.DownloadURL
					break
				}
			}

			music := MusicItem{
				Id:                musicItem.Id,
				Title:             musicItem.Audio.Title,
				Description:       musicItem.Description,
				Album:             musicItem.Audio.Album,
				AlbumArtist:       musicItem.Audio.AlbumArtist,
				AlbumThumbnailURL: matchedAlbumThumbnailURL,
				Duration:          musicItem.Audio.Duration,
				DurationDisplay:   getTimeDisplay(musicItem.Audio.Duration),
				DownloadURL:       musicItem.DownloadURL,
			}

			musicItems = append(musicItems, music)
		}
	}

	return musicItems
}
