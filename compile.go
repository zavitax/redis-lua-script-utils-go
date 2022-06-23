package redisLuaScriptUtils

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

type CompiledRedisScript struct {
	script      RedisScript
	scriptText  string
	keys        map[string]*RedisKey
	redisScript *redis.Script
	mx          sync.RWMutex
}

func CompileRedisScript(script *RedisScript, keys []*RedisKey) (*CompiledRedisScript, error) {
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

	result.scriptText = "\n" + result.scriptText

	argIndex := 0
	for arg, _ := range script.args {
		argIndex++

		result.scriptText = fmt.Sprintf("%s = ARGV[%d]\n%s", arg, argIndex, result.scriptText)
	}

	result.scriptText = "\n" + result.scriptText

	keyIndex := 0
	for key, _ := range keys {
		keyIndex++

		result.scriptText = fmt.Sprintf("%s = KEYS[%d]\n%s", key, keyIndex, result.scriptText)
	}

	return result, nil
}

func (this *CompiledRedisScript) Script() string {
	return this.scriptText
}

func (this *CompiledRedisScript) Keys(args *map[string]interface{}) []string {
	var result []string = []string{}

	for _, key := range this.keys {
		result = append(result, key.Value(args))
	}

	return result
}

func (this *CompiledRedisScript) Args(args *map[string]interface{}) ([]interface{}, error) {
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

func (this *CompiledRedisScript) Run(ctx context.Context, client *redis.Client, args *map[string]interface{}) *redis.Cmd {
	if this.redisScript == nil {
		this.mx.Lock()
		this.redisScript = redis.NewScript(this.scriptText)
		this.mx.Unlock()
	}

	if orderedArgsValues, err := this.Args(args); err != nil {
		return this.redisScript.Run(ctx, client, this.Keys(args), orderedArgsValues)
	} else {
		panic(err)
	}
}
