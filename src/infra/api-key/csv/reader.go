package csv

import (
	"context"
	csv2 "encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lugondev/signer-key-manager/src/auth/entities"
	apikey "github.com/lugondev/signer-key-manager/src/infra/api-key"
)

const (
	csvSeparator         = ','
	csvCommentsMarker    = '#'
	csvRowLen            = 4
	csvAPIKeyHash        = 0
	csvUserOffset        = 1
	csvPermissionsOffset = 2
	csvRolesOffset       = 3
)

type Reader struct {
	path string
}

var _ apikey.Reader = &Reader{}

func New(cfg *Config) (*Reader, error) {
	_, err := os.Stat(cfg.Path)
	if err != nil {
		return nil, err
	}

	return &Reader{path: cfg.Path}, nil
}

func (r *Reader) Load(_ context.Context) (map[string]*entities.UserClaims, error) {
	csvfile, err := os.Open(r.path)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	csvReader := csv2.NewReader(csvfile)
	csvReader.Comma = csvSeparator
	csvReader.Comment = csvCommentsMarker

	// Read each line from csv and fill claims
	claims := make(map[string]*entities.UserClaims)
	for {
		cells, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(cells) != csvRowLen {
			return nil, fmt.Errorf("invalid number of cells, should be %d", csvRowLen)
		}

		claims[cells[csvAPIKeyHash]] = &entities.UserClaims{
			Tenant:      cells[csvUserOffset],
			Permissions: strings.Split(cells[csvPermissionsOffset], " "),
			Roles:       strings.Split(cells[csvRolesOffset], " "),
		}
	}

	return claims, nil
}
