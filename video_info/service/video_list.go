package service

import (
	"encoding/json"
	"github.com/achekslime/flixray/video_info/info"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func VideoListController(ctx *gin.Context) {
	list := make([]*info.VideoInfo, 0)

	// for
	minecraft, err := GetVideoInfo("./minecraft.mp4.json")
	if err != nil {
		logrus.Error("minecraft config file error: ", err)
		return
	}
	list = append(list, minecraft)

	//jsonList, err := json.MarshalIndent(list, "", " ")
	ctx.IndentedJSON(200, list)
}

func GetVideoInfo(path string) (*info.VideoInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		logrus.Error("Error when opening file: ", err)
		return nil, err
	}

	var payload info.VideoInfo
	err = json.Unmarshal(content, &payload)
	if err != nil {
		logrus.Error("Error during Unmarshal(): ", err)
		return nil, err
	}

	return &payload, nil
}
