package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			callerFuncParts := strings.Split(f.Function, "/")
			callerFunc := callerFuncParts[len(callerFuncParts)-1]
			return fmt.Sprintf("%s()", callerFunc), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})

	// default report caller to be off
	log.SetReportCaller(false)

	// enable colors for Windows
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	buildGameList := flag.Bool("gamelist", false, "build gamelist.xml output file")
	configFilePtr := flag.String("config", "config.yaml", "path to config file")
	copyAllArt := flag.Bool("all-art", false, "copy ALL found image media (this can take awhile)")
	copyBoxArt := flag.Bool("box-art", false, "copy found box art")
	copyRoms := flag.Bool("roms", false, "copy found roms")
	copyVideos := flag.Bool("videos", false, "copy found videos")
	debugPtr := flag.Bool("debug", false, "enable debug logging")

	flag.Parse()

	if *copyBoxArt && *copyAllArt {
		log.Fatalln("both -box-art and -all-art cannot be used together! Please specify one!")
	}

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}

	config, err := NewConfig(*configFilePtr)
	if err != nil {
		log.Fatalln(err)
	}

	startTime := time.Now()

	for lbPlatform, esPlatform := range config.Platforms {
		log.Infoln("Starting platform:", lbPlatform)
		metaPath := filepath.Join(config.LaunchBoxRoot, "Data", "Platforms", fmt.Sprintf("%s.xml", lbPlatform))
		fp, err := os.Open(metaPath)
		if err != nil {
			log.Fatalln(err)
		}

		decoder := xml.NewDecoder(fp)
		decoder.Strict = false

		type LaunchBox struct {
			Games []LbGame `xml:"Game"`
		}

		lbGames := LaunchBox{}

		if err := decoder.Decode(&lbGames); err != nil {
			log.Fatalln(err)
		}

		log.Println("number of parsed LaunchBox games:", len(lbGames.Games))
		log.Println("Total number of games:", len(lbGames.Games))

		foundGames := make([]LbGame, 0)
		for _, game := range lbGames.Games {
			if game.Region == "" && config.IncludeUniversalGames {
				game.Region = "Universal"
				game.embedEsGame()
				foundGames = append(foundGames, game)
			}
			if checkIfInSlice(game.Region, config.Regions) {
				game.embedEsGame()
				foundGames = append(foundGames, game)
			}
		}
		fp.Close()

		log.Println("number of filtered region games:", len(foundGames))

		if *buildGameList {
			log.Println("writing gamelist.xml for", lbPlatform)
			if err := WriteEsGameList(config, foundGames, lbPlatform); err != nil {
				log.Fatalln(err)
			}
		}

		if *copyRoms {
			log.Println("copying rom files for", lbPlatform)
			for idx, lbGame := range foundGames {
				fmt.Printf("(%s - roms) [%d/%d] ", esPlatform, idx+1, len(foundGames))
				if err := config.copyRomFiles(lbPlatform, lbGame); err != nil {
					log.Errorln("error copying rom:", err)
					continue
				}
				fmt.Println(lbGame.EsGame.Name)
			}
		}

		if *copyBoxArt {
			log.Println("copying box-art")
			for idx, lbGame := range foundGames {
				fmt.Printf("(%s - box-art) [%d/%d] ", esPlatform, idx+1, len(foundGames))
				if err := config.copyBoxArtFiles(lbPlatform, lbGame); err != nil {
					log.Warnln(err)
					continue
				}
				fmt.Println(lbGame.Title)
			}
		}

		if *copyAllArt {
			log.Infoln("copying all art assets")
			for idx, game := range foundGames {
				counter := fmt.Sprintf("(%s - all-art: %s) [%d/%d]:\n", esPlatform, game.Title, idx+1, len(foundGames))
				fmt.Printf("\x1b[%dm%s\x1b[0m", 32, counter)
				if err := config.copyAllArtFiles(lbPlatform, game); err != nil {
					log.Warnf("%s - %v", game.Title, err)
					continue
				}
			}
		}

		if *copyVideos {
			log.Infoln("copying videos")
			for idx, lbGame := range foundGames {
				fmt.Printf("(%s - videos) [%d/%d] ", esPlatform, idx+1, len(foundGames))
				if err := config.copyAltVideoFiles(lbPlatform, lbGame.EsGame); err != nil {
					log.Warnln(err)
					continue
				}
				fmt.Println(lbGame.EsGame.Name)
			}
		}
	}

	log.Printf("entire process took %s\n", time.Now().Sub(startTime).String())
}
