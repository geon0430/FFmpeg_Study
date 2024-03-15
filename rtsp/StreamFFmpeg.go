package rtsp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GPUFFmpegStream(
	ctx context.Context,
	gpuDevice int,
	pipelineConfig util.PipelineInfo,
	videoWidth, videoHeight int,
	streamUrl string,
	logger *logrus.Logger,
) (*exec.Cmd, io.WriteCloser, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	FPS := pipelineConfig.RtspInfo.FPS
	encoder := pipelineConfig.RtspInfo.ENCODER

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-re",
		"-hwaccel", "cuda",
		"-hwaccel_device", fmt.Sprintf("%d", gpuDevice),
		"-f", "rawvideo",
		"-pixel_format", "bgr24",
		"-video_size", fmt.Sprintf("%dx%d", videoWidth, videoHeight),
		"-framerate", fmt.Sprintf("%f", FPS),
		"-i", "-",
		"-timeout", "300000",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "300",
		"-vf", fmt.Sprintf("scale=%d:%d", videoWidth, videoHeight),
		"-c:v", fmt.Sprint(encoder), //필요함
		"-preset", "p3",
		"-tune", "hq",
		"-maxrate", "4000k",
		"-g", "60",
		"-f", "rtsp",
		"-rtsp_transport", "tcp",
		fmt.Sprintf(streamUrl),
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Errorf("Error creating stdin pipe: %v", err)
		return nil, nil, err
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		logger.Errorf("Error starting command: %v", err)
		return nil, nil, err
	}

	logger.Infof("Streaming command started successfully")
	return cmd, stdin, nil
}

func CPUFFmpegStream(
	ctx context.Context,
	gpuDevice int,
	pipelineConfig util.PipelineInfo,
	videoWidth, videoHeight int,
	streamUrl string,
	logger *logrus.Logger,
) (*exec.Cmd, io.WriteCloser, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	FPS := pipelineConfig.RtspInfo.FPS
	// encoder := pipelineConfig.RtspInfo.ENCODER

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-re",
		"-f", "rawvideo",
		"-pixel_format", "bgr24",
		"-video_size", fmt.Sprintf("%dx%d", videoWidth, videoHeight),
		"-framerate", fmt.Sprintf("%f", FPS),
		"-i", "-",
		"-vf", fmt.Sprintf("scale=%d:%d", videoWidth, videoHeight),
		"-c:v", "libx264",
		"-timeout", "300000",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "300",
		"-preset", "fast",
		"-maxrate", "4000k",
		"-g", "60",
		"-f", "rtsp",
		"-rtsp_transport", "tcp",
		fmt.Sprintf(streamUrl),
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Errorf("Error creating stdin pipe: %v", err)
		return nil, nil, err
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		logger.Errorf("Error starting command: %v", err)
		return nil, nil, err
	}

	logger.Infof("Streaming command started successfully")
	return cmd, stdin, nil
}
