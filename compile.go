package redisLuaScriptUtils

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

type RedisScriptArguments map[string]interface{}
type CompiledRedisScript struct {
	script      RedisScript
	scriptText  string
	keys        []*RedisKey
	args        []string
	redisScript *redis.Script
	mx          sync.RWMutex
}

func getUniqueKeys(keys []*RedisKey) []*RedisKey {
	uniqueKeys := make(map[string]bool, 0)
	uniqueKeysSlice := []*RedisKey{}

	for _, key := range keys {
		if !uniqueKeys[key.Key()] {
			uniqueKeysSlice = append(uniqueKeysSlice, key)
			uniqueKeys[key.Key()] = true
		} else {
			panic("Duplicate key: " + key.Key())
		}
	}

	return uniqueKeysSlice
}

func CompileRedisScripts(scripts []*RedisScript, keys []*RedisKey) (*CompiledRedisScript, error) {
	suppliedKeys := make(map[string]*RedisKey)

	for _, key := range keys {
		suppliedKeys[key.Key()] = key
	}

	for _, key := range keys {
		if suppliedKeys[key.Key()] == nil {
			return nil, fmt.Errorf("Missing required LUA script key: %v", key)
		}
	}

	uniqueKeys := getUniqueKeys(keys)
	uniqueArgs := getScriptsUniqueArgNames(scripts)

	script := joinRedisScripts(scripts, uniqueKeys, uniqueArgs)

	result := &CompiledRedisScript{
		script:     *script,
		scriptText: script.scriptText,
		keys:       uniqueKeys,
		args:       uniqueArgs,
	}

	return result, nil
}

func (this *CompiledRedisScript) String() string {
	return this.scriptText
}

func (this *CompiledRedisScript) Keys(args *RedisScriptArguments) []string {
	var result []string = []string{}

	for _, key := range this.keys {
		result = append(result, key.Value(args))
	}

	return result
}

func (this *CompiledRedisScript) Args(args *RedisScriptArguments) ([]interface{}, error) {
	var result []interface{} = []interface{}{}

	for _, arg := range this.script.args {
		value, ok := (*args)[arg]

		if !ok {
			return nil, fmt.Errorf("Missing required Redis LUA script argument: %v", arg)
		}

		result = append(result, value)
	}

	return result, nil
}

func (this *CompiledRedisScript) Run(ctx context.Context, client *redis.Client, args *RedisScriptArguments) *redis.Cmd {
	if this.redisScript == nil {
		this.mx.Lock()
		this.redisScript = redis.NewScript(this.scriptText)
		this.mx.Unlock()
	}

	if orderedArgsValues, err := this.Args(args); err == nil {
		result := this.redisScript.Run(ctx, client, this.Keys(args), orderedArgsValues)

		if result.Err() != nil {
			panic(fmt.Sprintf("Script run error: %v\nKeys: %v\nArgs: %v\nScript: %v\n\n", result.Err(), this.Keys(args), orderedArgsValues, this.scriptText))
		}

		return result
	} else {
		panic(err)
	}
}
