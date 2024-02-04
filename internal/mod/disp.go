package mod

import (
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	rowConfig = table.RowConfig{
		AutoMerge: true,
	}
)

func renderModInfoTable(modInfoMap fetchModResult, directFileMap, depFileMap fetchFileResult) string {
	t := table.NewWriter()
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"ID", "Name", "Description", "Last Release", "Is Dependency"})
	appendMod(t, modInfoMap, depFileMap, true)
	appendMod(t, modInfoMap, directFileMap, false)
	return t.Render()
}

func appendMod(t table.Writer, modInfoMap fetchModResult, fileMap fetchFileResult, isDep bool) {
	for modID, result := range fileMap {
		info := modInfoMap[modID]
		if info.Err != nil {
			errMsg := info.Err.Error()
			t.AppendRow(table.Row{modID.Param(), errMsg, errMsg, errMsg, isDep}, rowConfig)
		} else {
			mod := info.Value
			date := "No release found"
			file := result.Value
			if file != nil {
				date = file.FileDate.Format(time.RFC3339)
			}
			t.AppendRow(table.Row{
				modID.Param(),
				mod.Name,
				mod.Summary,
				date,
				isDep,
			}, rowConfig)
		}
		t.AppendSeparator()
	}
}
