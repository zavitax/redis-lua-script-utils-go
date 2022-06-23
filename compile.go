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

	compiledArgs := ""

	argIndex := 0
	for _, arg := range script.args {
		argIndex++

		compiledArgs = compiledArgs + fmt.Sprintf("%s = ARGV[%d];\n", arg, argIndex)
	}

	compiledKeys := ""

	keyIndex := 0
	for _, key := range keys {
		keyIndex++

		compiledKeys = compiledKeys + fmt.Sprintf("%s = KEYS[%d];\n", key.Key(), keyIndex)
	}

	if len(compiledArgs) > 0 {
		result.scriptText = compiledArgs + "\n" + result.scriptText
	}

	if len(compiledKeys) > 0 {
		result.scriptText = compiledKeys + "\n" + result.scriptText
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

	if orderedArgsValues, err := this.Args(args); err != nil {
		return this.redisScript.Run(ctx, client, this.Keys(args), orderedArgsValues)
	} else {
		panic(err)
	}
}
