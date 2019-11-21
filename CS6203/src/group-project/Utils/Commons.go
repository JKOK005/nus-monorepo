package Utils

import (
	"hash/fnv"
	"os"
	"strconv"
	"strings"
)

func GetEnvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	valueInt, _ := strconv.Atoi(value)
	return valueInt
}

func GetEnvStrSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return strings.Split(value, ",")
}

func GetHashFunction() func(string) uint32 {
	/*
		Returns a hash function with implementation as defined
	*/
	return func(hashString string) uint32 {
		/*
			Applies a hash function over a key
			Returns the hashed value of the key to represent the hash group it should be assigned to
		*/
		h := fnv.New32a()
		_,_ = h.Write([]byte(hashString))
		return h.Sum32() % 3 +1
	}
}