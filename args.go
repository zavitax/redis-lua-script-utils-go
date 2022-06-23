package redisLuaScriptUtils

type RedisArg struct {
	Id          string
	staticValue string
}

func NewStaticArg(id string, value string) *RedisArg {
	return &RedisArg{
		Id:          id,
		staticValue: value,
	}
}

func (this *RedisArg) Value() string {
	return this.staticValue
}
