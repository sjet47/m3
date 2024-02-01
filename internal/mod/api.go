package mod

import "github.com/ASjet/go-curseforge"

func Init(apiKey string) {
	curseforge.InitDefault(apiKey)
}
