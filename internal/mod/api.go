package mod

import (
	"github.com/sjet47/go-curseforge/api"
	"github.com/sjet47/go-curseforge/schema"
	"github.com/sjet47/go-curseforge/schema/enum"
	"github.com/sjet47/m3/internal/index"
	"github.com/pkg/errors"
)

func getLatestModFile(modID schema.ModID, modLoader enum.ModLoader) (*schema.File, error) {
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
	if err != nil {
		return nil, err
	}
	return &resp.Data[0], nil
}

func getModInfo(modID schema.ModID) (*schema.Mod, error) {
	resp, err := api.Mod(modID)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
