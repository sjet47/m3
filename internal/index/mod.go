package index

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ASjet/go-curseforge/schema"
	"github.com/ASjet/m3/internal/util"
	"github.com/pkg/errors"
)

var (
	Mods        map[int]*Mod
	modsDirPath = filepath.Join(M3Root, "mods")
)

type Mod struct {
	ID      schema.ModID `json:"mod_id"`
	Name    string       `json:"mod_name"`
	Summary string       `json:"mod_summary"`
	File    struct {
		ID          schema.FileID `json:"id,omitempty"`
		Name        string        `json:"name,omitempty"`
		ReleaseType string        `json:"release_type,omitempty"`
		Hash        string        `json:"hash,omitempty"`
		Date        time.Time     `json:"date"`
		DownloadUrl string        `json:"download_url,omitempty"`
		GameVersion string        `json:"game_version"`
		ModLoader   string        `json:"mod_loader"`
	} `json:"file"`
}

func initMod() error {
	if !util.IsDirExists(modsDirPath) {
		if err := os.MkdirAll(modsDirPath, 0755); err != nil {
			return errors.Wrap(err, "create mods dir at %s error")
		}
	}
	return nil
}

func loadMods() error {
	Mods = make(map[int]*Mod)
	return filepath.WalkDir(modsDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return filepath.SkipDir
		}

		modIDs := strings.Split(filepath.Base(path), ".")[0]
		modID, err := strconv.Atoi(modIDs)
		if err != nil {
			return errors.Wrapf(err, "parse mod id %s from path error", modIDs)
		}

		m := new(Mod)
		if err := util.ReadJsonFromFile(path, m); err != nil {
			return errors.Wrapf(err, "read mod file at %s error", path)
		}

		Mods[modID] = m
		return nil
	})
}

func saveMods() error {
	if err := initMod(); err != nil {
		return err
	}

	for id, m := range Mods {
		path := filepath.Join(modsDirPath, fmt.Sprintf("%d.json", id))
		if err := util.WriteJsonToFile(path, m); err != nil {
			return errors.Wrapf(err, "write mod file at %s error", path)
		}
	}

	return nil
}
