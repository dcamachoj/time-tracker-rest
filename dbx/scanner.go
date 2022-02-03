package dbx

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

type eScanner struct {
	query    string
	args     []interface{}
	scanners []interface{}
}

func newScanner(query string, args ...interface{}) EntityScanner {
	return &eScanner{
		query: query,
		args:  args,
	}
}

func (s *eScanner) Scanners() []interface{} {
	return s.scanners
}
func (s *eScanner) AddScanners(scanners ...interface{}) {
	s.scanners = append(s.scanners, scanners...)
}
func (s *eScanner) Scan(ctx context.Context, rows *sql.Rows, entity Entity) (res interface{}, err error) {
	s.scanners = nil
	err = entity.PreScan(s)
	if err != nil {
		return nil, errors.Wrap(err, "Pre Scanning")
	}

	err = rows.Scan(s.scanners...)
	if err != nil {
		dbLog(s.query, s.args...)
		log.Println(err)
		return nil, errors.Wrap(err, "Scanning")
	}

	err = entity.PosScan(s)
	if err != nil {
		return nil, errors.Wrap(err, "Post Scanning")
	}

	return entity.Result(), nil
}
