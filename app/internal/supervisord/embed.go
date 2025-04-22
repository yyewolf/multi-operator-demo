package supervisord

import (
	"embed"
)

//go:embed static
var static embed.FS

func GetStaticFile(name string) ([]byte, error) {
	return static.ReadFile("static/" + name)
}
