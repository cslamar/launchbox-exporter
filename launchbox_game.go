package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

// LaunchBox Media Directories
const (
	LbBoxes         = "Box - 3D"
	LbBackCovers    = "Box - Back"
	LbCovers        = "Box - Front"
	LbFanArt        = "Fanart - Background"
	LbMarquees      = "Arcade - Marquee"
	LbPhysicalMedia = "Cart - Front"
	LbScreenShots   = "Screenshot - Gameplay"
	LbTitleScreens  = "Screenshot - Game Title"
	LbVideos        = "Videos"
)

type LbGame struct {
	XMLName         xml.Name `xml:"Game"`
	ApplicationPath string   `xml:"ApplicationPath"`
	Developer       string   `xml:"Developer"`
	Genre           string   `xml:"Genre"`
	MaxPlayers      string   `xml:"MaxPlayers"`
	Notes           string   `xml:"Notes"`
	Platform        string   `xml:"Platform"`
	Publisher       string   `xml:"Publisher"`
	ReleaseDate     string   `xml:"ReleaseDate"`
	Title           string   `xml:"Title"`
	Region          string   `xml:"Region"`
	StarRating      float32  `xml:"CommunityStarRating"`
	EsGame          EsGame
	ImagePath       string
	VideoPath       string
	NixRomPath      string
}

// embedEsGame generates EsGame for struct
func (l *LbGame) embedEsGame() {
	// parse rom file path
	appPath := strings.Split(l.ApplicationPath, `\`)
	esFilePath := appPath[len(appPath)-1]

	// parse image naming
	ext := filepath.Ext(filepath.Base(l.ImagePath))
	esExt := filepath.Ext(esFilePath)
	esImg := fmt.Sprintf("%s%s", strings.TrimSuffix(esFilePath, esExt), ext)

	// convert time to ES formatted time
	releaseDate, err := time.Parse(time.RFC3339, l.ReleaseDate)
	if err != nil {
		log.Warnf("error parsing time for %s! %v", l.Title, err)
	}

	// convert MaxPlayers to ES format
	var players string
	if l.MaxPlayers == "" {
		log.Warnf("error processing %s: MaxPlayers is blank", l.Title)
	} else if l.MaxPlayers == "1" || l.MaxPlayers == "0" {
		players = "1"
	} else {
		players = fmt.Sprintf("1-%s", l.MaxPlayers)
	}

	// parse ratings
	esRating := fmt.Sprintf("%.2f", (l.StarRating / 5))

	// parse Genre
	esGenre := strings.Replace(l.Genre, ";", ",", -1)

	l.EsGame = EsGame{
		Path:        fmt.Sprintf("./%s", esFilePath),
		Name:        l.Title,
		Description: l.Notes,
		Developer:   l.Developer,
		Genre:       esGenre,
		Players:     players,
		Publisher:   l.Publisher,
		ImagePath:   esImg,
		Rating:      esRating,
		ReleaseDate: releaseDate.Format(EsTimeFormat),
		//VideoPath:   l.VideoPath,
	}
}

// SetImagePath sets image path from Config object
func (l *LbGame) SetImagePath(config Config) error {
	var imgPath string

	nameGlob := fmt.Sprintf("%s-01.*", scrubTitle(l.Title))
	if l.Region == "Universal" {
		imgPath = filepath.Join(config.LaunchBoxRoot, "Images", l.Platform, "Box - Front")
	} else {
		imgPath = filepath.Join(config.LaunchBoxRoot, "Images", l.Platform, "Box - Front", l.Region)
	}
	imgMatches, err := filepath.Glob(filepath.Join(imgPath, nameGlob))
	if err != nil {
		return errors.New(fmt.Sprintf("bad pattern matching: %v", err))
	}
	if len(imgMatches) == 0 {
		return errors.New("no matches found for: " + l.Title)
	}

	l.ImagePath = imgMatches[0]
	log.Debugln("found:", l.ImagePath)

	// Set embedded EsGame image path
	ext := filepath.Ext(filepath.Base(l.ImagePath))
	esExt := filepath.Ext(l.EsGame.Path)
	esImg := fmt.Sprintf("%s%s", strings.TrimSuffix(l.EsGame.Path, esExt), ext)
	l.EsGame.ImagePath = esImg

	return nil
}

// ImageCategoryConfiguration configures the output and image type of art
// takes in imgCategory which is a const value starting with 'Lb'
// returns the src and dest image paths/files
func (l *LbGame) ImageCategoryConfiguration(config Config, imgCategory string) (srcPath string, imgFilename string, err error) {
	var imgPath string

	nameGlob := fmt.Sprintf("%s-01.*", scrubTitle(l.Title))
	if l.Region == "Universal" {
		imgPath = filepath.Join(config.LaunchBoxRoot, "Images", l.Platform, imgCategory)
	} else {
		imgPath = filepath.Join(config.LaunchBoxRoot, "Images", l.Platform, imgCategory, l.Region)
	}

	log.Debugf("Source image path: %s", imgPath)

	imgMatches, err := filepath.Glob(filepath.Join(imgPath, nameGlob))
	if err != nil {
		log.Errorln(err)
		return "", "", errors.New(fmt.Sprintf("bad pattern matching: %v", err))
	}
	if len(imgMatches) == 0 {
		//log.Warnln("no matches found for: " + l.Title)
		return "", "", errors.New("no matches found for: " + l.Title)
	}

	srcPath = imgMatches[0]
	log.Debugln("found source image:", srcPath)
	srcExt := filepath.Ext(srcPath)
	destExt := filepath.Ext(l.EsGame.Path)
	destFile := strings.TrimPrefix(l.EsGame.Path, "./")
	destFile = strings.TrimSuffix(destFile, destExt)
	imgFilename = fmt.Sprintf("%s%s", destFile, srcExt)

	return
}
