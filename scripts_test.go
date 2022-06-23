package redisLuaScriptUtils_test

import (
	"reflect"
	"testing"

	redisLuaScriptUtils "github.com/zavitax/redis-lua-script-utils-go"
)

func TestScripts(t *testing.T) {
	scriptText1 := `redis.call('SET', key1, arg1);`
	scriptText2 := `redis.call('SET', key2, arg2);`
	scriptText3 := `redis.call('SET', key2, arg2);`

	script1 := redisLuaScriptUtils.NewRedisScript(scriptText1, []string{"key1"}, []string{"arg1"})
	script2 := redisLuaScriptUtils.NewRedisScript(scriptText2, []string{"key2"}, []string{"arg2"})
	script3 := redisLuaScriptUtils.NewRedisScript(scriptText3, []string{"key2"}, []string{"arg2"})

	joinedScript := redisLuaScriptUtils.JoinRedisScripts([]*redisLuaScriptUtils.RedisScript{script1, script2, script3})

	if joinedScript.String() != (script1.String() + "\n" + script2.String() + "\n" + script3.String()) {
		t.Errorf("Joined script value does not match expected value")
		return
	}

	keys := joinedScript.Keys()
	expectedKeys := []string{"key1", "key2"}

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Joined script keys [%v] does not match expected value [%v]", keys, expectedKeys)
		return
	}

	args := joinedScript.Args()
	expectedArgs := []string{"arg1", "arg2"}

	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Joined script keys [%v] does not match expected value [%v]", args, expectedArgs)
		return
	}
}
