package util

import (
	"fmt"
	// "reflect"
	// "strings"

)

func InfoConverter(rtspInfo RtspInfo) RtspInfo {
	var pipelineInfo RtspInfo
	pipelineInfo.ID = rtspInfo.ID
	pipelineInfo.NAME = rtspInfo.NAME
	pipelineInfo.RTSP = rtspInfo.RTSP
	pipelineInfo.CODEC = rtspInfo.CODEC
	pipelineInfo.MODEL = rtspInfo.MODEL
	pipelineInfo.FPS = rtspInfo.FPS
	pipelineInfo.IN_WIDTH = rtspInfo.IN_WIDTH
	pipelineInfo.IN_HEIGHT = rtspInfo.IN_HEIGHT
	pipelineInfo.RtspServer = rtspInfo.RtspServer

	// pipelineInfo.ENCODER = returnEncoder(rtspInfo.CODEC, rtspInfo)
	// pipelineInfo.DECODER = returnDecoder(rtspInfo.CODEC, rtspInfo)
	_orgRtspAddr := returnStreamAddr(rtspInfo) 
	pipelineInfo.OrgRtspAddr = _orgRtspAddr
	pipelineInfo.BufferSize = rtspInfo.BufferSize
	pipelineInfo.Channels = rtspInfo.Channels
	pipelineInfo.LogPath = rtspInfo.LogPath

	return pipelineInfo
}


// func returnEncoder(CODEC string, rtspInfo RtspInfo) string {
//     codec := strings.ToUpper(CODEC) // codec = H264

//     r := reflect.ValueOf(rtspInfo.Encoder)
//     encoder := reflect.Indirect(r).FieldByName(codec)
//     return encoder.String()
// }

// func returnDecoder(CODEC string, rtspInfo RtspInfo) string {
//     codec := strings.ToUpper(CODEC) // codec = H264

//     r := reflect.ValueOf(rtspInfo.Decoder)
//     decoder := reflect.Indirect(r).FieldByName(codec)
//     return decoder.String()
// }


func returnStreamAddr(
	rtspInfo RtspInfo) (string) {
	name := rtspInfo.NAME
	rtspServer := rtspInfo.RtspServer
	in_height := rtspInfo.IN_HEIGHT

	OrgRtspAddr := fmt.Sprintf("%s/%s_%dp", rtspServer, name, in_height)

	return OrgRtspAddr
}
