package rtsp

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestFFmpegstream(t *testing.T) {

	TestRtspInfo := util.RtspInfo{
		ID:             int(12345),
		NAME:           "NAME_test",
		RTSP:           "rtsp_adress",
		CODEC:          "h264",
		MODEL:          "MODEL_test",
		FPS:            float64(30),
		IN_WIDTH:       int(1920),
		IN_HEIGHT:      int(1080),
		ENCODER:        "h264_nvenc",
		DECODER:        "h264_cuvid",
		OrgRtspAddr:    "rtsp://localhost:8444/NAME_test_1080p",

		BufferSize: 300,
		WaitCnt:    int(1.0 / 10.0 * 5.0 / 0.01),
		ChunkSize:  int(1000),
		Channels:   int(5),
		LogPath:    string("/tmp/log"),
	}
	TestRtspInfo.ON_TIME, _ = time.Parse("15:04", "19:00")
	TestRtspInfo.OFF_TIME, _ = time.Parse("15:04", "7:00")

	TestPipelineInfo := util.PipelineInfo{
		RtspInfo:     TestRtspInfo,
	}
	////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////
	MaxCounter := int(10)
	duration := time.Duration(1 * float64(time.Second))
	EndTicker := util.CreateTicker(duration)

	////////////////////////////////////////////////////////////////////////
	Context, cancel := context.WithCancel(context.Background())
	defer cancel()
	bufferSize := TestPipelineInfo.RtspInfo.BufferSize
	rtspReadBuf := make(chan *[]byte, bufferSize)
	errorChan := make(chan error, bufferSize)

	go ReadRtsp(
		Context,
		TestPipelineInfo,
		rtspReadBuf,
	)
	go StreamRtsp(
		Context,
		TestPipelineInfo,
		TestPipelineInfo.RtspInfo.IN_WIDTH, TestPipelineInfo.RtspInfo.IN_HEIGHT,
		TestPipelineInfo.RtspInfo.OrgRtspAddr,
		rtspReadBuf,
	)

	////////////////////////////////////////////////////////////////////////
	counter := int(1)
	for {
		select {
		case <-Context.Done():
			return
		case <-EndTicker.Ticker.C:
			fmt.Println(counter, "Second")
			if counter > MaxCounter {
				cancel()
			} else {
				counter += 1
			}
		}
	}
}
