package helper

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func RandomNumber(min, max int) int {

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min

}

func RandomFloat(min, max float64) float64 {

	rand.Seed(time.Now().UnixNano())

	num := rand.Float64()*(max-min) + min

	return math.Floor(num*100) / 100

}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func CheckSlice(arr []int, id int) bool {
	for _, item := range arr {
		if item == id {
			return true
		}
	}
	return false
}
