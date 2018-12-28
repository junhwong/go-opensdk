package opensdk

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	rnd := strconv.FormatInt(rand.Int63(), n)
	if len(rnd) > n {
		rnd = rnd[:n]
	}
	return fmt.Sprintf("%s", rnd)
}
