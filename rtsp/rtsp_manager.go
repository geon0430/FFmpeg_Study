package rtsp

import (
	"context"
	"fmt"
	util "go_vms/src/pipeline/util"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

func ReadRtsp(
	ctx context.Context,
	pipelineConfig util.PipelineInfo,
	rtspReadBuf chan<- *[]byte,
	errorChan chan<- error,
	logger *logrus.Logger,
) {
	channels := pipelineConfig.RtspInfo.Channels
	In_height := pipelineConfig.RtspInfo.IN_HEIGHT
	In_width := pipelineConfig.RtspInfo.IN_WIDTH
	frameSize := channels * In_height * In_width

	ReadUrl := pipelineConfig.RtspInfo.RTSP

	ControlDisplayOn := false

	// 화면 조정 이미지
	controlFrame, err := LoadImage(
		pipelineConfig.RtspInfo.ControlImage,
		In_width, In_height,
		"Cannot connect CAM",
		logger,
		errorChan,
	) //byte
	if err != nil {
		logger.Error("RTSP Read Control Image Load fail")
		return
	}

	waitTime := int(3)
	controlFrameDuration := time.Duration(waitTime) * time.Second
	controlTicker := util.CreateTicker(controlFrameDuration) // 하드코딩됨 TODO
	defer controlTicker.Ticker.Stop()

	FPS := pipelineConfig.RtspInfo.FPS
	//frameDelayDuration := time.Duration(1 / FPS * 2.0 * float64(time.Second))
	frameDuration := time.Duration(1 / FPS * 2.0 * float64(time.Second))
	frameTicker := util.CreateTicker(frameDuration)
	defer frameTicker.Ticker.Stop()

	logger.Infof(
		"Read RTSP Wait long : %fs, Wait short : %fs",
		float32(controlTicker.Period)/1e9, float32(frameTicker.Period)/1e9,
	)

	cmd, stdout, errFFmpegPrc := FFmpegRead(pipelineConfig, logger, errorChan, ctx)
	defer stdout.Close()
	defer cmd.Process.Kill()
	time.Sleep(500 * time.Millisecond)

	// 시간 측정을 위한 중간 channel
	frameChannel := make(chan *[]byte, cap(rtspReadBuf))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// FFmpegPrc를 에러 체크를 통해서 작동 유무 판단
				// 만약 동작하지 않으면 다시 실행
				if errFFmpegPrc != nil {
					logger.Errorf("Url: %s | FFmpeg Read processor is not running", ReadUrl)
					errorChan <- fmt.Errorf("Url: %s | FFmpeg Read processor is not running", ReadUrl)
					// Kill Previous FFmpegPrc
					stdout.Close()
					cmd.Process.Kill()
					// 재시작 하기 전에 시간 차이를 줌
					time.Sleep(1 * time.Second)
					// Restart FFmpegPrc
					cmd, stdout, errFFmpegPrc = FFmpegRead(pipelineConfig, logger, errorChan, ctx)
					defer stdout.Close()
					defer cmd.Process.Kill()
					logger.Errorf("Url: %s | FFmpeg Read processor Restart", ReadUrl)
					errorChan <- fmt.Errorf("Url: %s | FFmpeg Read processor Restart", ReadUrl)
					continue
				}

				// io.ReadFull은 데이터가 들어올 때 까지 대기함
				// 데이터 지연이 이루어 질 수 없음
				// async 작업을 통해서 계속 Wait해도 만들고 다른 곳에서 시간 체크를 해야 함
				frame := make([]byte, frameSize)
				_, errStdout := io.ReadFull(stdout, frame)
				if errStdout == nil {
					frameChannel <- &frame
					logger.Debugf("Url: %s | Read FFmpeg is success", ReadUrl)
				} else if errStdout == io.ErrShortBuffer {
					logger.Errorf("Url: %s | Failed to read data from stdout(ReadRtsp), Check frame is short | %v", ReadUrl, err)
					errorChan <- fmt.Errorf("Url: %s | Failed to read data from stdout(ReadRtsp), Check frame is short | %v", ReadUrl, err)
				} else {
					logger.Errorf("Url: %s | Failed to read data from stdout(ReadRtsp) | %v", ReadUrl, err)
					errorChan <- errStdout
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case tmpFrame := <-frameChannel:
			// Frame이 들어오면 Ticker를 reset해서 에러 메시지가 출력되지 않게 함
			// 지정 시간보다 Frame이 늦게 들어오면 아래 Ticker case가 발생
			logger.Debugf("Url: %s | Insert Frame at RtspReadBuffer", ReadUrl)
			rtspReadBuf <- tmpFrame
			if ControlDisplayOn {
				ControlDisplayOn = false
				logger.Infof("Url: %s | Stream Restarted", ReadUrl)
				logger.Infof("Url: %s | ControlDisplay %t", ReadUrl, ControlDisplayOn)
			}
			controlTicker.ResetTicker()
			frameTicker.ResetTicker()

		case <-controlTicker.Ticker.C:
			ControlDisplayOn = true
			logger.Warnf("Url: %s | %d second, timeout occurred while reading data from FFmpeg", ReadUrl, waitTime)
			logger.Warnf("Url: %s | ControlDisplay %t", ReadUrl, ControlDisplayOn)
			errorChan <- fmt.Errorf("%d second, timeout occurred while reading data from FFmpeg", waitTime)
		case <-frameTicker.Ticker.C:
			logger.Debugf("Url: %s | Frame Read is delayed ", ReadUrl)
			if ControlDisplayOn {
				rtspReadBuf <- &controlFrame
				logger.Debugf("Url: %s | %d second contorl display streaming", ReadUrl, frameDuration/1e9)
			} else {
				logger.Warnf("Url: %s | ControlDisplay %t", ReadUrl, ControlDisplayOn)
			}
		}
	}

}

// Stream 부분에도 ControlDisplay가 필요할거 같음
// CAM에서 잘못된 것인지, 중간 처리가 잘못된 것인지 알기 위해서
// ControlDisplay에 Text를 넣을 수 있게 만들어야 할거 같음
func StreamRtsp(
	ctx context.Context,
	pipelineConfig util.PipelineInfo,
	videoWidth, videoHeight int,
	streamUrl string,
	postprocessBuf <-chan *[]byte,
	errorChan chan<- error,
	logger *logrus.Logger,
) {

	logger.Infof("StreamRtsp function Start | %s", streamUrl)
	logger.Infof("Url: %s , W: %d, H: %d", streamUrl, videoWidth, videoHeight)

	// 화면조정 이미지
	ControlDisplayOn := false
	controlFrame, err := LoadImage(

		pipelineConfig.RtspInfo.ControlImage,
		videoWidth, videoHeight,
		"Processing Error",
		logger,
		errorChan,
	) //byte
	if err != nil {
		logger.Error("RTSP Stream Control Image Load fail")
		return
	}

	waitTime := int(4)
	controlFrameDuration := time.Duration(waitTime) * time.Second
	controlTicker := util.CreateTicker(controlFrameDuration) // 하드코딩됨 TODO
	defer controlTicker.Ticker.Stop()

	FPS := pipelineConfig.RtspInfo.FPS
	frameDuration := time.Duration(1 / FPS * 2.0 * float64(time.Second))
	frameTicker := util.CreateTicker(frameDuration)
	defer frameTicker.Ticker.Stop()

	logger.Infof(
		"Stream RTSP Wait long : %fs, Wait short : %fs",
		float32(controlTicker.Period)/1e9, float32(frameTicker.Period)/1e9,
	)
	cmd, stdin, errFFmpegPrc := FFmpegStream(
		pipelineConfig,
		videoWidth, videoHeight,
		streamUrl,
		logger,
		errorChan,
	)
	defer stdin.Close()
	defer cmd.Process.Kill()
	// Start delay time
	time.Sleep(1 * time.Second)

	// channel for write ffmpeg streaming
	streamWriteChan := make(chan *[]byte, cap(postprocessBuf))

	writerErrorCounter := 0
	for {
		// Err 확인을 통해서 FFmpegPrc 작동 유무 확인
		// 작동하지 않으면 다시 실행
		if errFFmpegPrc != nil {
			errorString := "FFmpeg stream is not running"
			logger.Errorf("Url: %s | %s", streamUrl, errorString)
			errorChan <- fmt.Errorf("Url: %s | %s", streamUrl, errorString)
			// Kill Previous FFmpegPrc
			stdin.Close()
			cmd.Process.Kill()
			// Restart delay time
			time.Sleep(1 * time.Second)
			// Restart FFmpegPrc
			cmd, stdin, errFFmpegPrc = FFmpegStream(
				pipelineConfig,
				videoWidth, videoHeight,
				streamUrl,
				logger,
				errorChan,
			)
			defer stdin.Close()
			defer cmd.Process.Kill()
			logger.Errorf("Url: %s | FFmpeg Stream restart", streamUrl)
			errorChan <- fmt.Errorf("Url: %s | FFmpeg Stream processor Restart", streamUrl)
			continue
		}

		select {
		case <-ctx.Done():
			return
		case frameFromBuf := <-postprocessBuf:
			streamWriteChan <- frameFromBuf
			controlTicker.ResetTicker()
			frameTicker.ResetTicker()
			ControlDisplayOn = false
			logger.Debugf("Url: %s | Get frame from preprocessBuf", streamUrl)
		case frame := <-streamWriteChan:
			_, errStdin := stdin.Write(*frame)
			if errStdin != nil {
				errorChan <- errStdin
				logger.Debugf("Url: %s | Stream write is error : %v", streamUrl, errStdin)
				if writerErrorCounter > int(FPS) {
					errFFmpegPrc = errStdin
					writerErrorCounter = 0
				} else {
					writerErrorCounter += 1
				}
				continue
			}
			logger.Debugf("Url: %s | Stream write is success", streamUrl)
		case <-controlTicker.Ticker.C:
			ControlDisplayOn = true
			logger.Warnf("Url: %s | %d second, timeout occurred while reading data from FFmpeg", streamUrl, waitTime)
			logger.Warnf("Url: %s | Stream ControlDisplay %t", streamUrl, ControlDisplayOn)
			errorChan <- fmt.Errorf("%d second, timeout occurred while reading data from FFmpeg", waitTime)
		case <-frameTicker.Ticker.C:
			logger.Debugf("Url: %s | Frame Read is delayed ", streamUrl)
			if ControlDisplayOn {
				streamWriteChan <- &controlFrame
				logger.Debugf("Url: %s | Stream Write is delayed", streamUrl)
			} else {
				logger.Warnf("Url: %s | Stream ControlDisplay %t", streamUrl, ControlDisplayOn)
			}
		}
	}

	///////////
}

