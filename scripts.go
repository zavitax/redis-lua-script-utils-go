package redisLuaScriptUtils

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
		}

		if len(result.scriptText) > 0 {
			result.scriptText = result.scriptText + "\n" + script.scriptText
		} else {
			result.scriptText = script.scriptText
		}
	}

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
