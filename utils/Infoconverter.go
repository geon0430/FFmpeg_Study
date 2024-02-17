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
	RtspInfo  PipelineConfig) PipelineInfo {

	var pipelineInfo PipelineInfo
	pipelineInfo.RtspInfo.ID = RtspInfo.ID
	pipelineInfo.RtspInfo.NAME = RtspInfo.NAME
	pipelineInfo.RtspInfo.RTSP = RtspInfo.RTSP
	pipelineInfo.RtspInfo.CODEC = RtspInfo.CODEC
	pipelineInfo.RtspInfo.MODEL = RtspInfo.MODEL
	pipelineInfo.RtspInfo.FPS = RtspInfo.FPS
	pipelineInfo.RtspInfo.IN_WIDTH = RtspInfo.IN_WIDTH
	pipelineInfo.RtspInfo.IN_HEIGHT = RtspInfo.IN_HEIGHT
	pipelineInfo.RtspInfo.GPU = RtspInfo.GPU

	pipelineInfo.RtspInfo.ENCODER = returnEncoder(RtspInfo.CODEC, globalConfig)
	pipelineInfo.RtspInfo.DECODER = returnDecoder(RtspInfo.CODEC, globalConfig)

	_orgRtspAddr, _resizeRtspAddr := returnStreamAddr(RtspInfo, globalConfig)
	pipelineInfo.RtspInfo.OrgRtspAddr = _orgRtspAddr
	pipelineInfo.RtspInfo.ResizeRtspAddr = _resizeRtspAddr
	pipelineInfo.RtspInfo.BufferSize = RtspInfo.General.BufferSize
	pipelineInfo.RtspInfo.Channels = RtspInfo.General.Channels
	pipelineInfo.RtspInfo.LogPath = RtspInfo.General.LogPath

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
	RtspInfo PipelineConfig) (string, string) {
	name := RtspInfo.NAME
	rtspServer := globalConfig.General.RtspServer
	out_height := RtspInfo.OUT_HEIGHT
	in_height := RtspInfo.IN_HEIGHT

	OrgRtspAddr := fmt.Sprintf("rtsp://%s/%s_%dp", rtspServer, name, in_height)
	ResizeRtspAddr := fmt.Sprintf("rtsp://%s/%s_%dp", rtspServer, name, out_height)

	return OrgRtspAddr, ResizeRtspAddr
}
