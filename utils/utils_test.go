package utils

import (
	"strings"
	"testing"
)

func TestInt32Ptr(t *testing.T) {
	// Test case 1: Positive integer
	input := int32(42)
	expected := &input

	result := Int32Ptr(input)
	if *result != *expected {
		t.Errorf("Expected %v, but got %v", input, result)
	}

	// Test case 2: Negative integer
	input = int32(-10)
	expected = &input

	result = Int32Ptr(input)
	if *result != *expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestUniqName(t *testing.T) {
	// Test case 1: Check if UniqName returns a non-empty string
	result := UniqName("mybase")
	if result == "" {
		t.Errorf("UniqName should not return an empty string")
	}

	// Test case 2: Check if UniqName handles base names with slashes
	result = UniqName("my/base/name")
	expected := "base-qd-"

	if !strings.HasPrefix(result, expected) {
		t.Errorf("UniqName should handle base names with slashes correctly")
	}

	// test case 3: length should be no longer than 27chars: "%s-qd-%s", shortBase (15), uniqueId(8), -qd- (4)
	result = UniqName("foo-bar-baz-brick-house.com/super-duper-crazy-long-image-name-goes-here:abcdefghijklmnop123")
	if len(result) > 27 {
		t.Errorf("Expected unique deployment name is too long")
	}
}

func TestShortName(t *testing.T) {
	// test case 1: length should be no longer than 15 chars
	result := shortName("foo-bar-baz-brick-house.com/super-duper-crazy-long-image-name-goes-here:abcdefghijklmnop123")
	if len(result) > 15 {
		t.Errorf("Expected short base name is too long")
	}
}
