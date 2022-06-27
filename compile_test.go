package redisLuaScriptUtils_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
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

func TestFunctions(t *testing.T) {
	scriptText1 := `
		return arg1
	`

	scriptText2 := `
		return arg2
	`

	script1 := redisLuaScriptUtils.NewRedisScript([]string{}, []string{"arg1"}, scriptText1)
	script2 := redisLuaScriptUtils.NewRedisScript([]string{}, []string{"arg2"}, scriptText2)

	compiled, err := redisLuaScriptUtils.CompileRedisScripts(
		[]*redisLuaScriptUtils.RedisScript{script1, script2},
		[]*redisLuaScriptUtils.RedisKey{},
	)

	if err != nil {
		t.Error(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	scriptArgs := make(redisLuaScriptUtils.RedisScriptArguments, 0)
	scriptArgs["arg1"] = "RESULT1"
	scriptArgs["arg2"] = "RESULT2"
	result, err2 := compiled.Run(context.TODO(), redisClient, &scriptArgs).StringSlice()

	if err2 != nil {
		t.Error(err2)
		return
	}

	expectedResult := []string{"RESULT1", "RESULT2"}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result [%v] to match expected result [%v]", result, expectedResult)
	}

	redisClient.Close()
}
