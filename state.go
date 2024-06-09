package dblite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const driver = "sqlite3"

type StateDB struct {
	DB *sql.DB
}

func New(filename string) (*StateDB, error) {
	db, err := sql.Open(driver, filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open database, err: %v", err)
	}

	createStmt := `create table if not exists state (
		key text not null primary key,
		value text
	)`
	if _, err := db.Exec(createStmt); err != nil {
		return nil, fmt.Errorf("cannot create state db, err: %v", err)
	}

	return &StateDB{DB: db}, nil
}

func (sdb *StateDB) Get(ctx context.Context, key string) (*State, error) {
	queryStmt := "select key, value from state where key = ?"
	rows, err := sdb.DB.QueryContext(ctx, queryStmt, key)
	if err != nil {
		return nil, fmt.Errorf("cannot save state, err: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var s State
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, fmt.Errorf("cannot scan state, err: %v", err)
		}
		return &s, nil
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows error, err: %v", err)
	}

	return nil, nil
}

func (sdb *StateDB) Save(ctx context.Context, s State) error {
	if s.Key == "" {
		return errors.New("key cannot be empty")
	}

	insertStmt := "insert into state(key,value) values (?,?)"
	if _, err := sdb.DB.ExecContext(ctx, insertStmt, s.Key, s.Value); err != nil {
		return fmt.Errorf("cannot save state, err: %v", err)
	}

	return nil
}

func (sdb *StateDB) SaveBulk(ctx context.Context, ss []State) error {
	holders := []string{}
	args := []any{}

	for _, s := range ss {
		if s.Key == "" {
			return errors.New("key cannot be empty")
		}

		holders = append(holders, "(?,?)")
		args = append(args, s.Key, s.Value)
	}

	insertStmt := fmt.Sprintf("insert into state(key,value) values %s", strings.Join(holders, ","))
	if _, err := sdb.DB.ExecContext(ctx, insertStmt, args...); err != nil {
		return fmt.Errorf("cannot bulk save state, err: %v", err)
	}

	return nil
}

func (sdb *StateDB) Delete(ctx context.Context, key string) error {
	deleteStmt := "delete from state where key = ?"
	if _, err := sdb.DB.ExecContext(ctx, deleteStmt, key); err != nil {
		return fmt.Errorf("cannot delete state, err: %v", err)
	}

	return nil
}

type State struct {
	Key   string
	Value string
}
