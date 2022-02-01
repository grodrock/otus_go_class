package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
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
	toCopyBytes := limit
	// если лимит не задан или превышает размер файла, копируем весь файл
	if limit == 0 || limit > fSize {
		toCopyBytes = fSize
	}
	// если лимит после смещения превышает размер файла, копируем только остаток файла
	if offset+limit > fSize {
		toCopyBytes = fSize - offset
	}

	file, err := os.Open(fromPath) // reader
	if err != nil {
		return err
	}
	defer file.Close()

	fileTo, err := os.Create(toPath) // writer
	if err != nil {
		return err
	}
	defer fileTo.Close()

	file.Seek(offset, 1) // move to offset
	buffSize := 1024
	if buffSize > int(toCopyBytes) {
		buffSize = int(toCopyBytes)
	}
	dataBuffer := make([]byte, buffSize) // копируем с буфером
	bar := pb.New(int(toCopyBytes)).SetUnits(pb.U_BYTES)
	bar.Start()
	for toCopyBytes > 0 {
		// Read data
		readed, err := file.Read(dataBuffer)
		if err != nil && err != io.EOF {
			return err
		}

		if readed < len(dataBuffer) {
			dataBuffer = dataBuffer[:readed]
		}
		// Write data
		_, err = fileTo.Write(dataBuffer)
		if err != nil {
			return err
		}
		toCopyBytes -= int64(readed)
		bar.Add(readed)
	}
	bar.Finish()

	return nil
}
