package help

import "strings"

func Render(appName string, verb string, txt string) string {
	return strings.ReplaceAll(strings.ReplaceAll(txt, "{app_name}", appName), "{verb}", verb)
}
