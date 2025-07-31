package dblib

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrLockNotAcquired = errors.New("lock has already been acquired")
var ErrLockAlreadyReleased = errors.New("lock has already been released")

type LockValue struct {
	Type  string
	Value string
	Time  time.Time
}

type Lock interface {
	AcquireLock(lockType string, timeout time.Duration) (*LockValue, error)
	ReleaseLock(lockValue *LockValue) error
}

type LockSql struct {
	ConnStr string
}

func (db *LockSql) GetConnectionString() string {
	return db.ConnStr
}

func (db *LockSql) AcquireLock(lockType string, timeout time.Duration) (*LockValue, error) {
	val := uuid.NewString()
	lockVal := LockValue{}
	err := OpenConnAndExec(db, func(dbc *sql.DB) error {
		row := dbc.QueryRow(`
			INSERT INTO sync_lock AS sl (id, lock_val, lock_time) 
			VALUES ($1, $2, CURRENT_TIMESTAMP)
			ON CONFLICT (id) DO UPDATE  
			SET lock_val = $2, 
			lock_time = CURRENT_TIMESTAMP 
			WHERE sl.id = $1
			AND (sl.lock_val IS NULL OR sl.lock_time + $3 < NOW())
			RETURNING id, lock_val, lock_time`, lockType, val, timeout)

		scanErr := row.Scan(&lockVal.Type, &lockVal.Value, &lockVal.Time)

		if scanErr == sql.ErrNoRows {
			return ErrLockNotAcquired
		}

		return scanErr
	})

	return &lockVal, err
}

func (db *LockSql) ReleaseLock(lockValue *LockValue) error {
	err := OpenConnAndExec(db, func(dbc *sql.DB) error {
		row := dbc.QueryRow(`
			UPDATE sync_lock 
			SET lock_val = NULL, lock_time = NULL 
			WHERE id = $1
			AND lock_val = $2 
			RETURNING id`, lockValue.Type, lockValue.Value)
		id := ""

		scanErr := row.Scan(&id)

		if scanErr == sql.ErrNoRows {
			return ErrLockAlreadyReleased
		}

		return scanErr
	})

	return err
}
