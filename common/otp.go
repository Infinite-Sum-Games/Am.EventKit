package common

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateOTP() (string, []string, error) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	num := r.Intn(900000) + 100000
	strNum := strconv.Itoa(num)
	strSlice := strings.Split(strNum, "")

	return strNum, strSlice, nil
}
