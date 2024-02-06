package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var (
	ErrCheckSumNotMatch = errors.New("checksum not match")
)

type DownloadTask struct {
	FileName string
	FileSize int64
	Url      string
	MD5Sum   string // Set to empty string to skip checksum verification
}

func Download(tasks ...*DownloadTask) (success int64) {
	successCnt := new(atomic.Int64)
	wg := new(sync.WaitGroup)
	wg.Add(len(tasks))
	proc := mpb.New(mpb.WithWaitGroup(wg))

	for _, task := range tasks {
		go func(task *DownloadTask) {
			defer wg.Done()
			bar := proc.AddBar(0,
				mpb.BarClearOnComplete(),
				mpb.TrimSpace(),
				mpb.PrependDecorators(
					decor.OnComplete(decor.Name("🔥", decor.WCSyncSpaceR), "✅"),
					decor.Name(task.FileName, decor.WCSyncSpaceR),
					decor.OnComplete(decor.CountersKibiByte("% 6.1f / % 6.1f ", decor.WCSyncSpace), ""),
					decor.OnComplete(decor.AverageSpeed(decor.UnitKiB, "% .2f", decor.WCSyncSpace), ""),
					decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_MMSS, 60, decor.WCSyncSpace), ""),
					decor.OnComplete(decor.Name("", decor.WCSyncSpaceR), ""),
				),
				mpb.AppendDecorators(
					decor.OnComplete(decor.Percentage(decor.WC{W: 6}), ""),
				),
			)

			skip, err := task.download(".", bar)
			if err != nil {
				log.Printf("Download %s error: %s", task.FileName, err)
				proc.Abort(bar, false)
				return
			}

			if skip {
				log.Printf("Skip download %s: file already exist", task.FileName)
				proc.Abort(bar, true)
			}

			successCnt.Add(1)
		}(task)
	}
	proc.Wait()
	return int64(successCnt.Load())
}

// Download downloads file from url and save it to fileName
// returns skip if file already exists and has the same md5sum
func (d *DownloadTask) download(path string, bar *mpb.Bar) (skip bool, err error) {
	filePath := filepath.Join(path, d.FileName)

	if passHashVerify, _ := verifyMD5Sum(filePath, d.MD5Sum); passHashVerify {
		bar.SetTotal(0, true)
		return true, nil
	}

	f, err := os.Create(filePath)
	if err != nil {
		return false, errors.Wrapf(err, "create file %s error", filePath)
	}
	defer f.Close()

	rsp, err := http.Get(d.Url)
	if err != nil {
		return false, errors.Wrapf(err, "request %s error", d.Url)
	}
	defer rsp.Body.Close()

	hasher := md5.New()
	var r io.Reader = rsp.Body
	if bar != nil {
		bar.SetTotal(rsp.ContentLength, false)
		r = bar.ProxyReader(r)
	}

	_, err = io.Copy(io.MultiWriter(f, hasher), r)
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
