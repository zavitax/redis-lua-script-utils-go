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

func getScriptsUniqueArgNames(scripts []*RedisScript) []string {
	uniqueArgs := make(map[string]bool, 0)
	uniqueArgsSlice := []string{}

	for _, script := range scripts {
		for _, key := range script.args {
			if !uniqueArgs[key] {
				uniqueArgsSlice = append(uniqueArgsSlice, key)
				uniqueArgs[key] = true
			}
		}
	}

	return uniqueArgsSlice
}

func joinRedisScripts(scripts []*RedisScript, keys []*RedisKey, args []string) *RedisScript {
	result := &RedisScript{}

	var functionCalls []string

	for scriptIndex, script := range scripts {
		compiledArgs := ""
		for argIndex, arg := range args {
			result.args = append(result.args, arg)
			compiledArgs = compiledArgs + fmt.Sprintf("local %s = ARGV[%d];\n", arg, argIndex+1)
		}

		compiledKeys := ""
		for keyIndex, key := range keys {
			result.keys = append(result.keys, key.Key())
			compiledKeys = compiledKeys + fmt.Sprintf("local %s = KEYS[%d];\n", key.Key(), keyIndex+1)
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
