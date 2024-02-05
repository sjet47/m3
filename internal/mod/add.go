package mod

import (
	"fmt"
	"sync"

	"github.com/ASjet/go-curseforge"
	"github.com/ASjet/go-curseforge/api"
	"github.com/ASjet/go-curseforge/schema"
	"github.com/ASjet/go-curseforge/schema/enum"
	"github.com/ASjet/m3/internal/index"
	"github.com/ASjet/m3/internal/util"
	"github.com/pkg/errors"
)

func Init(apiKey string) {
	curseforge.InitDefault(apiKey)
}

func Add(modLoaderStr string, optDep bool, ids ...int) error {
	modLoader, err := enum.ParseModLoader(modLoaderStr)
	if err != nil {
		return errors.Wrapf(err, "invalid mod loader %q", modLoaderStr)
	}

	// Fetch direct mod files
	modIDs := util.Map(func(id int) schema.ModID { return schema.ModID(id) }, ids...)
	modFileMap := fetchModFiles(modLoader, modIDs...)

	// Fetch dependencies files
	depIDs := util.Filter(func(id schema.ModID) bool {
		_, ok := modFileMap[id]
		return !ok
	}, extractDeps(optDep, modFileMap)...)
	depFileMap := fetchModFiles(modLoader, depIDs...)

	// Fetch mod info
	modMap := fetchMods(append(util.Keys(modFileMap), util.Keys(depFileMap)...)...)

	// Display mod info in table
	fmt.Println(renderModInfoTable(modMap, modFileMap, depFileMap))

	// Prompt user for download confirmation with mod info
	if promptDownload() {
		downloadMap := make(map[string]string, len(modFileMap))

		// Write to index
		for modID, result := range modFileMap {
			if file := result.Value; file != nil && file.DownloadURL != "" && file.FileName != "" {
				downloadMap[file.FileName] = file.DownloadURL
			}
			index.Mods[modID] = index.NewMod(modLoader, modMap[modID].Value, result.Value)
		}

		// TODO: add checksum verification
		downloadCnt := downloadFiles(downloadMap)
		fmt.Printf("(%d/%d) mod downloaded\n", downloadCnt, len(downloadMap))
	}
	return nil
}

type fetchModResult map[schema.ModID]util.Result[*schema.Mod]

func fetchMods(modIDs ...schema.ModID) fetchModResult {
	wg, mu := new(sync.WaitGroup), new(sync.Mutex)
	result := make(fetchModResult, len(modIDs))

	wg.Add(len(modIDs))
	for _, id := range modIDs {
		go func(modID schema.ModID) {
			defer wg.Done()
			resp, err := api.Mod(modID)

			var res util.Result[*schema.Mod]
			if err != nil {
				res = util.Err[*schema.Mod](err)
				res.Value = &schema.Mod{ID: modID}
			} else {
				res = util.Ok(&resp.Data)
			}

			mu.Lock()
			result[modID] = res
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	return result
}

type fetchFileResult map[schema.ModID]util.Result[*schema.File]

func fetchModFiles(modLoader enum.ModLoader, modIDs ...schema.ModID) fetchFileResult {
	wg, mu := new(sync.WaitGroup), new(sync.Mutex)
	result := make(fetchFileResult, len(modIDs))

	wg.Add(len(modIDs))
	for _, id := range modIDs {
		go func(modID schema.ModID) {
			defer wg.Done()
			resp, err := api.ModFiles(modID,
				api.ModFiles.WithGameVersion(index.Meta.GameVersion),
				api.ModFiles.WithModLoader(modLoader),
				api.ModFiles.WithIndex(0),
				api.ModFiles.WithPageSize(1),
			)
			if err == nil && len(resp.Data) == 0 {
				err = errors.Errorf("mod %d has no files for game version %s and mod loader %s",
					modID, index.Meta.GameVersion, modLoader)
			}

			var res util.Result[*schema.File]
			if err != nil {
				res = util.Err[*schema.File](err)
			} else {
				res = util.Ok(&resp.Data[0])
			}

			mu.Lock()
			result[modID] = res
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	return result
}

func extractDeps(option bool, results fetchFileResult) []schema.ModID {
	dt := NewDepTree[schema.ModID]()
	for modID, result := range results {
		if result.Value == nil {
			continue
		}
		dt.AddNode(Dep(modID, util.Map(
			func(dep schema.FileDependency) schema.ModID { return dep.ModID },
			util.Filter(
				func(dep schema.FileDependency) bool {
					switch dep.RelationType {
					case enum.RequiredDependency:
						return true
					case enum.OptionalDependency:
						return option
					default:
						return false
					}
				},
				result.Value.Dependencies...,
			)...,
		)...))
	}
	return dt.TopSort()
}
