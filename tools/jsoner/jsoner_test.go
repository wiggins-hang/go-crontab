package jsoner

import (
	"fmt"
	"testing"
)

type User struct {
	Id       uint64  `json:"id"`
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Money    float64 `json:"money"`
	MoneyStr string  `json:"moneyStr"`
}

func Test_IntString(t *testing.T) {
	userStr := `{
		"id": "1",
		"age": 99,
		"money": "123",
		"moneyStr": 100000.898
	}`
	userDemo := &User{}
	err := UnmarshalByte([]byte(userStr), userDemo)
	fmt.Println("err", err)
	fmt.Printf("%+v", userDemo)
}
