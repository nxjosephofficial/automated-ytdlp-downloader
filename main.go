package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	ytdlpPath, err := check_ytdlp()
	if err != nil {
		log.Fatal(err)
	}

	link, err := getLink()
	if err != nil {
		log.Fatal(err)
	}

	ytdlpArgs, err := getArgs()
	if err != nil {
		log.Fatal(err)
	}

	downloadLink(ytdlpPath, ytdlpArgs, link)
}

func check_ytdlp() (string, error) {
	path, err := exec.LookPath("yt-dlp")
	if err != nil {
		return "", err
	}
	return path, nil
}

func getLink() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter link: ")
	link, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(link), nil
}

func downloadLink(ytdlpPath string, ytdlpArgs []string, link string) (string, error) {
	cmd := exec.Command(ytdlpPath, append(ytdlpArgs, link)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil
	}
	return string(output), nil
}

func getArgs() ([]string, error) {
	var ytdlpArgs []string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose type: 1) Audio\t2) Video\n")
	contentType, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	contentType = strings.TrimSpace(contentType)
	if contentType == "1" {
		ytdlpArgs = []string{"-x", "--audio-format", "mp3"}
	} else if contentType == "2" {
		ytdlpArgs = []string{"-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"}
	} else {
		return nil, errors.New("invalid download type")
	}
	return ytdlpArgs, nil
}
