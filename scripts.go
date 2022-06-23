package redisLuaScriptUtils

type RedisScript struct {
	scriptText string
	args       []string
	keys       []string
}

func NewRedisScript(scriptText string, keys []string, args []string) *RedisScript {
	return &RedisScript{
		scriptText: scriptText,
		args:       args,
		keys:       keys,
	}
}

func JoinRedisScripts(scripts []*RedisScript) *RedisScript {
	result := &RedisScript{}

	uniqueKeys := make(map[string]bool, 0)
	uniqueArgs := make(map[string]bool, 0)

	for _, script := range scripts {
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

			result.scriptText = result.scriptText + "\n" + result.scriptText
		}
	}

	return result
}
