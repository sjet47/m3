package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

var (
	ErrCheckSumNotMatch = errors.New("checksum not match")
)

type DownloadTask struct {
	FileName string
	Url      string
	MD5Sum   string // Set to empty string to skip checksum verification
}

func Download(tasks ...*DownloadTask) (success int64) {
	successCnt := new(atomic.Int64)
	wg := new(sync.WaitGroup)
	wg.Add(len(tasks))
	for _, task := range tasks {
		go func(task *DownloadTask) {
			defer wg.Done()

			skip, err := task.download(".")
			if err != nil {
				log.Printf("Download %s error: %s", task.FileName, err)
				return
			}

			if skip {
				log.Printf("Skip download %s: file already exist", task.FileName)
			}

			successCnt.Add(1)
		}(task)
	}
	wg.Wait()
	return successCnt.Load()
}

// Download downloads file from url and save it to fileName
// returns skip if file already exists and has the same md5sum
func (d *DownloadTask) download(path string) (skip bool, err error) {
	filePath := filepath.Join(path, d.FileName)

	if passHashVerify, _ := verifyMD5Sum(filePath, d.MD5Sum); passHashVerify {
		return true, nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		return false, errors.Wrapf(err, "create file %s error", filePath)
	}
	defer f.Close()

	fmt.Printf("Downloading %s\n", d.FileName)
	rsp, err := http.Get(d.Url)
	if err != nil {
		return false, errors.Wrapf(err, "request %s error", d.Url)
	}
	defer rsp.Body.Close()

	hasher := md5.New()
	_, err = io.Copy(io.MultiWriter(hasher, f), rsp.Body)
	if err != nil {
		return false, errors.Wrapf(err, "write file %s error", filePath)
	}

	if len(d.MD5Sum) > 0 && !strings.EqualFold(hex.EncodeToString(hasher.Sum(nil)), d.MD5Sum) {
		return false, ErrCheckSumNotMatch
	}

	return false, nil
}

func verifyMD5Sum(path, md5sum string) (bool, error) {
	if len(path) == 0 || len(md5sum) == 0 {
		return false, nil
	}

	if !IsFileExist(path) {
		return false, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	hasher := md5.New()
	io.Copy(hasher, f)
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	return strings.EqualFold(fileHash, md5sum), nil
}
