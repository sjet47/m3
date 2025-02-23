package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/sjet47/go-curseforge/schema"
	"github.com/sjet47/go-curseforge/schema/enum"
	"github.com/sjet47/m3/internal/index"
	"github.com/sjet47/m3/internal/util"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func Update(confirm bool) error {
	mu := new(sync.Mutex)
	update := make(index.ModIndexes)
	downloadCnt := new(atomic.Int64)

	proc := mpb.New()
	spn := proc.AddSpinner(int64(len(index.Mods)),
		mpb.SpinnerOnLeft,
		mpb.BarClearOnComplete(),
		mpb.PrependDecorators(
			decor.OnComplete(decor.Name("🔥", decor.WCSyncSpaceR), "✅"),
			decor.OnComplete(decor.Name("Fetching updates", decor.WCSyncSpaceR), "Fetch updates"),
			decor.OnComplete(decor.CountersNoUnit("%d/%d", decor.WCSyncSpace), ""),
		),
	)
	for modID, mod := range index.Mods {
		go func(modID schema.ModID, mod *index.Mod) {
			defer spn.Increment()
			modLoader, _ := enum.ParseModLoader(mod.ModLoader)
			file, err := getLatestModFile(modID, modLoader)
			if err == nil {
				defer mod.Update(file)
			}

			if !needUpdate(mod, file) {
				return
			}

			mu.Lock()
			if file != nil {
				downloadCnt.Add(1)
				os.Remove(filePath(mod.File.Name))
			}
			update[modID] = mod
			mu.Unlock()
		}(modID, mod)
	}
	proc.Wait()

	if downloadCnt.Load() == 0 {
		fmt.Println("🎉 All mods are up to date")
		return nil
	}

	fmt.Println(update.String())

	if downloadCnt.Load() > 0 && (confirm || util.Prompt("Update mods?")) {
		downloadMods := make([]*util.DownloadTask, 0, len(update))

		for _, mod := range update {
			// Add to download list
			if len(mod.File.DownloadUrl) > 0 && len(mod.File.Name) > 0 {
				downloadMods = append(downloadMods, &util.DownloadTask{
					FileName: mod.File.Name,
					Url:      mod.File.DownloadUrl,
					MD5Sum:   mod.File.HashMD5,
				})
			}
		}

		downloadCnt := util.Download(".", downloadMods...)
		fmt.Printf("%d mod(s) updated\n", downloadCnt)
	}
	return nil
}

func needUpdate(old *index.Mod, fetched *schema.File) bool {
	if fetched == nil {
		return true
	}
	if fetched.FileDate.After(old.File.Date) && old.File.Name != fetched.FileName {
		return true
	}

	var md5Hash string
	for _, hash := range fetched.Hashes {
		if hash.Algo == enum.HashAlgoMD5 {
			md5Hash = hash.Value
			break
		}
	}
	return !util.VerifyMD5Sum(filePath(old.File.Name), md5Hash, nil)
}

func filePath(fileName string) string {
	return filepath.Join(".", fileName)
}
