package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	ErrMoveFile          = errors.New("ファイルを移動できません")
	ErrOpenFile          = errors.New("ファイルをオープンできません")
	ErrCreateFile        = errors.New("ファイルを作成できません")
	ErrFileAlreadyExists = errors.New("ファイルが既に存在します")
)

// Move File は srcPath を dstPath に移動する。
// ドライブが異なっていても移動可能。
func MoveFile(srcPath, dstPath string) error {
	if len(srcPath) == 0 {
		return fmt.Errorf("src path is empty")
	}
	if len(dstPath) == 0 {
		return fmt.Errorf("dst path is empty")
	}

	inFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("%w: %w: %s", ErrOpenFile, err, srcPath)
	}
	defer inFile.Close()

	if _, err := os.Stat(dstPath); err == nil {
		return fmt.Errorf("%w", ErrFileAlreadyExists)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("%s: %w", dstPath, err)
	}

	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("%w: %w: %s", ErrCreateFile, err, dstPath)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return fmt.Errorf("failed to copy to dst from src: %w", err)
	}

	inFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(srcPath)
	if err != nil {
		return fmt.Errorf("failed to remove src file: %w", err)
	}

	return nil
}

func MoveFileToDir(srcPath, dstDirPath string) error {
	name := filepath.Base(srcPath)
	dstPath := filepath.Join(dstDirPath, name)
	if err := MoveFile(srcPath, dstPath); err != nil {
		return err
	}
	return nil
}
