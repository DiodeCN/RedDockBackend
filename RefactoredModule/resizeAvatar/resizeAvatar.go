package resizeavatar

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

func compressImage(imageBytes []byte, maxSize int) ([]byte, error) {
	// 解码图像
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Println("failed to decode image:", err)
		return nil, err
	}

	// 缩放图像
	img = resize.Resize(0, uint(maxSize), img, resize.Lanczos3)

	// 使用JPEG编码压缩图像
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		log.Println("failed to encode image:", err)
		return nil, err
	}

	// 检查图像大小是否超过最大限制
	if buf.Len() > maxSize {
		// 如果超过最大限制，则使用imaging库继续压缩图像
		img, err = imaging.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			log.Println("failed to decode image:", err)
			return nil, err
		}

		img = imaging.Fit(img, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)

		buf.Reset()
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			log.Println("failed to encode image:", err)
			return nil, err
		}
	}

	// 返回压缩后的图像数据
	return buf.Bytes(), nil
}
