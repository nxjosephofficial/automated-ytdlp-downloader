package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/adrg/xdg"
)

func main() {

	ytdlpPath, err := check_ytdlp()
	if err != nil {
		log.Fatal(err)
	}

	check_dir(xdg.UserDirs.Music)
	check_dir(xdg.UserDirs.Videos)

	for {
		isPlaylist, link, err := getLink()
		if err != nil {
			log.Fatal(err)
		}

		if link == "" {
			break
		}

		ytdlpArgs, err := getArgs(isPlaylist)
		if err != nil {
			log.Fatal(err)
		}

		downloadLink(ytdlpPath, ytdlpArgs, link)
	}
}

func check_ytdlp() (string, error) {
	path, err := exec.LookPath("yt-dlp")
	if err != nil {
		return "", err
	}
	return path, nil
}

func check_dir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func getLink() (bool, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter link: ")
	link, err := reader.ReadString('\n')
	if err != nil {
		return false, "", err
	}
	var isPlaylist bool
	if strings.Contains(link, "playlist") {
		isPlaylist = true
	}
	return isPlaylist, strings.TrimSpace(link), nil
}

func downloadLink(ytdlpPath string, ytdlpArgs []string, link string) (string, error) {
	cmd := exec.Command(ytdlpPath, append(ytdlpArgs, link)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil
	}
	return string(output), nil
}

func getArgs(isPlaylist bool) ([]string, error) {
	var ytdlpArgs []string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose type: 1) Audio\t2) Video\n")
	contentType, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	contentType = strings.TrimSpace(contentType)
	if contentType == "1" {
		path := xdg.UserDirs.Music + "/%(title)s.%(ext)s"
		var format string
		fmt.Print("Choose format:\n1) mp3\t2) m4a\t3) wav\n4) flac\t5) opus\t6) ogg\n")
		contentFormat, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		contentFormat = strings.TrimSpace(contentFormat)
		switch contentFormat {
		case "1":
			format = "mp3"
		case "2":
			format = "m4a"
		case "3":
			format = "wav"
		case "4":
			format = "flac"
		case "5":
			format = "opus"
		case "6":
			format = "vorbis"
		default:
			return nil, errors.New("invalid content format")
		}
		if isPlaylist {
			path := xdg.UserDirs.Music
			ytdlpArgs = []string{"-x", "--audio-format", format, "--output", fmt.Sprintf("%s/%%(playlist|)s/%%(playlist_index)s - %%(title)s.%%(ext)s", path)}
		} else {
			ytdlpArgs = []string{"-x", "--audio-format", format, "--output", path}
		}
	} else if contentType == "2" {
		path := xdg.UserDirs.Videos + "/%(title)s.%(ext)s"
		ytdlpArgs = []string{"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best", "--output", path}
	} else {
		return nil, errors.New("invalid download type")
	}
	return ytdlpArgs, nil
}
