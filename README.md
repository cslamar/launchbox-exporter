# LaunchBox Exporter

This is a small application that will allow you to export your ROMs, media, and metadata to a format that can be used by EmulationStation DE (so far the only tested frontend, it may work with others.)

_This is a work in progress_

## The Why

I love [LaunchBox](https://www.launchbox-app.com/), it does everything I need it to do (and more) for my desktop and arcade needs.  I recently got a SteamDeck though and wanted to have all the great things I have on my desktop, but on the go.  Enter [EmuDeck](https://github.com/dragoonDorise/EmuDeck), this is a great program/set of programs that enables all your retro needs on the SteamDeck, the only downside is that the scraper for metadata, images, and videos only works one game at a time.  

Also the filtering is not what I was looking for. I wanted a way to export the ROMs that I wanted to have with me at all time.  I wanted more than a hand full, but didn't need all 3700 SNES games with me at all times.  

LaunchBox handles all of this very well.  I can organize my games into playlists, regions, favorites, whatever I want.  The media pulled is top notch, and well, I already have it locally and didn't want to pull it again on the SteamDeck.  The naming conventions and metadata didn't match up exactly between LaunchBox and ES, so I wrote this translator/exporter to make it possible for me to get what I wanted.

Make sure you have a PAID copy of LaunchBox, because the work they put into it, and what you get out of it, is amazing and they deserve compensation for their hard work.

## Configuration

The exporter has a [single config file](config-dist.yaml) that handles all of the details.  It's in YAML format and is fairly easy to understand.  You can comment and uncomment different parts depending on what system(s) you're looking to export.

```yaml
launch_box_root: E:\Games\LaunchBox
output_directory: E:\Games\export-directory
platforms:
  "Nintendo Game Boy": gb
#  "Nintendo Game Boy Color": gbc
#  "Nintendo Entertainment System": nes
  "Nintendo Game Boy Advance": gba
  "Super Nintendo Entertainment System": snes
include_universal_games: true
regions:
  #  - Asia
  #  - China
  #  - Europe
  #  - France
  #  - Germany
  #  - Japan
  #  - Korea
  - North America
  #  - "North America, Europe"
  #  - Spain
  #  - "United Kingdom"
  - "United States"
#  - World

alt_video_locations:
  "Super Nintendo Entertainment System": "/Volumes/Junk/Games/LaunchBox/Videos/Super Nintendo Entertainment System"
  "Nintendo Entertainment System": "/Volumes/Junk/Games/LaunchBox/Videos/Nintendo Entertainment System"
  "Nintendo Game Boy Advance": "/Volumes/Junk/Games/LaunchBox/Videos/Nintendo Game Boy Advance"
```

The config reads like so:

* `launch_box_root` define the root of the LaunchBox application and data
* `output_directory` is where you want the exports to be saved
* The `platforms` map is a key value pairing where the key is the full name of the system that is in LaunchBox (note the `"` marks around the name of the system becuase of the spaces), and the value after the colon is the name of the system as it would be parsed by EmulationStation
* `include_universal_games` is the option to export games that aren't tied to a specific region
* The `regions` array is where you can specify what regions you wish to export.  The values MUST be the same as they appear in LaunchBox's directory hierarchy.
* The `alt_video_locations` is a temporary map of system name to filesystem locations that contain video previews for the games.  Note: this is a temporary option as I don't have an EmuMovies account and can't speak to how LaunchBox names saved videos that it automatically pulls.  This config option currently scrapes the specified directory and matches it to the name of the ROM in question.  This will probably be removed, augmented, or replaced with something that addresses how LaunchBox stores its video files.

## Running the Application

Once the config is in order, simply run the command with one, or more, of the following flags:

```shell
  -all-art
    	copy ALL found image media (this can take awhile)
  -box-art
    	copy found box art
  -config string
    	path to config file (default "config.yaml")
  -debug
    	enable debug logging
  -gamelist
    	build gamelist.xml output file
  -roms
    	copy found roms
  -videos
    	copy found videos
```

This will copy all the needed data to the output directory, and all you'll need to do is add the files to your SteamDeck's config folders.  (More information will come around how to do so, so stay tuned!)

The folders in question are:
* `~/.emulationstation/gamelists/` - this is where your newly genereated `gamelist.xml` files will go.
* `Emulation/tools/downloaded_media/` - this will be in the root of where you installed EmuDeck (mine is in the root of the SD card I installed to).  Inside are the different directories that hold media for each system.
* `Emulation/roms/` - this is where your ROMs will be stored.

## Notes

* This is a work in progress.  It's working fine for me and what I'm looking to do, but do so at your own warning.  You're source LaunchBox data and configs should not be touched, but running the application multiple times **will** overwrite the content in the output directory, so be warned!
* The majority of this has been written on a Mac using the LaunchBox root on a removable hard drive.  If there are issues, they probably have to do with how the pathing structure is different in *nix and Windows.  I've tested most of the functionality on a Windows PC and it worked fine.  The aim is to have this be usable cross-platform.
