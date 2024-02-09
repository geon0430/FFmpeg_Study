package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gocv.io/x/gocv"
)

var (
	width     = 1920
	height    = 1080
	fps       = 30
	videoSrc1 = "rtsp://admin:qazwsx123!@192.168.10.71/0/1080p/media.smp"
)

func main() {
	go readAndProcessStream(videoSrc1, "stream")

	select {}
}

func readAndProcessStream(src string, streamName string) {
	cap, err := gocv.OpenVideoCapture(src)
	if err != nil {
		log.Fatalf("Error opening video capture device: %v\n", src)
		return
	}
	defer cap.Close()

	cap.Set(3, float64(width))
	cap.Set(4, float64(height))
	cap.Set(5, float64(fps))

	cmd := exec.Command("ffmpeg",
		"-re",
		// "-hwaccel", "cuda",
		"-f", "rawvideo",
		"-pixel_format", "bgr24",
		"-video_size", fmt.Sprintf("%dx%d", width, height),
		"-framerate", fmt.Sprint(fps),
		"-i", "-",
		"-c:v", "libx264",
		"-preset", "p4",
		"-f", "rtsp",
		"-rtsp_transport", "tcp",
		fmt.Sprintf("rtsp://localhost:8444/%s", streamName),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error creating stdin pipe: %v", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error starting command: %v", err)
		return
	}

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := cap.Read(&img); !ok {
			fmt.Println("Error reading frame from stream")
			break
		}

		if img.Empty() {
			fmt.Println("Empty frame received")
			continue
		}

		// 바이트 슬라이스로 변환
		imgData := img.ToBytes()

		// FFmpeg에 전달
		_, err = stdin.Write(imgData)
		if err != nil {
			log.Fatalf("Error writing to stdin: %v", err)
			return
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Error waiting for command to finish: %v", err)
	}
}

