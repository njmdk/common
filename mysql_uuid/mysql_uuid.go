package mysql_uuid

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/edwingeng/wuid/mysql/wuid"

	"github.com/njmdk/common/crypt"
	"github.com/njmdk/common/db"
	"github.com/njmdk/common/utils"
)

var baseRandBytes = []byte("ghijklmnopqrstuvwxyz")
var lenBase = len(baseRandBytes)

type UUID struct {
	*wuid.WUID
}

func (u *UUID) Next() int64 {
	return u.WUID.Next()
}

func (u *UUID) NextString() string {
	return strconv.Itoa(int(u.WUID.Next()))
}

func (u *UUID) NextBase62() string {
	n := u.WUID.Next()
	return crypt.Base62Encode(n)
}

func (u *UUID) NextBase34() string {
	n := u.WUID.Next()
	return crypt.Base34(uint64(n))
}

func (u *UUID) Next16String() string {
	bs := []byte(fmt.Sprintf("%016x", u.Next()))
	for k := range bs {
		if bs[k] == '0' {
			bs[k] = baseRandBytes[rand.Intn(lenBase)]
		}
	}
	return utils.BytesToString(bs)
}

func (u *UUID) Next32String() string {
	bs := []byte(fmt.Sprintf("%032x", u.Next()))
	for k := range bs {
		if bs[k] == '0' {
			bs[k] = baseRandBytes[rand.Intn(lenBase)]
		}
	}
	return utils.BytesToString(bs)
}

func CreateUUIDCreator(mysql *db.MySQL) (*UUID, error) {
	rand.Seed(time.Now().UnixNano())
	newDB := func() (*sql.DB, bool, error) {
		return mysql.DB.DB, false, nil
	}

	// Setup
	g := wuid.NewWUID("default", nil)
	err := g.LoadH28FromMysql(newDB, "wuid")
	if err != nil {
		return nil, err
	}
	return &UUID{WUID: g}, nil
}
