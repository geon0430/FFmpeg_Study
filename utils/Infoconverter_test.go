package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInfoConverter(t *testing.T) {
	TestApiConfig := apipkg.RTSPstruct{
		ID:          int(12345),
		NAME:        "NAME_test",
		RTSP:        "rtsp://test.test",
		CODEC:       "h264",
		MODEL:       "MODEL_test",
		FPS:         float64(10),
		IN_WIDTH:    1999,
		IN_HEIGHT:   1888,
	}
	TestGeneralConfig := GeneralConfig{
		BufferSize:         99,
		Channels:           5,
		LogPath:            "/tmp/log",
		RtspServer:         "rtsp://localtest:8554",
	}
	TestEncoder := Encoder{
		H264: "h264_encoder_test",
		H265: "h265_encoder_test",
	}
	TestDecoder := Decoder{
		H264: "h264_decoder_test",
		H265: "h265_decoder_test",
	}
	TestConfig := PipelineConfig{
		General: TestGeneralConfig,
		Encoder: TestEncoder,
		Decoder: TestDecoder,
	}

	TestRtspInfo := RtspInfo{
		ID:             int(12345),
		NAME:           "NAME_test",
		RTSP:           "rtsp://test.test",
		CODEC:          "h264",
		MODEL:          "MODEL_test",
		FPS:            float64(10),
		IN_WIDTH:       int(1999),
		IN_HEIGHT:      int(1888),
		GPU:            int(1),
		ENCODER:        "h264_encoder_test",
		DECODER:        "h264_decoder_test",
		OrgRtspAddr:    "rtsp://localtest:8554/NAME_test_1888p",

		BufferSize:   99,
		Channels:     int(5),
		LogPath:      string("/tmp/log"),
	}


	TestPipelineInfo := PipelineInfo{
		RtspInfo: TestRtspInfo,
		// ModelControl: TestModelControl,
	}

	TestpipelineInfo := InfoConverter(TestApiConfig, TestConfig)
	if TestpipelineInfo.RtspInfo != TestPipelineInfo.RtspInfo {
		t.Error("Wrong result")
		vResult := reflect.ValueOf(TestpipelineInfo.RtspInfo)
		vTruth := reflect.ValueOf(TestPipelineInfo.RtspInfo)
		for i := 0; i < vResult.NumField(); i++ {
			if vResult.Field(i).Interface() != vTruth.Field(i).Interface() {
				fmt.Println(vResult.Type().Field(i).Name, " : ", vResult.Field(i).Interface())
				fmt.Println(vTruth.Type().Field(i).Name, " : ", vTruth.Field(i).Interface())
			}
		}

	}
}

