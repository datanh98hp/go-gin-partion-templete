package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}
var allowedMimmTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

const (
	maxFileSize = 5 << 20
)

func ValidateAndSaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	//validate file
	// check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExts[ext] {
		return "", errors.New("unsupport file type")
	}
	// check file size
	if fileHeader.Size > maxFileSize {
		return "", errors.New("file is too large.Max size is 5MB")
	}
	// check content file type
	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("Can not open file")

	}
	defer file.Close()          // close file
	buffer := make([]byte, 512) // doc file voi dang byte
	_, e := file.Read(buffer)
	if e != nil {
		return "", errors.New("Can not read file")
	}
	//dectect file type
	mimmType := http.DetectContentType(buffer)
	mime := http.DetectContentType(buffer)
	if !allowedMimmTypes[mimmType] {
		return "", fmt.Errorf("Invalid file content type %s, only support : %s", mime, keys(allowedMimmTypes))
	}
	//change filename
	fimeName := fmt.Sprintf("%d%s", uuid.New().ID(), ext)
	// create folder if not exist
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		return "", errors.New("Can not create upload folder")
	}
	savePath := filepath.Join(uploadDir, fimeName)
	if err := saveFile(fileHeader, savePath); err != nil {
		return "", err
	} // save file to disksaveFile(fileHeader, savePath)
	return fimeName, nil
}

func saveFile(fileHeader *multipart.FileHeader, des string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err

	}
	defer src.Close()
	// create new file

	out, err := os.Create(des)
	if err != nil {
		return err

	}
	defer out.Close()
	//copy data to file
	_, err = io.Copy(out, src)
	if err != nil {

	}
	return err
}
