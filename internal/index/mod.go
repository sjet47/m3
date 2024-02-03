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
	"github.com/ASjet/go-curseforge/schema/enum"
	"github.com/ASjet/m3/internal/util"
	"github.com/pkg/errors"
)

var (
	Mods        map[schema.ModID]*Mod
	modsDirPath = filepath.Join(M3Root, "mods")
)

type Mod struct {
	ID          schema.ModID `json:"mod_id"`
	Name        string       `json:"mod_name"`
	Summary     string       `json:"mod_summary"`
	GameVersion string       `json:"game_version"`
	ModLoader   string       `json:"mod_loader"`
	File        struct {
		ID           schema.FileID `json:"id,omitempty"`
		Name         string        `json:"name,omitempty"`
		ReleaseType  string        `json:"release_type,omitempty"`
		Hash         string        `json:"hash,omitempty"`
		Date         time.Time     `json:"date"`
		DownloadUrl  string        `json:"download_url,omitempty"`
		IsServerPack bool          `json:"is_server_pack"`
	} `json:"file"`
}

func NewMod(modLoader enum.ModLoader, mod *schema.Mod, file *schema.File) *Mod {
	m := new(Mod)
	m.ID = mod.ID
	m.Name = mod.Name
	m.Summary = mod.Summary
	m.GameVersion = string(Meta.GameVersion)
	m.ModLoader = modLoader.String()
	m.File.ID = file.ID
	m.File.Name = file.FileName
	m.File.ReleaseType = file.ReleaseType.String()
	m.File.Hash = file.Hashes[0].Value
	m.File.Date = file.FileDate
	m.File.DownloadUrl = file.DownloadURL
	m.File.IsServerPack = file.IsServerPack
	return m
}

func EmptyMod(modLoader enum.ModLoader, modID schema.ModID) *Mod {
	mod := new(Mod)
	mod.ID = modID
	mod.GameVersion = string(Meta.GameVersion)
	mod.ModLoader = modLoader.String()
	mod.File.Date = time.Now()
	return nil
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
	Mods = make(map[schema.ModID]*Mod)
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

		Mods[schema.ModID(modID)] = m
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
