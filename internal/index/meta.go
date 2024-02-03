package index

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/ASjet/m3/internal/util"

	"github.com/ASjet/go-curseforge/schema"
	"github.com/pkg/errors"
)

var (
	Meta *meta

	versionPtn   = regexp.MustCompile(`^(\d+\.\d+)(\.\d+)*$`)
	metaFilePath = filepath.Join(M3Root, "meta.json")
)

type meta struct {
	GameVersion schema.GameVersionStr `json:"game_version"`
}

func initMeta(gameVersion string) error {
	if !isValidVersion(gameVersion) {
		return fmt.Errorf("invalid game version %q", gameVersion)
	}

	Meta = &meta{GameVersion: schema.GameVersionStr(gameVersion)}
	return saveMeta()
}

func loadMeta() error {
	m := new(meta)
	if err := util.ReadJsonFromFile(metaFilePath, m); err != nil {
		return errors.Wrapf(err, "read index meta file at %s error", metaFilePath)
	}
	Meta = m

	return nil
}

func saveMeta() error {
	if err := util.WriteJsonToFile(metaFilePath, Meta); err != nil {
		return errors.Wrapf(err, "write index meta at %s error", metaFilePath)
	}
	return nil
}

func isValidVersion(version string) bool {
	return versionPtn.MatchString(version)
}
