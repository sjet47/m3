package mod

import "github.com/ASjet/go-curseforge"

var (
	cli *curseforge.Client
)

func Init(apiKey string) {
	cli = curseforge.NewClient(apiKey)
}
