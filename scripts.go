package redisLuaScriptUtils

import (
	"fmt"
	"strings"
)

type RedisScript struct {
	scriptText string
	args       []string
	keys       []string
}

func NewRedisScript(keys []string, args []string, scriptText string) *RedisScript {
	return &RedisScript{
		scriptText: scriptText,
		args:       args,
		keys:       keys,
	}
}

func joinRedisScripts(scripts []*RedisScript) *RedisScript {
	result := &RedisScript{}

	uniqueKeys := make(map[string]bool, 0)
	uniqueArgs := make(map[string]bool, 0)

	var functionCalls []string

	for scriptIndex, script := range scripts {
		for _, key := range script.keys {
			if !uniqueKeys[key] {
				result.keys = append(result.keys, key)
				uniqueKeys[key] = true
			}
		}

		for _, key := range script.args {
			if !uniqueArgs[key] {
				result.args = append(result.args, key)
				uniqueArgs[key] = true
			}
		}

		compiledArgs := ""
		for argIndex, arg := range result.args {
			compiledArgs = compiledArgs + fmt.Sprintf("local %s = ARGV[%d];\n", arg, argIndex+1)
		}

		compiledKeys := ""
		for keyIndex, key := range result.keys {
			compiledKeys = compiledKeys + fmt.Sprintf("local %s = KEYS[%d];\n", key, keyIndex+1)
		}

		functionName := fmt.Sprintf("____joinedRedisScripts_%d____", scriptIndex)

		envelopedScriptText := fmt.Sprintf("local function %s()\n%s\n%s\n%s\nend", functionName, compiledKeys, compiledArgs, script.scriptText)

		functionCalls = append(functionCalls, fmt.Sprintf("%s()", functionName))

		if len(result.scriptText) > 0 {
			result.scriptText = result.scriptText + "\n" + envelopedScriptText
		} else {
			result.scriptText = envelopedScriptText
		}
	}

	result.scriptText = result.scriptText + "\n" + "return {" + strings.Join(functionCalls, ", ") + "}\n"

	return result
}

func (this *RedisScript) String() string {
	return this.scriptText
}

func (this *RedisScript) Keys() []string {
	return this.keys
}

func (this *RedisScript) Args() []string {
	return this.args
}
