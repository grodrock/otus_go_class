package main

import (
	"errors"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy копирует данные из файла fromPath в toPath
//
// offset - смещение, скоторого читаем файл, по умолчанию 0
// при превышении размера файла возвращает ошибку ErrOffsetExceedsFileSize;
// limit - количество байт для считывания из файла, по умолчанию 0 - весь файл
func Copy(fromPath, toPath string, offset, limit int64) error {

	// проверка превышения смещением размера файла
	fInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if offset > fInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	file, err := os.Open(fromPath)
	defer file.Close()
	if err != nil {
		return err
	}

	return nil
}
