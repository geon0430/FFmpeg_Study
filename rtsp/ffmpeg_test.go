package rtsp

import (
	"context"
	"fmt"
	util "go_vms/src/pipeline/util"
	"strconv"
	"testing"
	"time"
)

func TestRtspHookServer(t *testing.T) {

	TestRtspInfo := util.RtspInfo{
		ID:             int(12345),
		NAME:           "NAME_test",
		RTSP:           "rtsp://admin:qazwsx123!@192.168.10.70/0/1080p/media.smp",
		CODEC:          "h264",
		MODEL:          "MODEL_test",
		FPS:            float64(30),
		IN_WIDTH:       int(1920),
		IN_HEIGHT:      int(1080),
		OUT_WIDTH:      int(1280),
		OUT_HEIGHT:     int(720),
		GPU:            int(0),
		ENCODER:        "h264_nvenc",
		DECODER:        "h264_cuvid",
		OrgRtspAddr:    "rtsp://localhost:8444/NAME_test_1080p",
		ResizeRtspAddr: "rtsp://localhost:8444/NAME_test_720p",

		BufferSize: 300,
		POS_FLAG:   float32(1),
		NEG_FLAG:   float32(2),
		WaitCnt:    int(1.0 / 10.0 * 5.0 / 0.01),
		ChunkSize:  int(1000),
		Channels:   int(5),
		FlagBytes:  int(8),
		LogPath:    string("/tmp/log"),
		//MaxPipelinesPerGPU int,
		GPU_NAME:     "3090",
		ControlImage: "/volume/go-test/go_vms/src/loading.jpg",
	}
	TestRtspInfo.ON_TIME, _ = time.Parse("15:04", "19:00")
	TestRtspInfo.OFF_TIME, _ = time.Parse("15:04", "7:00")

	//var TestModelControl ModelControl
	//TestModelControl.OnButtonChan <- false
	//TestModelControl.OffButtonChan <- true
	TestModelControl := util.ModelControl{
		OnButtonChan:  make(chan bool),
		OffButtonChan: make(chan bool),
	}
	go func() {
		for {
			select {
			case onVal := <-TestModelControl.OnButtonChan:
				fmt.Println("Received on OnButtonChan:", onVal)
			case offVal := <-TestModelControl.OffButtonChan:
				fmt.Println("Received on OffButtonChan:", offVal)
			}
		}
	}()
	TestModelControl.OnButtonChan <- false
	TestModelControl.OffButtonChan <- true

	TestPipelineInfo := util.PipelineInfo{
		RtspInfo:     TestRtspInfo,
		ModelControl: TestModelControl,
	}
	////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////
	MaxCounter := int(10)
	duration := time.Duration(1 * float64(time.Second))
	EndTicker := util.CreateTicker(duration)

	////////////////////////////////////////////////////////////////////////
	var Logpath = TestPipelineInfo.RtspInfo.LogPath
	var model = TestPipelineInfo.RtspInfo.MODEL
	var name = TestPipelineInfo.RtspInfo.NAME
	var idStr = TestPipelineInfo.RtspInfo.ID
	var id = strconv.Itoa(idStr)
	var logLevel = "debug"
	logger := util.SetupLogging(Logpath, model, name, logLevel, id)
	Context, cancel := context.WithCancel(context.Background())
	defer cancel()
	bufferSize := TestPipelineInfo.RtspInfo.BufferSize
	rtspReadBuf := make(chan *[]byte, bufferSize)
	errorChan := make(chan error, bufferSize)

	go ReadRtsp(
		Context,
		TestPipelineInfo,
		rtspReadBuf,
		errorChan,
		logger,
	)
	go StreamRtsp(
		Context,
		TestPipelineInfo,
		TestPipelineInfo.RtspInfo.IN_WIDTH, TestPipelineInfo.RtspInfo.IN_HEIGHT,
		TestPipelineInfo.RtspInfo.OrgRtspAddr,
		rtspReadBuf,
		errorChan,
		logger,
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
