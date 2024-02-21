package rtsp

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GPUVideoFFmpegRead(
	pipelineConfig util.PipelineInfo,
	logger *logrus.Logger,
	errorChan chan<- error,
	ctx context.Context,
) (*exec.Cmd, io.ReadCloser, error) {

	videoFilePath := pipelineConfig.RtspInfo.RTSP 
	FPS := pipelineConfig.RtspInfo.FPS
	In_width := pipelineConfig.RtspInfo.IN_WIDTH
	In_height := pipelineConfig.RtspInfo.IN_HEIGHT

	cmd := exec.Command("ffmpeg",
		"-hwaccel", "cuda",
		"-i", videoFilePath,
		"-preset", "fast",
		"-framerate", fmt.Sprintf("%f", FPS),
		"-vf", fmt.Sprintf("scale=%d:%d", In_width, In_height),
		"-f", "rawvideo",
		"-pix_fmt", "bgr24",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf("Failed to create stdout pipe: %v", err)
		errorChan <- err
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		logger.Errorf("Failed to start ffmpeg: %v", err)
		errorChan <- err
		return nil, nil, err
	}
	logger.Debugf("GPUFFmpegRead started successfully with video file")

	return cmd, stdout, nil
}

func CPUVideoFFmpegRead(
	pipelineConfig util.PipelineInfo,
	logger *logrus.Logger,
	errorChan chan<- error,
	ctx context.Context,
) (*exec.Cmd, io.ReadCloser, error) {

	videoFilePath := pipelineConfig.RtspInfo.RTSP 
	FPS := pipelineConfig.RtspInfo.FPS
	In_width := pipelineConfig.RtspInfo.IN_WIDTH
	In_height := pipelineConfig.RtspInfo.IN_HEIGHT

	cmd := exec.Command("ffmpeg",
		"-i", videoFilePath,
		"-preset", "fast",
		"-framerate", fmt.Sprintf("%f", FPS),
		"-vf", fmt.Sprintf("scale=%d:%d", In_width, In_height),
		"-f", "rawvideo",
		"-pix_fmt", "bgr24",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf("Failed to create stdout pipe: %v", err)
		errorChan <- err
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		logger.Errorf("Failed to start ffmpeg: %v", err)
		errorChan <- err
		return nil, nil, err
	}
	logger.Debugf("CPUFFmpegRead started successfully with video file")

	return cmd, stdout, nil
}
