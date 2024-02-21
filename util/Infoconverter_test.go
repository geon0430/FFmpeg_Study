package util

import (
	"testing"
	"reflect"
	"fmt"
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
			BufferSize: 99,
			Channels:   5,
			LogPath:    "/tmp/log",
			RtspServer: "rtsp://localtest:8554",
			// ENCODER:    "h264_encoder_test",
			// DECODER:    "h264_decoder_test",

		},
		// Encoder: Encoder{
		// 	H264: "h264_encoder_test",
		// 	H265: "h265_encoder_test",
		// },
		// Decoder: Decoder{
		// 	H264: "h264_decoder_test",
		// 	H265: "h265_decoder_test",
		// },
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
		RtspServer: "rtsp://localtest:8554",
		// ENCODER:     "h264_encoder_test",
		// DECODER:     "h264_decoder_test",
		OrgRtspAddr: "rtsp://localtest:8554/NAME_test_1888p",
		BufferSize: 99,
		Channels:   5,
		LogPath:    "/tmp/log",
	}
	resultRtspInfo := InfoConverter(TestConfig.RtspInfo)

    if !reflect.DeepEqual(resultRtspInfo, expectedRtspInfo) {
        t.Errorf("InfoConverter result does not match expected result")
        vResult := reflect.ValueOf(resultRtspInfo)
        vExpected := reflect.ValueOf(expectedRtspInfo)
        for i := 0; i < vResult.NumField(); i++ {
            if !reflect.DeepEqual(vResult.Field(i).Interface(), vExpected.Field(i).Interface()) {
                fmt.Printf("Field mismatch: %s\nExpected: %+v\nGot: %+v\n", 
                    vResult.Type().Field(i).Name, 
                    vExpected.Field(i).Interface(), 
                    vResult.Field(i).Interface())
            }
        }
    }
}