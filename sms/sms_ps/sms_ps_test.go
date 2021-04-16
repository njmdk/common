package sms_ps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPost(t *testing.T) {
	r := require.New(t)
	resp, err := Post("11111111111","1234")
	r.NoError(err)
	fmt.Println(*resp)
}

func TestGetResult(t *testing.T) {
	r := require.New(t)
	err := GetResult()
	r.NoError(err)

}
