package builtin

import (
	"os"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/parser"
)

// BasicMaker represents common members for any maker
type BasicMaker struct {
	TypeHolder *parser.TypeHolder
	Output     *io.GoFile
	Template   string
}

func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
