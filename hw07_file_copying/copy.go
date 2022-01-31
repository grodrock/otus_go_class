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
// limit - количество байт для считывания из файла, по умолчанию 0 - весь файл.
func Copy(fromPath, toPath string, offset, limit int64) error {
	// проверка превышения смещением размера файла
	fInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fSize := fInfo.Size()
	if offset > fSize {
		return ErrOffsetExceedsFileSize
	}
	if fSize == 0 {
		return ErrUnsupportedFile
	}

	// устанавливаем необходимое кол-во данных для копирования
	if limit == 0 || limit > fSize {
		limit = fSize
	}
	if offset+limit > fSize {
		limit = fSize - offset
	}

	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read data
	data := make([]byte, limit)
	file.Seek(offset, 1)
	file.Read(data)

	// Write data
	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer fileTo.Close()
	_, err = fileTo.Write(data)
	if err != nil {
		return err
	}
	return nil
}
