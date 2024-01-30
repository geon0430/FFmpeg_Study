func FFmpegRead(
	pipelineConfig util.PipelineInfo,
	logger *logrus.Logger,
	errorChan chan<- error,
	ctx context.Context,
) (*exec.Cmd, io.ReadCloser, error) {

	videoSrc := pipelineConfig.RtspInfo.RTSP
	FPS := pipelineConfig.RtspInfo.FPS
	In_width := pipelineConfig.RtspInfo.IN_WIDTH
	In_height := pipelineConfig.RtspInfo.IN_HEIGHT

	// RTSP exec command
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-hwaccel", "cuda",
		"-i", videoSrc,
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "300",
		"-timeout", "300000",
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
	logger.Debugf("ReadFFmpeg return cmd, stdout successfully")

	return cmd, stdout, nil
}
