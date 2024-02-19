package util

import (
	"testing"
	"reflect"
)

func TestInfoConverter(t *testing.T) {
	TestConfig := PipelineConfig{
		RtspInfo: RtspInfo{
			ID:        12345,
			NAME:      "NAME_test",
			RTSP:      "rtsp://test.test",
			CODEC:     "h264",
			MODEL:     "MODEL_test",
			FPS:       10,
			IN_WIDTH:  1999,
			IN_HEIGHT: 1888,
			GPU:       1,
		},
		General: GeneralConfig{
			BufferSize: 99,
			Channels:   5,
			LogPath:    "/tmp/log",
			RtspServer: "rtsp://localtest:8554",
		},
		Encoder: Encoder{
			H264: "h264_encoder_test",
			H265: "h265_encoder_test",
		},
		Decoder: Decoder{
			H264: "h264_decoder_test",
			H265: "h265_decoder_test",
		},
	}

	expectedRtspInfo := RtspInfo{
		ID:          12345,
		NAME:        "NAME_test",
		RTSP:        "rtsp://test.test",
		CODEC:       "h264",
		MODEL:       "MODEL_test",
		FPS:         10,
		IN_WIDTH:    1999,
		IN_HEIGHT:   1888,
		GPU:         1,
		ENCODER:     "h264_encoder_test",
		DECODER:     "h264_decoder_test",
		OrgRtspAddr: "rtsp://localtest:8554/NAME_test_1888p",
		BufferSize: 99,
		Channels:   5,
		LogPath:    "/tmp/log",
	}

	resultPipelineInfo := InfoConverter(TestConfig)

	if !reflect.DeepEqual(resultPipelineInfo.RtspInfo, expectedRtspInfo) {
		t.Errorf("InfoConverter result does not match expected result")
	}
}
