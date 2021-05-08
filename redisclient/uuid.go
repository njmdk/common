package redisclient

import (
	"math"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"

	"github.com/njmdk/common/crypt"
	"github.com/njmdk/common/timer"
)

const redisUUIDKey = "redis_uuid_key"

type RedisUUID struct {
	uid   uint64
	index uint32

	smallUid int64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var defaultGen *RedisUUID

func InitDefaultUUIDCreator(client redis.Cmdable) error {
	var err error
	defaultGen, err = NewRedisUUID(client)
	if err != nil {
		return err
	}
	return nil
}

func UUIDNext() string {
	return defaultGen.Next()
}

func UUIDNextN(n int32) []string {
	return defaultGen.NextN(n)
}

func UUIDNext32() string {
	return defaultGen.Next32()
}

func UUIDNext11() string {
	return defaultGen.Next11()
}

func UUIDNext28() string {
	return defaultGen.Next28()
}

func NewRedisUUID(client redis.Cmdable) (*RedisUUID, error) {
	v, e := client.Incr(redisUUIDKey).Result()
	if e != nil {
		return nil, e
	}
	uid := v<<16 + int64(rand.Int31n(math.MaxInt16))
	return &RedisUUID{uid: uint64(uid) << 32, smallUid: v << 32}, nil
}

func (this_ *RedisUUID) Next() string {
	low := atomic.AddUint32(&this_.index, 1)
	uuid := this_.uid + uint64(low)

	return timer.Now().Format("060102150405") + strconv.Itoa(int(uuid))
}

func (this_ *RedisUUID) Next11() string {
	const baseStr = "00000000000"
	low := atomic.AddUint32(&this_.index, 1)
	uuid := this_.smallUid + int64(low)
	uuidStr := crypt.Base62Encode(uuid)
	uuidStrLen := len(uuidStr)
	if uuidStrLen > 11 {
		return uuidStr[:11]
	}
	return baseStr[uuidStrLen:] + uuidStr
}

func (this_ *RedisUUID) Next32() string {
	const baseStr = "00000000000000000000000000000000"
	low := atomic.AddUint32(&this_.index, 1)
	uuid := this_.uid + uint64(low)
	uuidStr := timer.Now().Format("060102150405") + strconv.Itoa(int(uuid))
	uuidStrLen := len(uuidStr)
	if uuidStrLen > 32 {
		return uuidStr[:32]
	}
	return uuidStr + baseStr[uuidStrLen:]
}

func (this_ *RedisUUID) Next28() string {
	const baseStr = "0000000000000000000000000000"
	low := atomic.AddUint32(&this_.index, 1)
	uuid := this_.smallUid + int64(low)
	uuidStr := crypt.Base62Encode(uuid)
	uuidStrLen := len(uuidStr)
	if uuidStrLen > 28 {
		return uuidStr[:28]
	}
	return baseStr[uuidStrLen:] + uuidStr
}

func (this_ *RedisUUID) NextN(n int32) []string {
	if n <= 0 {
		panic("n must more than the 0")
	}
	last := atomic.AddUint32(&this_.index, uint32(n))
	out := make([]string, 0, int(n))
	t := timer.Now().Format("060102150405")
	for low := last - uint32(n) + 1; low <= last; low++ {
		uuid := this_.uid + uint64(low)
		out = append(out, t+strconv.Itoa(int(uuid)))
	}
	return out
}
