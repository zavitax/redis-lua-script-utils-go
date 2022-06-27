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
	keys        map[string]*RedisKey
	redisScript *redis.Script
	mx          sync.RWMutex
}

func CompileRedisScripts(scripts []*RedisScript, keys []*RedisKey) (*CompiledRedisScript, error) {
	script := joinRedisScripts(scripts)

	suppliedKeys := make(map[string]*RedisKey)

	for _, key := range keys {
		suppliedKeys[key.Key()] = key
	}

	result := &CompiledRedisScript{
		script:     *script,
		scriptText: script.scriptText,
		keys:       make(map[string]*RedisKey, 0),
	}

	for _, key := range script.keys {
		if suppliedKeys[key] == nil {
			return nil, fmt.Errorf("Missing required LUA script key: %v", key)
		}

		result.keys[key] = suppliedKeys[key]
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
