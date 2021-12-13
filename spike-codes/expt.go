package spikecodes

import (
	"encoding/json"
	"fmt"

	// "os"
	"time"
)

type TimeStruct struct {
	Sometime time.Time
}

func Drr() {
	// fmt.Println(os.Getwd())
	s := fmt.Sprint(time.Now())

	ss := []byte(s)

	var t time.Time
	err := json.Unmarshal(ss, &t)
	fmt.Println(err)
	fmt.Println(t)

}
