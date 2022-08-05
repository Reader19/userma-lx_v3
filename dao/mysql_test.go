package dao

import (
	"strconv"
	"testing"
	"usermaLX4/utils"
)

func TestCreateAccount(t *testing.T) {
	password := utils.MD5("123")
	for i := 0; i < 10000000; i++ {
		username := "user" + strconv.Itoa(i)
		if err := InsertUser(username, password); err != nil {
			t.Errorf("can't insert the account(username: %s, err: %q )", username, err)
		}
		//if i%100000 == 0 {
		//	t.Logf("now is %d", i)
		//}
	}
}
