package index

import (
	"fmt"
	"m3/internal/util"
	"path/filepath"
	"regexp"

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

	m := &meta{GameVersion: schema.GameVersionStr(gameVersion)}
	if err := util.WriteJsonToFile(metaFilePath, m); err != nil {
		return errors.Wrapf(err, "write index meta at %s err:", metaFilePath)
	}
	Meta = m

	return nil
}

func loadMeta() error {
	m := new(meta)
	if err := util.ReadJsonFromFile(metaFilePath, m); err != nil {
		return errors.Wrapf(err, "read index meta file at %s err:", metaFilePath)
	}
	Meta = m

	return nil
}

func isValidVersion(version string) bool {
	return versionPtn.MatchString(version)
}
