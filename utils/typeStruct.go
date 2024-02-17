package util

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type PipelineConfig struct {
	General    RtspInfo
	Encoder    Encoder
	Decoder    Decoder
}

type Encoder struct {
	H264 string
	H265 string
}

type Decoder struct {
	H264 string
	H265 string
}

type RtspInfo struct {
	ID              int
	NAME            string
	RTSP            string
	CODEC           string
	MODEL           string
	FPS             float64
	IN_WIDTH        int
	IN_HEIGHT       int
	ENCODER         string
	DECODER         string
	OrgRtspAddr     string
	ResizeRtspAddr  string
	BufferSize int
	Channels   int
	LogPath    string
}

