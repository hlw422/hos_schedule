-- KEYS[1]: schedule:{id}:remain
-- ARGV[1]: 扣减数量（通常为1）
local key = KEYS[1]
local count = tonumber(ARGV[1])
local remain = tonumber(redis.call('get', key))
if remain >= count then
    redis.call('decrby', key, count)
    return 1
else
    return 0
end
