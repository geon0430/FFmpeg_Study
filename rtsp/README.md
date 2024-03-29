# FFmpeg stream option
----

### pix_fmt Option
- yuv420p: YUV 색공간의 4:2:0 크로마 서브샘플링을 사용합니다. 이는 스탠다드 디피니션(SD), 하이 디피니션(HD), 그리고 풀 하이 디피니션(FHD) 비디오에 널리 사용되며, 대부분의 비디오 코덱과 호환됩니다.
- rgb24: RGB 색공간을 사용하며, 각 색상 채널에 8비트를 할당합니다. 총 24비트의 색상 깊이를 제공합니다.
- yuv444p: YUV 색공간의 4:4:4 크로마 서브샘플링을 사용합니다. 이 포맷은 색상 정보를 보다 정확하게 보존하지만, yuv420p보다 더 큰 파일 크기를 가집니다.
- yuv422p: YUV 색공간의 4:2:2 크로마 서브샘플링을 사용합니다. 이는 yuv420p와 yuv444p의 중간 정도의 색상 정보를 보존하며, 전문 비디오 편집 및 방송에 자주 사용됩니다.

### f (format)
- rawvideo: 비압축 원시 비디오 데이터를 의미합니다. 이 포맷을 사용할 때는 추가적으로 픽셀 포맷(-pix_fmt), 해상도 등의 비디오 스트림 정보를 명시해야 할 수 있습니다.
- image2: 일련의 이미지 파일로부터 비디오를 생성하거나, 비디오로부터 이미지를 추출할 때 사용합니다. 예를 들어, JPEG 또는 PNG 파일 시퀀스를 처리할 수 있습니다.
- mp4, avi, mkv 등: 특정 컨테이너 포맷으로 미디어 파일을 읽거나 쓸 때 사용됩니다. 이러한 포맷은 비디오와 오디오 스트림, 자막 등 다양한 미디어 데이터를 하나의 파일에 포함할 수 있습니다.
- flv: 플래시 비디오 포맷으로, 주로 스트리밍 서비스에 사용됩니다.
- hls: HTTP Live Streaming을 위한 포맷입니다. 비디오 스트림을 여러 작은 파일로 분할하고, 이를 순차적으로 전송하여 실시간 스트리밍을 가능하게 합니다.
- rtsp, rtmp: 네트워크 스트리밍 프로토콜을 위한 포맷입니다. 실시간 데이터 전송에 주로 사용됩니다.
- mjpeg: 모션 JPEG 비디오 스트림을 처리할 때 사용합니다. 각 프레임이 독립적인 JPEG 이미지로 인코딩되어 있습니다.
- webp, gif: 이미지 파일 포맷으로, 특히 gif는 애니메이션 이미지 생성에 사용됩니다.
- srt, ass: 자막 파일 포맷입니다. 비디오 파일과 함께 자막을 처리할 때 사용됩니다.

### preset
- ultrafast: 가장 빠른 인코딩 속도를 제공합니다. 실시간 스트리밍이나 빠른 인코딩이 필요한 경우에 적합합니다.
- superfast
- veryfast
- faster
- fast
- medium: 기본값입니다. 속도와 품질 사이에 균형을 맞춥니다.
- slow
- slower
- veryslow: 가장 느린 인코딩 속도를 제공합니다. 최대한의 압축률과 품질을 얻을 수 있으며, 파일 크기를 줄이는 데 유리합니다.
