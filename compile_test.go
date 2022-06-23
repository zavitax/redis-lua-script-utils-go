package redisLuaScriptUtils_test

import (
	"reflect"
	"testing"

	redisLuaScriptUtils "github.com/zavitax/redis-lua-script-utils-go"
)

func TestCompileScript(t *testing.T) {
	scriptText1 := `redis.call('SET', key1, arg1);`
	scriptText2 := `redis.call('SET', key2, arg2);`
	scriptText3 := `redis.call('SET', key2, arg2);`

	script1 := redisLuaScriptUtils.NewRedisScript([]string{"key1"}, []string{"arg1"}, scriptText1)
	script2 := redisLuaScriptUtils.NewRedisScript([]string{"key2"}, []string{"arg2"}, scriptText2)
	script3 := redisLuaScriptUtils.NewRedisScript([]string{"key2"}, []string{"arg2"}, scriptText3)

	compiled, err := redisLuaScriptUtils.CompileRedisScripts(
		[]*redisLuaScriptUtils.RedisScript{script1, script2, script3},
		[]*redisLuaScriptUtils.RedisKey{
			redisLuaScriptUtils.NewStaticKey("key1", "keyName1"),
			redisLuaScriptUtils.NewStaticKey("key2", "keyName2"),
		},
	)

	if err != nil {
		t.Error(err)
		return
	}

	expectedScript := `local key1 = KEYS[1];
local key2 = KEYS[2];

local arg1 = ARGV[1];
local arg2 = ARGV[2];

redis.call('SET', key1, arg1);
redis.call('SET', key2, arg2);
redis.call('SET', key2, arg2);`

	if compiled.String() != expectedScript {
		t.Errorf("Expected compiled.Script() to match expected script: [%s] != [%s]", compiled.String(), expectedScript)
		return
	}

	argsValues := make(redisLuaScriptUtils.RedisScriptArguments, 0)
	argsValues["arg1"] = "arg1_expected_value"
	argsValues["arg2"] = "arg2_expected_value"
	if _, err := compiled.Args(&argsValues); err != nil {
		t.Error(err)
	}

	keys := compiled.Keys(&argsValues)

	if !reflect.DeepEqual(keys, []string{"keyName1", "keyName2"}) {
		t.Error("Expected keys to match expected keys")
	}

	delete(argsValues, "arg2")

	if _, err := compiled.Args(&argsValues); err == nil {
		t.Error("Expected missing argument to yield error")
	}
}
