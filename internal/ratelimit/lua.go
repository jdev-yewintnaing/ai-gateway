package ratelimit

import "github.com/redis/go-redis/v9"

var IncrementAndCheckLua = redis.NewScript(`
local key = KEYS[1]
local tokens = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local ttl = tonumber(ARGV[3])

local current = redis.call("GET", key)
if current and tonumber(current) + tokens > limit then
    return 0
end

redis.call("INCRBY", key, tokens)
redis.call("EXPIRE", key, ttl)
return 1
`)
