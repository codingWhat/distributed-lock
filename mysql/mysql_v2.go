package mysql

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type LockerV2 struct {
	store       *sql.DB
	tx          *sql.Tx
	isStartedTx bool
}

func NewLockerV2() (*LockerV2, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_practise")
	if err != nil {
		return nil, err
	}

	return &LockerV2{
		store: db,
	}, nil
}

func (l *LockerV2) Acquire(ctx context.Context, k string) error {
	if !l.isStartedTx {
		tx, err := l.store.Begin()
		if err != nil {
			return err
		}

		l.isStartedTx = true
		l.tx = tx
	}

	_, err := l.tx.Query("select * from tbl_lock_info_v2 where lock_key = ? FOR UPDATE ", k)
	return err
}

func (l *LockerV2) Release(ctx context.Context) error {
	if !l.isStartedTx {
		return errors.New("没有开启事物")
	}

	if err := l.tx.Commit(); err != nil {
		return err
	}

	l.isStartedTx = false
	return nil
}

func (l *LockerV2) TryLock(ctx context.Context, k string) error {
	if !l.isStartedTx {
		tx, err := l.store.Begin()
		if err != nil {
			return err
		}

		l.isStartedTx = true
		l.tx = tx
	}

	_, err := l.tx.Query("select * from tbl_lock_info_v2 where lock_key = ? FOR UPDATE NOWAIT ", k)
	return err
}
