package mod

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

func downloadFiles(fileUrlMap map[string]string) int64 {
	successCnt := new(atomic.Int64)
	wg := new(sync.WaitGroup)
	wg.Add(len(fileUrlMap))
	for fileName, url := range fileUrlMap {
		go func(fileName, url string) {
			defer wg.Done()
			if err := download(fileName, url); err != nil {
				log.Printf("Download %s error: %s", url, err)
				return
			}
			successCnt.Add(1)
		}(fileName, url)
	}
	wg.Wait()
	return successCnt.Load()
}

func download(fileName, url string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return errors.Wrapf(err, "create file %s error", fileName)
	}
	defer f.Close()

	fmt.Printf("Downloading %s\n", fileName)
	rsp, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "get %s error", url)
	}
	defer rsp.Body.Close()
	_, err = io.Copy(f, rsp.Body)
	return err
}
