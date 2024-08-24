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

	ytdlpPath, err := check_dependency("yt-dlp")
	if err != nil {
		log.Fatal(err)
	}

	err = check_dir(xdg.UserDirs.Music)
	if err != nil {
		log.Fatal(err)
	}
	err = check_dir(xdg.UserDirs.Videos)
	if err != nil {
		log.Fatal(err)
	}

	for {
		isPlaylist, link, err := getLink()
		if err != nil {
			log.Fatal(err)
			continue
		}

		if link == "" {
			fmt.Println("No link entered, exiting.")
			break
		}

		ytdlpArgs, err := getArgs(isPlaylist)
		if err != nil {
			log.Fatal(err)
			continue
		}

		output, err := downloadLink(ytdlpPath, ytdlpArgs, link)
		if err != nil {
			log.Fatal(err)
		} else {
			if strings.Contains(output, "has already been downloaded") {
				fmt.Println("It has already been downloaded.")
			} else {
				fmt.Println("Download successful!")
			}
		}
	}
}

func check_dependency(dep string) (string, error) {
	path, err := exec.LookPath(dep)
	if err != nil {
		return "", errors.New("dependency is not found: " + dep)
	}
	return path, nil
}

func check_dir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return errors.New("couldn't check/create directory: " + dir)
		}
	}
	return nil
}

func getLink() (bool, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter link: ")
	link, err := reader.ReadString('\n')
	if err != nil {
		return false, "", errors.New("couldn't read link")
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
		return "", errors.New("couldn't download link")
	}
	return string(output), err
}

func getFormat(reader *bufio.Reader, prompt string, formatMap map[string]string) (string, error) {
	fmt.Print(prompt)
	contentFormat, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.New("couldn't read content format")
	}
	contentFormat = strings.TrimSpace(contentFormat)
	format, ok := formatMap[contentFormat]
	if !ok {
		return "", errors.New("invalid content format")
	}
	return format, nil
}

func getContentType(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	contentType, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.New("couldn't get content type")
	}
	contentType = strings.TrimSpace(contentType)
	return contentType, nil
}

func getArgs(isPlaylist bool) ([]string, error) {
	var ytdlpArgs []string
	reader := bufio.NewReader(os.Stdin)
	contentType, err := getContentType(reader, "Choose type: 1) Audio\t2) Video\n")
	if err != nil {
		return nil, err
	}
	if contentType == "1" {
		path := xdg.UserDirs.Music + "/%(title)s.%(ext)s"
		formatMap := map[string]string{
			"1": "mp3",
			"2": "m4a",
			"3": "wav",
			"4": "flac",
			"5": "opus",
			"6": "vorbis",
		}
		format, err := getFormat(reader, "Choose format:\n1) mp3\t2) m4a\t3) wav\n4) flac\t5) opus\t6) vorbis\n", formatMap)
		if err != nil {
			return nil, err
		}
		if isPlaylist {
			path := xdg.UserDirs.Music
			ytdlpArgs = []string{"-x", "--audio-format", format, "--output", fmt.Sprintf("%s/%%(playlist|)s/%%(playlist_index)s - %%(title)s.%%(ext)s", path)}
		} else {
			ytdlpArgs = []string{"-x", "--audio-format", format, "--output", path}
		}
	} else if contentType == "2" {
		path := xdg.UserDirs.Videos + "/%(title)s.%(ext)s"
		formatMap := map[string]string{
			"1": "mp4",
			"2": "mkv",
			"3": "webm",
		}
		format, err := getFormat(reader, "Choose format:\n1) mp4\t2) mkv\t3) webm\n", formatMap)
		if err != nil {
			return nil, err
		}
		if isPlaylist {
			path := xdg.UserDirs.Videos
			ytdlpArgs = []string{"--merge-output-format", format, "--output", fmt.Sprintf("%s/%%(playlist|)s/%%(playlist_index)s - %%(title)s.%%(ext)s", path)}
		} else {
			ytdlpArgs = []string{"--merge-output-format", format, "--output", path}
		}
	} else {
		return nil, errors.New("invalid download type")
	}
	return ytdlpArgs, nil
}
