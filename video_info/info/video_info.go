package info

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var FFProbePath = "ffprobe"

type VideoInfo struct {
	FileName string `json:"file_name"`
	Duration string `json:"duration"`
}

func GenerateVideoInfoConfig(path string) {
	infoMap, err := getFFMPEGJson(path)
	if err != nil {
		logrus.Error(err)
		return
	}
	format := infoMap["format"].(map[string]interface{})

	videoInfo := VideoInfo{
		FileName: format["filename"].(string),
		Duration: format["duration"].(string),
	}

	infoJson, _ := json.MarshalIndent(videoInfo, "", " ")
	//err = os.WriteFile(fmt.Sprintf("%s.json", videoInfo.FileName[:len(videoInfo.FileName)-4]), infoJson, 0644)
	err = os.WriteFile(fmt.Sprintf("%s.json", videoInfo.FileName), infoJson, 0644)
	if err != nil {
		logrus.Error(err)
		return
	}
	//fmt.Println(string(infoJson))
}

func getFFMPEGJson(path string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, FFProbePath, "-v", "quiet", "-print_format", "json", "-show_format", path)
	var info map[string]interface{}
	if err := execAndGetStdoutJson(cmd, &info); err != nil {
		return nil, fmt.Errorf("error getting JSON from ffprobe output for file '%v': %v", path, err)
	}

	return info, nil
}

func execAndGetStdoutJson(cmd *exec.Cmd, v interface{}) error {
	b, err := execAndGetStdoutBytes(cmd)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func execAndGetStdoutBytes(cmd *exec.Cmd) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := execAndWriteStdout(cmd, b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func execAndWriteStdout(cmd *exec.Cmd, w io.Writer) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error opening stdout of command: %v", err)
	}
	defer stdout.Close()
	logrus.Debugf("Executing: %v %v", cmd.Path, cmd.Args)
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}
	if _, err := io.Copy(w, stdout); err != nil {
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		return fmt.Errorf("error copying stdout to buffer: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed %v", err)
	}
	return nil
}
