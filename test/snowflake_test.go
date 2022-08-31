package test

import (
	"snowflake/snowflake"
	"testing"
	"time"
)

func TestNextId(t *testing.T) {
	t.Run("generate id", func(t *testing.T) {
		snow := snowflake.Snowflake{}
		workId := 2
		id, err := snow.NextId(workId)
		if err != nil {
			t.Errorf("geneated id failed, error is %v", err)
		}
		seqId := snowflake.SEQUENCE_MASK & id
		if seqId != 1 {
			t.Errorf("seqId should be 1, but got %v", seqId)
		}

		workId2 := (int64(-1) ^ (int64(-1) << snowflake.TIME_SHIFT) - int64(snowflake.SEQUENCE_MASK)) & id >> snowflake.SEQUENCE_BITS
		if workId2 != int64(workId) {
			t.Errorf("want work id %v, but got %v", workId, workId2)
		}

		times := ((int64(-1) << snowflake.TIME_SHIFT) & id) >> snowflake.TIME_SHIFT
		now_diff := time.Now().UnixMilli() - snowflake.SYS_START_TIME
		if times > now_diff {
			t.Errorf("time(%v) can't great then now(%v)", times, now_diff)
		}
	})
}
