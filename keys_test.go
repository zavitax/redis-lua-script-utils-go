package redisLuaScriptUtils_test

import (
	"testing"

	redisLuaScriptUtils "github.com/zavitax/redis-lua-script-utils"
)

func TestStaticKeys(t *testing.T) {
	key := redisLuaScriptUtils.NewStaticKey("key", "value")

	if key.Key() != "key" {
		t.Error("key.Key() return value should be 'key'")
		return
	}

	argsMap := make(redisLuaScriptUtils.RedisScriptArguments, 0)

	if key.Value(&argsMap) != "value" {
		t.Error("key.Value() return value should be 'value'")
		return
	}
}

func TestDynamicKeys(t *testing.T) {
	key := redisLuaScriptUtils.NewDynamicKey("key", func(args *redisLuaScriptUtils.RedisScriptArguments) string {
		return (*args)["arg1"].(string)
	})

	if key.Key() != "key" {
		t.Error("key.Key() return value should be 'key'")
		return
	}

	argsMap := make(redisLuaScriptUtils.RedisScriptArguments, 0)

	argsMap["arg1"] = "expected_argument_value"

	if key.Value(&argsMap) != "expected_argument_value" {
		t.Error("key.Value() return value should be 'expected_argument_value'")
		return
	}
}
