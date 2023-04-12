package main

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// ES Media Directories
const (
	EsBoxes         = "3dboxes"
	EsBackCovers    = "backcovers"
	EsCovers        = "covers"
	EsFanArt        = "fanart"
	EsMarquees      = "marquees"
	EsMixImages     = "miximages"
	EsPhysicalMedia = "physicalmedia"
	EsScreenShots   = "screenshots"
	EsTitleScreens  = "titlescreens"
	EsVideos        = "videos"
)

type EsGame struct {
	XMLName     xml.Name `xml:"game"`
	Path        string   `xml:"path"`
	Name        string   `xml:"name"`
	Description string   `xml:"desc"`
	Rating      string   `xml:"rating,omitempty"`
	ReleaseDate string   `xml:"releasedate,omitempty"`
	Developer   string   `xml:"developer,omitempty"`
	Publisher   string   `xml:"publisher,omitempty"`
	Players     string   `xml:"players,omitempty"`
	Genre       string   `xml:"genre,omitempty"`
	//Image       string   `xml:"image,omitempty"`
	//Video       string   `xml:"video,omitempty"`
	ImagePath string
	VideoPath string
}

// WriteEsGameList creates a gamelist.xml file from parse LbGame list
func WriteEsGameList(config Config, lbGames []LbGame, platform string) error {
	esPlatform := filepath.Join(config.OutputDirectory, config.Platforms[platform])

	if err := os.MkdirAll(esPlatform, 0755); err != nil {
		return err
	}

	esGames := make([]EsGame, 0)
	for _, lbGame := range lbGames {
		esGame := lbGame.EsGame
		esGames = append(esGames, esGame)
	}

	log.Printf("original count: %d, es count: %d", len(lbGames), len(esGames))

	type container struct {
		XMLName  xml.Name `xml:"gameList"`
		GameList []EsGame
	}

	c := container{GameList: esGames}

	out, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	convertEscaped(&out)

	finalOutput := []byte(xml.Header + string(out))

	if err := os.WriteFile(filepath.Join(esPlatform, "gamelist.xml"), finalOutput, 0644); err != nil {
		return err
	}

	return nil
}

// parseEsImageName returns an image path for EsGame
func parseEsImagePath(srcImgPath, esPath string) string {
	ext := filepath.Ext(filepath.Base(srcImgPath))
	esExt := filepath.Ext(esPath)
	esImg := fmt.Sprintf("%s%s", strings.TrimSuffix(esPath, esExt), ext)

	return esImg
}
