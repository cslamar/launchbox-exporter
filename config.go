package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// media map between LaunchBox and ES folder mapping
type media struct {
	Lb string
	Es string
}

// MediaMatches LaunchBox and ES folder mapping group
var MediaMatches []media

// Config options structure
type Config struct {
	AltVideoLocation      string            `yaml:"alt_video_location"`
	AltVideoLocations     map[string]string `yaml:"alt_video_locations"`
	LaunchBoxRoot         string            `yaml:"launch_box_root"`
	IncludeUniversalGames bool              `yaml:"include_universal_games"`
	OutputDirectory       string            `yaml:"output_directory"`
	Platforms             map[string]string `yaml:"platforms"`
	Regions               []string          `yaml:"regions"`
}

// copyAltVideoFiles copy to configured output directory
func (c *Config) copyAltVideoFiles(platform string, esGame EsGame) error {
	outputDir := filepath.Join(c.OutputDirectory, c.Platforms[platform], "media", "videos")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	fpath := filepath.Join(c.AltVideoLocations[platform], (esGame.ImagePath + "*"))
	vidMatches, err := filepath.Glob(fpath)
	if err != nil {
		return errors.New(fmt.Sprintf("bad pattern matching: %v", err))
	}
	if len(vidMatches) == 0 {
		return errors.New("no matches found for: " + esGame.Name)
	}

	inFile, err := os.Open(vidMatches[0])
	if err != nil {
		return err
	}
	defer inFile.Close()
	inputFile, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(outputDir, filepath.Base(vidMatches[0])), inputFile, 0644); err != nil {
		return err
	}
	return nil
}

// copyRomFiles copy to configured output directory
func (c *Config) copyRomFiles(platform string, lbGame LbGame) error {
	outputDir := filepath.Join(c.OutputDirectory, c.Platforms[platform], "roms")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Detect OS to determine how paths should be parsed
	var romPath string
	if runtime.GOOS == "windows" {
		romPath = lbGame.ApplicationPath
	} else {
		romPath = strings.Replace(lbGame.ApplicationPath, `\`, "/", -1)
	}

	inFile, err := os.Open(filepath.Join(c.LaunchBoxRoot, romPath))
	if err != nil {
		return err
	}
	defer inFile.Close()
	inputFile, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(outputDir, lbGame.EsGame.Path), inputFile, 0644); err != nil {
		return err
	}

	return nil
}

// copyBoxArtFiles copy to configured output directory
func (c *Config) copyBoxArtFiles(platform string, lbGame LbGame) error {
	outputDir := filepath.Join(c.OutputDirectory, c.Platforms[platform], "media", EsCovers)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	if err := lbGame.SetImagePath(*c); err != nil {
		return err
	}

	inFile, err := os.Open(lbGame.ImagePath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	inputFile, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(outputDir, lbGame.EsGame.ImagePath), inputFile, 0644); err != nil {
		return err
	}

	return nil
}

// copyAllArtFiles copies all found art to output directory
func (c *Config) copyAllArtFiles(platform string, lbGame LbGame) error {
	for _, mediaType := range MediaMatches {
		srcFile, destFile, err := lbGame.ImageCategoryConfiguration(*c, mediaType.Lb)
		if err != nil {
			log.Debugln(err)
			continue
		}

		outputDir := filepath.Join(c.OutputDirectory, c.Platforms[platform], "media", mediaType.Es)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}

		outfile := filepath.Join(outputDir, destFile)
		if srcFile != "" && outfile != "" {
			log.Debugln("srcFile:", srcFile)
			log.Debugln("outfile:", outfile)
			log.Infof("copying %s -> %s\n", srcFile, outfile)
			if err := copyFile(srcFile, outfile); err != nil {
				log.Errorln(err)
			}
		}
	}

	return nil
}

// NewConfig creates new Config object
func NewConfig(configFile string) (config Config, err error) {
	fp, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}

	data, err := io.ReadAll(fp)
	if err != nil {
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	// expose media matchers
	MediaMatches = []media{
		{
			Es: EsBoxes,
			Lb: LbBoxes,
		},
		{
			Es: EsBackCovers,
			Lb: LbBackCovers,
		},
		{
			Es: EsCovers,
			Lb: LbCovers,
		},
		{
			Es: EsFanArt,
			Lb: LbFanArt,
		},
		{
			Es: EsMarquees,
			Lb: LbMarquees,
		},
		{
			Es: EsPhysicalMedia,
			Lb: LbPhysicalMedia,
		},
		{
			Es: EsScreenShots,
			Lb: LbScreenShots,
		},
		{
			Es: EsTitleScreens,
			Lb: LbTitleScreens,
		},
	}

	return config, nil
}
