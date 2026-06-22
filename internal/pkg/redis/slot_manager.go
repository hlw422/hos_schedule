package redis

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

//go:embed slot_deduction.lua
var luaScript string

type SlotManager struct {
	rdb *redis.Client
}

func NewSlotManager(rdb *redis.Client) *SlotManager {
	return &SlotManager{rdb: rdb}
}

func (m *SlotManager) InitSlot(ctx context.Context, scheduleID int64, remainCount int) error {
	key := fmt.Sprintf("schedule:%d:remain", scheduleID)
	return m.rdb.Set(ctx, key, remainCount, 0).Err()
}

func (m *SlotManager) DeductSlot(ctx context.Context, scheduleID int64) (bool, error) {
	key := fmt.Sprintf("schedule:%d:remain", scheduleID)
	result, err := m.rdb.Eval(ctx, luaScript, []string{key}, 1).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (m *SlotManager) ReleaseSlot(ctx context.Context, scheduleID int64) error {
	key := fmt.Sprintf("schedule:%d:remain", scheduleID)
	return m.rdb.Incr(ctx, key).Err()
}

func (m *SlotManager) GetRemain(ctx context.Context, scheduleID int64) (int, error) {
	key := fmt.Sprintf("schedule:%d:remain", scheduleID)
	val, err := m.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}
