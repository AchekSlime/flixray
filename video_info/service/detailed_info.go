package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type VideoInfoDetailed struct {
	FileName        string `json:"file_name"`
	Duration        string `json:"duration"`
	CurrentDuration int64  `json:"current_duration"`
}

func DetailedInfoController(ctx *gin.Context) {
	fileName := ctx.Query("file_name")
	detailedInfo, _ := GetDetailedVideoInfo(fileName)

	ctx.IndentedJSON(200, detailedInfo)
}

func GetDetailedVideoInfo(fileName string) (*VideoInfoDetailed, error) {
	content, err := os.ReadFile(fmt.Sprintf("./%s.json", fileName))
	if err != nil {
		logrus.Error("Error when opening file: ", err)
		return nil, err
	}

	var payload VideoInfoDetailed
	err = json.Unmarshal(content, &payload)
	if err != nil {
		logrus.Error("Error during Unmarshal(): ", err)
		return nil, err
	}

	floatDuration, _ := strconv.ParseFloat(payload.Duration, 64)
	floatDuration = floatDuration * 1000
	payload.CurrentDuration = int64(floatDuration)

	return &payload, nil
}
