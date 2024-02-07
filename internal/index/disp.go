package index

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func renderMods(mods ModIndexes) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"ModID", "Name", "Latest Release Date", "Indirect"})
	t.SortBy([]table.SortBy{{Name: "Name", Mode: table.Asc}, {Name: "ModID", Mode: table.Asc}})
	t.SetAutoIndex(true)

	for modID, mod := range mods {
		date := "⛔Release Not Found⛔"
		if len(mod.Name) == 0 {
			errMsg := "⛔Mod Not Found⛔"
			t.AppendRow(table.Row{modID, errMsg, date, mod.IsDependency})
		} else {
			if len(mod.File.Name) > 0 {
				date = mod.File.Date.Format(time.RFC3339)
			}
			t.AppendRow(table.Row{
				modID,
				mod.Name,
				date,
				mod.IsDependency,
			})
		}
	}
	return t.Render()
}
