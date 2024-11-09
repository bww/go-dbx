package persist

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/bww/go-dbx/v1/errors"
	"github.com/bww/go-dbx/v1/persist/pql"
	"github.com/jmoiron/sqlx"
)

type Stmt interface {
	Select(interface{}, ...interface{}) error
}

type stmt struct {
	// The following fields are immutable
	pst *persister
	pql string

	// The following fields are mutable and synchronized by the Once; the
	// statement should only be accessed via the stmt() method, which will
	// deal with synchronization and initialization.
	init sync.Once
	sql  string
	raw  *sqlx.Stmt
	typ  reflect.Type
	cols []string
	many bool
}

func (s *stmt) stmt(val reflect.Value) (*sqlx.Stmt, error) {
	var err error
	s.init.Do(func() {
		s.typ, s.many = s.pst.typeMeta(val)
		s.cols = s.pst.typeCols(s.typ)
		prg, err := pql.Parse(s.pql)
		if err != nil {
			err = errors.NewWithSQL(err, s.pql)
			return
		}
		s.sql, err = prg.Text(pql.Context{Columns: s.cols})
		if err != nil {
			err = errors.NewWithSQL(err, s.pql)
			return
		}
		s.raw, err = s.pst.Preparex(s.sql)
		if err != nil {
			err = errors.NewWithSQL(err, s.sql)
			return
		}
	})
	return s.raw, err
}

func (s *stmt) Select(ent interface{}, args ...interface{}) error {
	val := reflect.ValueOf(ent)
	raw, err := s.stmt(val)
	if err != nil {
		return fmt.Errorf("Could not prepare statement: %w", err)
	}
	typ, _ := s.pst.typeMeta(reflect.ValueOf(ent))
	if !typ.AssignableTo(s.typ) {
		return fmt.Errorf("Initialization entity and parameter entity do not agree; you must always provide a destination pointer of the same type to a statement: %v != %v", s.typ, typ)
	}
	if s.many {
		raws, err := raw.Queryx(args...)
		if err != nil {
			return errors.NewWithSQL(err, s.sql)
		}
		return s.pst.selectMany(ent, val, s.cols, s.sql, raws)
	} else {
		return s.pst.selectOne(ent, val, s.cols, s.sql, s.raw.QueryRowx(args...))
	}
}
