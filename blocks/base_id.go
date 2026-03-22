package blocks

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var idCounter uint64

// GenerateID creates a unique ID suitable for action_id or block_id.
// Format: prefix_timestamp_counter_random
// Ensures uniqueness across concurrent calls and process restarts.
func GenerateID(prefix string) string {
	count := atomic.AddUint64(&idCounter, 1)

	b := make([]byte, 4)
	_, _ = rand.Read(b)

	ts := time.Now().UnixNano() / 1e6

	if prefix == "" {
		prefix = "id"
	}

	return fmt.Sprintf("%s_%d_%d_%s", prefix, ts, count, hex.EncodeToString(b))
}

// GenerateActionID creates a unique action_id with "act" prefix.
func GenerateActionID() string {
	return GenerateID("act")
}

// GenerateBlockID creates a unique block_id with "blk" prefix.
func GenerateBlockID() string {
	return GenerateID("blk")
}
