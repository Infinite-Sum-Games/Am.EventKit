package pkg

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Totally over-engineered OTP generation function :(
func GenerateOTP() (string, []string, error) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(900000) + 100000
	strNum := strconv.Itoa(num)
	strSlice := strings.Split(strNum, "")
	fmt.Println(strSlice)
	return strNum, strSlice, nil
}
