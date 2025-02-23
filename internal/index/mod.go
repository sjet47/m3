package index

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sjet47/go-curseforge/schema"
	"github.com/sjet47/go-curseforge/schema/enum"
	"github.com/sjet47/m3/internal/util"
	"github.com/pkg/errors"
)

var (
	Mods        ModIndexes
	modsDirPath = filepath.Join(M3Root, "mods")
)

type ModIndexes map[schema.ModID]*Mod

func (i ModIndexes) String() string {
	return renderMods(i)
}

type Mod struct {
	ID           schema.ModID `json:"mod_id"`
	Name         string       `json:"mod_name"`
	Summary      string       `json:"mod_summary"`
	GameVersion  string       `json:"game_version"`
	ModLoader    string       `json:"mod_loader"`
	IsDependency bool         `json:"is_dependency"`
	File         struct {
		ID           schema.FileID `json:"id,omitempty"`
		Name         string        `json:"name,omitempty"`
		ReleaseType  string        `json:"release_type,omitempty"`
		HashMD5      string        `json:"hash_md5,omitempty"`
		HashSHA1     string        `json:"hash_sha1,omitempty"`
		Date         time.Time     `json:"date"`
		DownloadUrl  string        `json:"download_url,omitempty"`
		IsServerPack bool          `json:"is_server_pack"`
	} `json:"file"`
}

func NewMod(modLoader enum.ModLoader, mod *schema.Mod, file *schema.File, isDep bool) *Mod {
	m := new(Mod)
	m.ID = mod.ID
	m.Name = mod.Name
	m.Summary = mod.Summary
	m.GameVersion = string(Meta.GameVersion)
	m.ModLoader = modLoader.String()
	m.IsDependency = isDep
	m.Update(file)
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

func (m *Mod) Update(file *schema.File) {
	if file == nil {
		m.File.Date = time.Now()
		return
	}

	m.File.ID = file.ID
	m.File.Name = file.FileName
	m.File.ReleaseType = file.ReleaseType.String()
	for _, hash := range file.Hashes {
		switch hash.Algo {
		case enum.HashAlgoMD5:
			m.File.HashMD5 = hash.Value
		case enum.HashAlgoSHA1:
			m.File.HashSHA1 = hash.Value
		}
	}
	m.File.Date = file.FileDate
	m.File.DownloadUrl = file.DownloadURL
	m.File.IsServerPack = file.IsServerPack
}

func Remove(modIDs ...int) error {
	deleted := 0
	for _, id := range modIDs {
		mod, ok := Mods[schema.ModID(id)]
		if !ok {
			log.Printf("No such mod id: %d", id)
			continue
		}
		os.Remove(mod.File.Name)
		os.Remove(getModIndexByID(mod.ID))
		delete(Mods, schema.ModID(id))
		deleted++
		fmt.Printf("[%d]%s removed\n", mod.ID, mod.Name)
	}
	fmt.Printf("%d/%d mod(s) removed\n", deleted, len(modIDs))
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
			return nil
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
		path := getModIndexByID(id)
		if err := util.WriteJsonToFile(path, m); err != nil {
			return errors.Wrapf(err, "write mod file at %s error", path)
		}
	}

	return nil
}

func getModIndexByID(modID schema.ModID) string {
	return filepath.Join(modsDirPath, fmt.Sprintf("%d.json", modID))
}
