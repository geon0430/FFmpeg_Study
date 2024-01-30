package main

import (
	"context"
	"fmt"
	"image"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

var (
	cmd             *exec.Cmd
	running         bool
	mu              sync.Mutex
	readSuccessChan = make(chan bool)
	dataChan        = make(chan []byte, 1)
	errChan         = make(chan error, 1)
)

func ReadFFmpeg(videoSrc string, In_width, In_height int, logger *logrus.Logger, errorChan chan<- error) (*exec.Cmd, io.ReadCloser) {
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-hwaccel", "cuda",
		"-i", videoSrc,
		"-reconnect", "1",
		"-reconnect_at_eof", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "15",
		"-preset", "fast",
		"-framerate", "30",
		"-vf", fmt.Sprintf("scale=%d:%d", In_width, In_height),
		"-f", "rawvideo",
		"-pix_fmt", "bgr24",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorf("Failed to create stdout pipe: %v", err)
		errorChan <- err
		return nil, nil
	}

	if err := cmd.Start(); err != nil {
		logger.Errorf("Failed to start ffmpeg: %v", err)
		errorChan <- err
		return nil, nil
	}
	logger.Debugf("ReadFFmpeg return cmd, stdout successfully")
	return cmd, stdout
}

func ReadRTSP(In_width, In_height int, ctx context.Context, videoSrc string, logger *logrus.Logger, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Info("Starting RTSP stream processing")
			cmd, stdout := ReadFFmpeg(videoSrc, In_width, In_height, logger, errorChan)
			if cmd == nil || stdout == nil {
				readSuccessChan <- false
				return
			}

			buf := make([]byte, In_width*In_height*3)
			for {
				timeout := time.NewTimer(2 * time.Second)
				readChan := make(chan error, 1)

				go func() {
					_, err := io.ReadFull(stdout, buf)
					readChan <- err
				}()

				select {
				case err := <-readChan:
					if err != nil {
						logger.Errorf("Failed to read data from FFmpeg: %v", err)
						readSuccessChan <- false
						continue
					}
					dataChan <- buf
					readSuccessChan <- true
					timeout.Stop()
				case <-timeout.C:
					logger.Error("Timeout occurred while reading data from FFmpeg")
					readSuccessChan <- false
					continue
				}
			}
		}
	}
}

func processRTSPStream(In_width, In_height int, ctx context.Context, logger *logrus.Logger, errorChan chan<- error) {
	window := gocv.NewWindow("RTSP Stream")
	window.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowNormal)
	defer window.Close()

	resizedWidth := 1280
	resizedHeight := 720

	blackImage := gocv.NewMatWithSize(resizedHeight, resizedWidth, gocv.MatTypeCV8UC3)
	defer blackImage.Close()
	blackImage.SetTo(gocv.NewScalar(0, 0, 0, 0))

	for {
		select {
		case success := <-readSuccessChan:
			logger.Infof("Read success channel: %t", success)
			if !success {
				window.IMShow(blackImage)
				if window.WaitKey(1) >= 0 {
					return
				}
				continue
			}
		case buf := <-dataChan:

			img, err := gocv.NewMatFromBytes(In_height, In_width, gocv.MatTypeCV8UC3, buf)
			if err != nil {
				logger.Error("Failed to convert bytes to Mat: %v", err)
				errorChan <- err
				continue
			}

			resizedImg := gocv.NewMat()
			defer resizedImg.Close()
			gocv.Resize(img, &resizedImg, image.Pt(resizedWidth, resizedHeight), 0, 0, gocv.InterpolationDefault)

			window.IMShow(resizedImg)
			if window.WaitKey(1) >= 0 {
				return
			}
			img.Close()

		}
	}
}

func main() {
	logger := logrus.New()
	errorChan := make(chan error)
	ctx := context.Background()

	videoSrc := "rtsp://admin:qazwsx123!@192.168.10.70/0/720p/media.smp"
	In_width := 1280
	In_height := 720

	go ReadRTSP(In_width, In_height, ctx, videoSrc, logger, errorChan)
	go processRTSPStream(In_width, In_height, ctx, logger, errorChan)

	for {
		select {
		case err := <-errorChan:
			logger.Errorf("Error received: %v", err)
		}
	}
}
