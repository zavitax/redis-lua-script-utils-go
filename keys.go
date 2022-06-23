package redisLuaScriptUtils

type RedisKeyValueGeneratorFunc func(args *map[string]interface{}) string

type RedisKey struct {
	id           string
	staticValue  string
	dynamicValue RedisKeyValueGeneratorFunc
}

func NewStaticKey(id string, value string) *RedisKey {
	return &RedisKey{
		id:           id,
		staticValue:  value,
		dynamicValue: nil,
	}
}

func NewDynamicKey(id string, generator RedisKeyValueGeneratorFunc) *RedisKey {
	return &RedisKey{
		id:           id,
		dynamicValue: generator,
	}
}

func (this *RedisKey) Key() string {
	return this.id
}

func (this *RedisKey) Value(args *map[string]interface{}) string {
	if this.dynamicValue == nil {
		return this.staticValue
	} else {
		return this.dynamicValue(args)
	}
}
