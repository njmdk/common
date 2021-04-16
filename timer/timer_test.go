package timer

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/njmdk/common/utils"
)

func TestNow(t *testing.T) {
	//loc := time.FixedZone("EDT", -4*60*60)
	//fmt.Println(Now().In(loc))
	//fmt.Println(Now())
	//tt, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-12-18 22:14:57", loc)
	//r.NoError(err)
	//fmt.Println(tt.In(time.Local))
	//fmt.Println(UTCAdd8ToUTCSub4Str(NowToString()))
	//uu, err := UTCAdd8ToUTCSub4Str(NowToString())
	//r.NoError(err)
	//fmt.Println(UTCSub4ToUTCAdd8Str(uu))
	//fmt.Println(UTCAdd8ToUTCSub4Time(Now()))
	//fmt.Println(UTCSub4ToUTCAdd8Time(UTCAdd8ToUTCSub4Time(Now())))
	////day := time.Now().Day()
	//
	//y, m, _ := time.Now().Date()
	//fmt.Println(time.Date(y, m, 5, 0, 0, 0, 0, time.Local))
	//fmt.Println(StringToDate(strings.Split("2019-11-11", " ")[0]))
	//fmt.Println(NowUnixMillisecond())
	//fmt.Println(NowUnixNanoSecond())
	//fmt.Println(NowUnixMicrosecond())
	//fmt.Println(NowUnixSecond())
	//fmt.Println(ToString(time.Now()))
	//fmt.Println(TimeToDayInt64(time.Now()))
	//fmt.Println(NowToDayInt64())
	//fmt.Println(TimeToDayInt64(time.Now().Add(-time.Hour * 14)))
	//fmt.Println(TimeToDayInt64(time.Now().Add(-time.Hour * 15)))
	//fmt.Println(TimeToDayInt64(time.Now().Add(time.Hour * 9)))
	//fmt.Println(TimeToDayInt64(time.Now().Add(time.Hour * 10)))
	//fmt.Println(time.Now().Format("2006-01-02 00:00:00"))
	//
	//a := []int{1, 2, 3, 4, 5, 6}
	//lenA := len(a)
	//for i := 0; i < lenA; i++ {
	//	if a[i] == 6 {
	//		a = append(a[:i], a[i+1:]...)
	//		i--
	//		lenA--
	//	}
	//}
	//fmt.Println(a, a[5:])
	
	//ti,err:=time.ParseInLocation("2006-01-02T15:04:05",strings.ReplaceAll("2020-08-04T15:00:00-04:00","-04:00",""),time.FixedZone("UTC",-4*3600))
	//r.NoError(err)
	//ti = ti.In(time.Local)
	//fmt.Println(ti)
	fmt.Println(time.Now().Add(time.Hour*8).Add(-time.Hour*11).UnixNano()/int64(time.Hour*24))
}
func HidePhone(phone string) string {
	lp := len(phone)
	if lp == 11 {
		return phone[:3] + "****" + phone[7:]
	}
	if lp < 11 && lp > 4 {
		return strings.Repeat("*", lp-4) + phone[lp-4:]
	}
	if lp > 11 {
		return phone[:3] + strings.Repeat("*", lp-7) + phone[lp-7+3:]
	}
	return phone
}
func ConvertMapName(name string) string {
	if strings.HasPrefix(name, "de_") {
		name = name[len("de_"):]
		n := []byte(name)
		n[0] -= 'a' - 'A'
		name = utils.BytesToString(n)
	}
	return name
}
