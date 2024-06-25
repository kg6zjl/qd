package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Int32Ptr(i int32) *int32 { return &i }

func UniqName(base string) string {
	// generate a unique name
	rand.NewSource(time.Now().UnixNano())
	randInt := rand.Int()
	uniqueId := strconv.Itoa(randInt)

	// get short name (63 chars and limited special chars)
	shortBase := shortName(base)

	// verify the unique bit is at least 12 chars or use the whatever we've got
	if len(uniqueId) > 8 {
		uniqueId = uniqueId[:8]
	}

	// format the name
	uniqueName := fmt.Sprintf("%s-qd-%s", shortBase, uniqueId)

	return uniqueName
}

func shortName(base string) string {
	var shortStr string
	// Check if "/" is in the string
	if strings.Contains(base, "/") {
		// Split the string on "/"
		splitStr := strings.Split(base, "/")
		shortStr = splitStr[1]
	} else {
		shortStr = base
	}

	// Santize down to allowed chars
	reg, _ := regexp.Compile("[^a-zA-Z0-9-_.]+")
	shortStr = reg.ReplaceAllString(shortStr, "")

	// Shorten it to 15 characters
	if len(shortStr) > 15 {
		shortStr = shortStr[:15]
	}

	return shortStr
}
