package cache

import (
	"github.com/gomodule/redigo/redis"
	"testing"
)

const (
	testEntry       = "hello_world_test"
	testEntryResult = "how_are_you_starshine?"
)

func TestCreate(t *testing.T) {
	conn, errConn := create(nil)
	if conn != nil {
		t.Error("nil parameters should return nil")
	}
	if errConn == nil {
		t.Error("nil paramters should return error")
	}
}

func TestExec(t *testing.T) {
	setCommands := []interface{}{"SET", testEntry, testEntryResult}
	entry, errEntry := Exec(&setCommands, nil)
	if errEntry != nil {
		t.Fail()
		t.Logf(errEntry.Error())
	}

	if entry == nil {
		t.Fail()
		t.Logf("setter.Set should retrun an entry")
	}

	getCommands := []interface{}{"GET", testEntry}
	getterEntry, errGetterEntry := redis.String(Exec(&getCommands, nil))
	if errGetterEntry != nil {
		t.Fail()
		t.Logf(errGetterEntry.Error())
	}

	if getterEntry != testEntryResult {
		t.Fail()
		t.Logf("setter.Get should equal found count")
	}

}