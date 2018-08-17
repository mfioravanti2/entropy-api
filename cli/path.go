package cli

import (
	"context"
	"fmt"
	"errors"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

const (
	PATH_MODEL		string = "entropy.path.models"
	PATH_OPENAPI	string = "entropy.path.openapi"
)

type Path struct {
	Name string		`json:"name"`
	Path string		`json:"path"`
}

type Paths []Path

func (p *Path) DefaultConfig() {
	p.Name = ""
	p.Path = ""
}

// Create a default empty path
func NewPath() (*Path, error) {
	var p = &Path{}
	p.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewPath", "cli" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default path configuration",
	)

	return p, nil
}

func NewPaths() (Paths, error) {
	var ps = Paths{}
	var p *Path

	p, _ = NewPath()
	p.Name = PATH_MODEL
	ps = append( ps, *p )

	p, _ = NewPath()
	p.Name = PATH_OPENAPI
	ps = append( ps, *p )

	return ps, nil
}

func (ps *Paths) GetPath( name string ) (*Path, error) {
	var p Path

	for _, p = range *ps {
		if p.Name == name {
			return &p, nil
		}
	}

	s := fmt.Sprintf("find path failed (path name: %s)", name )
	return nil, errors.New(s)
}

