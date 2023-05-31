package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoPlaylist struct {
	list map[string]*Video
}

func NewVideoPlaylist() *VideoPlaylist {
	return &VideoPlaylist{list: make(map[string]*Video, 0)}
}

func (playlist VideoPlaylist) StartVideo(ctx *gin.Context) {
	fileName := ctx.Query("file_name")
	detailedInfo, _ := GetDetailedVideoInfo(fileName)

	playlist.list[fileName] = NewVideo(detailedInfo.CurrentDuration)
	playlist.list[fileName].Start()
	ctx.IndentedJSON(200, gin.H{"duration": playlist.list[fileName].duration})
}

func (playlist VideoPlaylist) GetCurrentDuration(ctx *gin.Context) {
	fileName := ctx.Query("file_name")

	val, ok := playlist.list[fileName]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s not found", fileName)})
		ctx.Abort()
		return
	}

	duration := val.GetCurrentDuration()
	ctx.IndentedJSON(200, gin.H{"duration": duration})
}
