package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LockerV1 struct {
	store *sql.DB

	nodeName string
}

func NewLockerV1(name string) (*LockerV1, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_practise")
	if err != nil {
		return nil, err
	}

	return &LockerV1{
		store:    db,
		nodeName: name,
	}, nil
}

func (l *LockerV1) Lock(ctx context.Context, k string, ttl int) error {

	now := time.Now().Unix()
	_, err := l.store.Exec("insert into tbl_lock_info_v1  (holder, lock_key, ttl, last_lock_time) values (?, ?, ?, ?)", l.nodeName, k, ttl, now)
	if err == nil {
		return nil
	}

	//锁已经存在
	if isDupError(err) {

		expireTime := now - int64(ttl)

		//判断任务是否过期,若过期抢占锁
		r2, err := l.store.Exec("update tbl_lock_info_v1 set last_lock_time = ?, holder = ?  where lock_key = ? and last_lock_time <= ?", now, l.nodeName, k, expireTime)
		if err != nil {
			return err
		}
		preempRet, err := r2.RowsAffected()
		if err != nil {
			return err
		}
		if preempRet == 1 {
			return nil
		}

		return errors.New("抢占失败")
	}
	return err
}
func (l *LockerV1) Renewal(ctx context.Context, k string, ttl int) error {
	now := time.Now().Unix()
	//续约 - 延长过期时间，
	expireTime := now - int64(ttl)
	r1, err := l.store.Exec("update tbl_lock_info_v1 set last_lock_time = ? where holder = ? and lock_key = ? and last_lock_time >= ?", now, l.nodeName, k, expireTime)
	if err != nil {
		return err
	}
	reEntrance, err := r1.RowsAffected()
	if err != nil {
		return err
	}
	if reEntrance == 1 {
		return nil
	}
	//到这说明续约失败
	return errors.New("续约失败")
}
func (l *LockerV1) UnLock(ctx context.Context, k string) error {

	r, err := l.store.Exec("update tbl_lock_info_v1 set holder = ?, last_lock_time = ? where lock_key = ? ", "", 0, k)
	if err != nil {
		return err
	}
	releaseRet, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if releaseRet == 0 {
		return errors.New("释放失败")
	}
	return nil
}

func isDupError(e error) bool {
	if e == nil {
		return false
	}

	return strings.Contains(e.Error(), "Duplicate")
}
