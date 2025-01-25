package env

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val

}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToInt(val, fallback)
}
func StringToInt(val string, fallback int) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return res
}
func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToBool(val, fallback)
}
func StringToBool(val string, fallback bool) bool {
	res, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return res
}
func GetStringSlice(key string, fallback []string) []string {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToStringSlice(val, fallback)
}
func StringToStringSlice(val string, fallback []string) []string {
	res := strings.Split(val, ",")
	if len(res) == 0 {
		return fallback
	}
	return res
}
func GetIntSlice(key string, fallback []int) []int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToIntSlice(val, fallback)
}
func StringToIntSlice(val string, fallback []int) []int {
	res := strings.Split(val, ",")
	if len(res) == 0 {
		return fallback
	}
	intSlice := make([]int, len(res))
	for i, v := range res {
		intSlice[i] = StringToInt(v, 0)
	}
	return intSlice
}
func GetBoolSlice(key string, fallback []bool) []bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToBoolSlice(val, fallback)
}
func StringToBoolSlice(val string, fallback []bool) []bool {
	res := strings.Split(val, ",")
	if len(res) == 0 {
		return fallback
	}
	boolSlice := make([]bool, len(res))
	for i, v := range res {
		boolSlice[i] = StringToBool(v, false)
	}
	return boolSlice
}
func GetFloat64(key string, fallback float64) float64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToFloat64(val, fallback)
}
func StringToFloat64(val string, fallback float64) float64 {
	res, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fallback
	}
	return res
}
func GetFloat64Slice(key string, fallback []float64) []float64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToFloat64Slice(val, fallback)
}
func StringToFloat64Slice(val string, fallback []float64) []float64 {
	res := strings.Split(val, ",")
	if len(res) == 0 {
		return fallback
	}
	float64Slice := make([]float64, len(res))
	for i, v := range res {
		float64Slice[i] = StringToFloat64(v, 0)
	}
	return float64Slice
}
func GetStringMap(key string, fallback map[string]interface{}) map[string]interface{} {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return StringToStringMap(val, fallback)
}
func StringToStringMap(val string, fallback map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	err := json.Unmarshal([]byte(val), &res)
	if err != nil {
		return fallback
	}
	return res
}
