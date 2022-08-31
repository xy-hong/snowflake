package snowflake

import (
	"fmt"
	"time"
)

const (
	TIME_BITS      = 41
	WORKID_BITS    = 10
	SEQUENCE_BITS  = 12
	WORKID_SHIFT   = SEQUENCE_BITS
	TIME_SHIFT     = WORKID_BITS + SEQUENCE_BITS
	WORKID_MAX     = -1 ^ (-1 << WORKID_BITS)
	SEQUENCE_MASK  = -1 ^ (-1 << SEQUENCE_BITS)
	SYS_START_TIME = 1661916745000
)

type Snowflake struct {
	LastTimestamp int64
	SequenceId    int
}

/**
 * ⚠️ 注意事项
 * 1）workid 不能超过其最大值
 * 2) 如果是同一毫秒内，则递增 sequence id，如果溢出则放到下一毫秒
 *
 */
func (snowflake *Snowflake) NextId(workid int) (int64, error) {
	if workid > WORKID_MAX || workid < 0 {
		return 0, fmt.Errorf("work id should bettwen 0 and %v, but got %v", WORKID_BITS, workid)
	}

	now := time.Now().UnixMilli()
	seqId := 0
	if now < snowflake.LastTimestamp {
		return 0, fmt.Errorf("clock moved backwords, refusing to generate id")
	}

	if now == snowflake.LastTimestamp {
		seqId = (snowflake.SequenceId + 1) & SEQUENCE_MASK
		if seqId == 0 {
			now = tilNextMill(snowflake.LastTimestamp)
		}
		seqId = 1

	} else {
		snowflake.LastTimestamp = now
		seqId = 1
	}
	snowflake.SequenceId = seqId

	return (now-SYS_START_TIME)<<TIME_SHIFT | int64(workid)<<WORKID_SHIFT | int64(seqId), nil
}

func tilNextMill(lastMill int64) int64 {
	now := time.Now().UnixMilli()
	for now <= lastMill {
		now = time.Now().UnixMilli()
	}
	return now
}
