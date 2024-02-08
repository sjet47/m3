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

func Download(dir string, tasks ...*DownloadTask) (success int64) {
	successCnt := new(atomic.Int64)
	wg := new(sync.WaitGroup)
	wg.Add(len(tasks))
	proc := mpb.New(mpb.WithWaitGroup(wg))

	for _, task := range tasks {
		go func(task *DownloadTask) {
			defer wg.Done()
			if err := task.download(dir, proc); err != nil {
				log.Printf("Download %s error: %s", task.FileName, err)
				return
			}
			successCnt.Add(1)
		}(task)
	}
	proc.Wait()
	return int64(successCnt.Load())
}

// Download downloads file from url and save it to fileName
// returns skip if file already exists and has the same md5sum
func (d *DownloadTask) download(path string, proc *mpb.Progress) error {
	filePath := filepath.Join(path, d.FileName)

	checkLocalBar := proc.AddSpinner(0,
		mpb.SpinnerOnLeft,
		mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.Name("🔍", decor.WCSyncSpaceR),
			decor.Name(d.FileName, decor.WCSyncSpaceR),
			decor.Name("Checking local file...", decor.WCSyncSpaceR),
		),
	)
	passHashVerify := VerifyMD5Sum(filePath, d.MD5Sum, checkLocalBar)
	checkLocalBar.SetTotal(0, true)

	downloadBar := proc.AddBar(0,
		mpb.BarClearOnComplete(),
		mpb.BarParkTo(checkLocalBar),
		mpb.OptionOnCondition(mpb.PrependDecorators(
			decor.OnComplete(decor.Name("🔥", decor.WCSyncSpaceR), "✅"),
			decor.Name(d.FileName, decor.WCSyncSpaceR),
			decor.OnComplete(decor.CountersKibiByte("% 6.1f / % 6.1f ", decor.WCSyncSpace), ""),
			decor.OnComplete(decor.AverageSpeed(decor.UnitKiB, "% .2f", decor.WCSyncSpace), ""),
			decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_MMSS, 60, decor.WCSyncSpace), ""),
			decor.OnComplete(decor.Name("", decor.WCSyncSpaceR), ""),
		), func() bool {
			return !passHashVerify
		}),
		mpb.OptionOnCondition(mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 6}), ""),
		), func() bool {
			return !passHashVerify
		}),

		mpb.OptionOnCondition(mpb.PrependDecorators(
			decor.Name("✅", decor.WCSyncSpaceR),
			decor.Name(d.FileName, decor.WCSyncSpaceR),
		), func() bool {
			return passHashVerify
		}),
	)

	if passHashVerify {
		downloadBar.SetTotal(0, true)
		return nil
	}

	defer func() {
		if !downloadBar.Completed() {
			proc.Abort(downloadBar, true)
		}
	}()

	f, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "create file %s error", filePath)
	}
	defer f.Close()

	rsp, err := http.Get(d.Url)
	if err != nil {
		return errors.Wrapf(err, "request %s error", d.Url)
	}
	defer rsp.Body.Close()

	hasher := md5.New()
	var r io.ReadCloser = rsp.Body
	if downloadBar != nil {
		downloadBar.SetTotal(rsp.ContentLength, false)
		r = downloadBar.ProxyReader(r)
	}

	_, err = io.Copy(io.MultiWriter(f, hasher), r)
	if err != nil {
		return errors.Wrapf(err, "write file %s error", filePath)
	}

	if len(d.MD5Sum) > 0 && !strings.EqualFold(hex.EncodeToString(hasher.Sum(nil)), d.MD5Sum) {
		return ErrCheckSumNotMatch
	}

	return nil
}

func VerifyMD5Sum(path, md5sum string, bar *mpb.Bar) bool {
	if len(path) == 0 || len(md5sum) == 0 {
		return false
	}

	if !IsFileExist(path) {
		return false
	}

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	var r io.Reader = f
	if bar != nil {
		if stat, err := f.Stat(); err == nil {
			bar.SetTotal(stat.Size(), false)
			r = bar.ProxyReader(r)
		}
	}

	hasher := md5.New()
	io.Copy(hasher, r)
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	return strings.EqualFold(fileHash, md5sum)
}
