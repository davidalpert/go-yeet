package resources

import (
	"os"
)

var fs fileSystem = osFS{}

type fileSystem interface {
	Getwd() (string, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Getwd() (string, error) { return os.Getwd() }
