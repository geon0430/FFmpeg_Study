package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	apipkg "go_vms/src/api/util"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

func InfoConverter(
	apiConfig apipkg.RTSPstruct,
	globalConfig PipelineConfig) PipelineInfo {

	var pipelineInfo PipelineInfo
	pipelineInfo.RtspInfo.ID = apiConfig.ID
	pipelineInfo.RtspInfo.NAME = apiConfig.NAME
	pipelineInfo.RtspInfo.RTSP = apiConfig.RTSP
	pipelineInfo.RtspInfo.CODEC = apiConfig.CODEC
	pipelineInfo.RtspInfo.MODEL = apiConfig.MODEL
	pipelineInfo.RtspInfo.FPS = apiConfig.FPS
	pipelineInfo.RtspInfo.IN_WIDTH = apiConfig.IN_WIDTH
	pipelineInfo.RtspInfo.IN_HEIGHT = apiConfig.IN_HEIGHT
	pipelineInfo.RtspInfo.GPU = apiConfig.GPU

	pipelineInfo.RtspInfo.ENCODER = returnEncoder(apiConfig.CODEC, globalConfig)
	pipelineInfo.RtspInfo.DECODER = returnDecoder(apiConfig.CODEC, globalConfig)

	_orgRtspAddr, _resizeRtspAddr := returnStreamAddr(apiConfig, globalConfig)
	pipelineInfo.RtspInfo.OrgRtspAddr = _orgRtspAddr
	pipelineInfo.RtspInfo.ResizeRtspAddr = _resizeRtspAddr
	pipelineInfo.RtspInfo.BufferSize = globalConfig.General.BufferSize
	pipelineInfo.RtspInfo.Channels = globalConfig.General.Channels
	pipelineInfo.RtspInfo.LogPath = globalConfig.General.LogPath

	return pipelineInfo
}

func returnEncoder(CODEC string, globalConfig PipelineConfig) string {
	codec := strings.ToUpper(CODEC) // codec = H264

	r := reflect.ValueOf(globalConfig.Encoder)
	encoder := reflect.Indirect(r).FieldByName(codec)
	return encoder.String()
}

func returnDecoder(CODEC string, globalConfig PipelineConfig) string {
	codec := strings.ToUpper(CODEC) // codec = H264

	r := reflect.ValueOf(globalConfig.Decoder)
	decoder := reflect.Indirect(r).FieldByName(codec)
	return decoder.String()
}

func returnStreamAddr(
	apiconfig apipkg.RTSPstruct,
	globalConfig PipelineConfig) (string, string) {
	name := apiconfig.NAME
	rtspServer := globalConfig.General.RtspServer
	out_height := apiconfig.OUT_HEIGHT
	in_height := apiconfig.IN_HEIGHT

	OrgRtspAddr := fmt.Sprintf("rtsp://%s/%s_%dp", rtspServer, name, in_height)
	ResizeRtspAddr := fmt.Sprintf("rtsp://%s/%s_%dp", rtspServer, name, out_height)

	return OrgRtspAddr, ResizeRtspAddr
}
