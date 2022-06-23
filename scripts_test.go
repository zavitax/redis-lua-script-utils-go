package redisLuaScriptUtils_test

import (
	"reflect"
	"testing"

	redisLuaScriptUtils "github.com/zavitax/redis-lua-script-utils-go"
)

func TestScripts(t *testing.T) {
	scriptText1 := `redis.call('SET', key1, arg1);`

	script := redisLuaScriptUtils.NewRedisScript([]string{"key1"}, []string{"arg1"}, scriptText1)

	keys := script.Keys()
	expectedKeys := []string{"key1"}

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Joined script keys [%v] does not match expected value [%v]", keys, expectedKeys)
		return
	}

	args := script.Args()
	expectedArgs := []string{"arg1"}

	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Joined script keys [%v] does not match expected value [%v]", args, expectedArgs)
		return
	}
}
