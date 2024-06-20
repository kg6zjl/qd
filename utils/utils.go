package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func Int32Ptr(i int32) *int32 { return &i }

func UniqName(base string) string {
	// generate a unique name
	rand.NewSource(time.Now().UnixNano())
	randInt := rand.Int()
	uniqueId := strconv.Itoa(randInt)

	// verify the unique bit is at least 12 chars or use the whatever we've got
	if len(uniqueId) > 8 {
		uniqueId = uniqueId[:8]
	}

	// format the name
	uniqueName := fmt.Sprintf("%s-qd-%s", base, uniqueId)

	return uniqueName
}
