package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
)

const maxRetry = 2

var defaultMysql *MySQL

type MySQL struct {
	*sqlx.DB
	log *logger.Logger
}

func Init(cfg *Config, log *logger.Logger) {
	defaultMysql = NewMySQL(cfg, log)

	err := defaultMysql.Ping()
	if err != nil {
		panic(err)
	}
}

func NewMySQL(cfg *Config, log *logger.Logger) *MySQL {
	config := mysql.NewConfig()
	config.ParseTime = true
	config.Net = "tcp"
	config.Addr = cfg.Addr
	config.User = cfg.User
	config.Passwd = cfg.Password
	config.DBName = cfg.Database
	config.Loc = time.Local
	config.Collation = "utf8mb4_general_ci"
	config.AllowNativePasswords = true
	config.InterpolateParams = true
	dsn := config.FormatDSN()
	db := sqlx.MustConnect("mysql", dsn)
	db.SetConnMaxLifetime(cfg.MaxLifeTime.Duration)
	db.SetMaxIdleConns(cfg.MaxIdleConnS)
	db.SetMaxOpenConns(cfg.MaxOpenConnS)
	db.Mapper = reflectx.NewMapperFunc("json", func(s string) string {
		return s
	})

	return &MySQL{DB: db, log: log}
}

func MustQuerySlice(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) (foundRows int64) {
	return defaultMysql.MustQuerySlice(dealRows, query, args...)
}

func (this_ *MySQL) GetLogger() *logger.Logger {
	return this_.log
}

func (this_ *MySQL) MustQuerySlice(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) (foundRows int64) {
	var (
		rows *sqlx.Rows
		err  error
	)

	for i := 0; i < maxRetry; i++ {
		rows, err = this_.Queryx(query, args...)
		if err != mysql.ErrInvalidConn {
			break
		}
	}

	if err == mysql.ErrInvalidConn {
		rows, err = this_.Queryx(query, args...)
	}

	if err != nil {
		this_.log.Panic("query slice error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			this_.log.Error("query rows.close error", zap.Error(err))
		}
	}()

	for rows.Next() {
		foundRows++

		if dealRows != nil {
			err = dealRows(rows)
			if err != nil {
				this_.log.Panic("query slice error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
			}
		}
	}

	return foundRows
}

func QuerySlice(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) error {
	return defaultMysql.QuerySlice(dealRows, query, args...)
}

func (this_ *MySQL) QuerySlice(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) error {
	var (
		rows *sqlx.Rows
		err  error
	)

	for i := 0; i < maxRetry; i++ {
		rows, err = this_.Queryx(query, args...)
		if err != mysql.ErrInvalidConn {
			break
		}
	}

	if err == mysql.ErrInvalidConn {
		rows, err = this_.Queryx(query, args...)
	}

	if err != nil {
		this_.log.Error("query slice error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
		return err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			this_.log.Error("query rows.close error", zap.Error(err))
		}
	}()

	if dealRows != nil {
		for rows.Next() {
			err = dealRows(rows)
			if err != nil {
				this_.log.Warn("query slice error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
				return err
			}
		}
	}

	return nil
}

func MustQuery(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) (foundRows int64) {
	return defaultMysql.MustQuery(dealRows, query, args...)
}

func (this_ *MySQL) MustQuery(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) (foundRows int64) {
	var (
		rows *sqlx.Rows
		err  error
	)

	for i := 0; i < maxRetry; i++ {
		rows, err = this_.Queryx(query, args...)
		if err != mysql.ErrInvalidConn {
			break
		}
	}

	if err == mysql.ErrInvalidConn {
		rows, err = this_.Queryx(query, args...)
	}

	if err != nil {
		this_.log.Panic("query error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			this_.log.Error("query rows.close error", zap.Error(err))
		}
	}()

	if rows.Next() {
		foundRows++

		err := dealRows(rows)
		if err != nil {
			this_.log.Panic("query error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
		}
	}

	return foundRows
}

func Query(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) (foundRows int64, err error) {
	return defaultMysql.Query(dealRows, query, args...)
}

func (this_ *MySQL) Query(dealRows func(rows *sqlx.Rows) error, query string, args ...interface{}) (foundRows int64, err error) {
	var rows *sqlx.Rows
	for i := 0; i < maxRetry; i++ {
		rows, err = this_.Queryx(query, args...)
		if err != mysql.ErrInvalidConn {
			break
		}
	}

	if err == mysql.ErrInvalidConn {
		rows, err = this_.Queryx(query, args...)
	}

	if err != nil {
		this_.log.Error("query error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
		return 0, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			this_.log.Error("query rows.close error", zap.Error(err))
		}
	}()

	if rows.Next() {
		foundRows++

		err := dealRows(rows)
		if err != nil {
			this_.log.Warn("query error", zap.Error(err), zap.String("sql", query), zap.Any("args", args))
			return 0, err
		}
	}

	return foundRows, nil
}

func MustExec(query string, args ...interface{}) (lastInsertID, rowsAffected int64) {
	return defaultMysql.MustExec(query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return defaultMysql.Exec(query, args...)
}

func (this_ *MySQL) LogExec(query string, args ...interface{}) (lastInsertID, rowsAffected int64, err error) {
	r, err := this_.DB.Exec(query, args...)
	if err != nil {
		this_.log.Error("exec sql error:", zap.String("sql", query), zap.Any("param", args), zap.Error(err))
		return 0, 0, err
	}
	lastInsertID, _ = r.LastInsertId()
	rowsAffected, _ = r.RowsAffected()
	return lastInsertID, rowsAffected, nil
}
func (this_ *MySQL) MustExec(query string, args ...interface{}) (lastInsertID, rowsAffected int64) {
	var (
		result sql.Result
		err    error
	)

	for i := 0; i < maxRetry; i++ {
		result, err = this_.DB.Exec(query, args...)
		if err != mysql.ErrInvalidConn {
			break
		}
	}

	if err == mysql.ErrInvalidConn {
		result, err = this_.DB.Exec(query, args...)
	}

	if err != nil {
		this_.log.Panic("MustExec panic", zap.Error(err), zap.String("query", query), zap.Any("args", args))
	}

	lastInsertID, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()

	return
}

func (this_ *MySQL) BeginTx(f func(tx *sqlx.Tx) error) (err error) {
	if f == nil {
		return nil
	}

	var tx *sqlx.Tx
	for i := 0; i < maxRetry; i++ {
		tx, err = this_.Beginx()
		if err != mysql.ErrInvalidConn {
			break
		}
	}

	if err == mysql.ErrInvalidConn {
		tx, err = this_.Beginx()
	}

	if err != nil {
		return err
	}

	defer func() {
		if e := recover(); e != nil {
			this_.log.Error("exec tx panic", zap.Any("panic info", e))
			err = fmt.Errorf("%v", e)
			_ = tx.Rollback()
		}
	}()

	err = f(tx)
	if err != nil {
		if err != sql.ErrNoRows {
			this_.log.Error("exec tx error", zap.Error(err))
		}
		_ = tx.Rollback()

		return err
	}

	_ = tx.Commit()

	return nil
}

func IsDuplicateKeyError(err error) bool {
	e, ok := err.(*mysql.MySQLError)
	if ok && e.Number == 1062 {
		return true
	}

	return false
}

func BeginWithTx(f func(tx *sqlx.Tx) error) (err error) {
	if f == nil {
		return nil
	}

	err = defaultMysql.BeginTx(f)
	if err != nil {
		defaultMysql.log.Panic("BeginWithTx panic", zap.Error(err))
	}

	return
}
