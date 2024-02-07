package mod

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func promptDownload(confirm bool) bool {
	return confirm || strings.ToLower(input("Download mods? [y/N]", "N")) == "y"
}

func input(str, byDefault string) string {
	fmt.Print(str)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if text := scanner.Text(); len(text) > 0 {
			return text
		} else {
			return byDefault
		}
	}
	return ""
}
