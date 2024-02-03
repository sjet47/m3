package index

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ASjet/m3/internal/util"

	"github.com/pkg/errors"
)

const (
	M3Root = ".m3"
)

func Init(gameVersion string) error {
	abs, _ := filepath.Abs(M3Root)
	reinit := util.IsDirExists(M3Root)
	if err := os.MkdirAll(M3Root, 0755); err != nil {
		return errors.Wrapf(err, "create m3 index root %s error", abs)
	}

	if err := initMeta(gameVersion); err != nil {
		return errors.Wrap(err, "init m3 index error")
	}

	if err := initMod(); err != nil {
		return errors.Wrap(err, "init m3 index error")
	}

	if reinit {
		fmt.Printf("Reinitialized existing m3 index in %s/\n", abs)
	} else {
		fmt.Printf("Initialized m3 index in %s/\n", abs)
	}
	return nil
}

func Load() error {
	if err := doWith(
		loadMeta,
		loadMods,
	); err != nil {
		return errors.Wrap(err, "Load m3 index error")
	}
	return nil
}

func Save() error {
	if err := doWith(
		saveMeta,
		saveMods,
	); err != nil {
		return errors.Wrap(err, "Save m3 index error")
	}
	return nil
}

func doWith(fs ...func() error) error {
	for _, f := range fs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
