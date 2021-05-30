package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mostafasolati/leviathan/contracts"

	"github.com/disintegration/imaging"
)

// ImageURL returns absolute url of the image
func ImageURL(baseUrl, image string) string {
	if image == "" {
		return ""
	}
	return fmt.Sprintf("%s/v1/image/%s", baseUrl, image)
}

// BannerURL returns the fully-qualified URL for a banner.
func BannerURL(baseURL, image string) string {
	if image == "" {
		return ""
	}
	return fmt.Sprintf("%s/v1/image/banner/%s", baseURL, image)
}

func UploadFile(imagesRoot string, h *multipart.FileHeader) (string, error) {
	src, err := h.Open()
	if err != nil {
		return "", err
	}

	var destFilename string
	var filename string

	ext, err := ParseImageExtFromMIMEType(h.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	path := filepath.Join(imagesRoot, "original")
	_ = os.MkdirAll(path, 0755)
	filename = fmt.Sprintf("%d.%s", time.Now().Unix(), ext)
	destFilename = fmt.Sprintf("%s/%s", path, filename)

	dest, err := os.Create(destFilename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(dest, src)

	return filename, err

}

func Thumb(imagesRoot, file string, w, h int) (string, error) {
	if w > 1000 || h > 1000 {
		return "", contracts.ErrFileDimensionQuota
	}

	thumb := filepath.Join(imagesRoot, "thumbs")
	_ = os.Mkdir(thumb, 0755)
	thumb += fmt.Sprintf("/%d_%d_%s", w, h, file)

	_, err := os.Stat(thumb)
	if os.IsNotExist(err) {
		imagePath := filepath.Join(imagesRoot, "original/"+file)
		image, err := imaging.Open(imagePath)
		if err != nil {
			return "", fmt.Errorf("error in opening %s file: %v", file, err)
		}
		image = imaging.Fill(image, w, h, imaging.Center, imaging.Lanczos)
		err = imaging.Save(image, thumb)
		if err != nil {
			return "", fmt.Errorf("error in saving %s file %v", file, err)
		}
	}

	return thumb, nil
}

// ParseImageExtFromMIMEType determines the image file extension based on its
// MIME type.
func ParseImageExtFromMIMEType(mimeType string) (string, error) {
	mime := strings.Split(mimeType, "/")
	//if mime[0] != "image" {
	//	return "", errors.New("mime type is not an image")
	//}
	return mime[1], nil
}
