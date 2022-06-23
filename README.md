# redis-lua-script-utils-go

Utilities for composing and executing Redis LUA scripts

## Motivation

Redis LUA scripts are great! They allow implementing complex functionality with logic as atomic operations.
But what do you do when you want to execute several scripts as a single logical unit?

One way to do it is to compose longer LUA scripts manually, keeping close track of keys and arguments
which should be passed to each resuling script variation.

Will this work? Yes.
Is it maintainable? Likely not.
So what do you do?

## `redis-lua-script-utils-go` to the rescue!

This package allows maintaining (and testing) separate chunks of LUA scripts, merging them and their Keys and Arguments parameters into one single script unit in your program code.

Once merged, it expects KEYS and ARGV values to be supplied as maps, validating input and executing the complete, merged script against redis in an efficient manner.

## Example

```go
package main

import (
  "context"
  "github.com/go-redis/redis/v8"
  redisLuaScriptUtils "github.com/zavitax/redis-lua-script-utils-go"
)

func main() {
	scriptText1 := `redis.call('SET', key1, arg1);`
	scriptText2 := `redis.call('SET', key2, arg2);`
	scriptText3 := `redis.call('SET', key2, arg2);`

	script1 := redisLuaScriptUtils.NewRedisScript(scriptText1, []string{"key1"}, []string{"arg1"})
	script2 := redisLuaScriptUtils.NewRedisScript(scriptText2, []string{"key2"}, []string{"arg2"})
	script3 := redisLuaScriptUtils.NewRedisScript(scriptText3, []string{"key2"}, []string{"arg2"})

	joinedScript := redisLuaScriptUtils.JoinRedisScripts([]*redisLuaScriptUtils.RedisScript{script1, script2, script3})

	compiled, err := redisLuaScriptUtils.CompileRedisScript(
		joinedScript,
		[]*redisLuaScriptUtils.RedisKey{
			redisLuaScriptUtils.NewStaticKey("key1", "keyName1"),
			redisLuaScriptUtils.NewStaticKey("key2", "keyName2"),
		},
	)

	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})

	scriptArgs := make(redisLuaScriptUtils.RedisScriptArguments, 0)
	scriptArgs["arg1"] = "arg1_expected_value"
	scriptArgs["arg2"] = "arg2_expected_value"

	joinedScript.Run(context.TODO(), redisClient, &scriptArgs).Result()

	redisClient.Close()
}
```