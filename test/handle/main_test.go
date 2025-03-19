package testhanle

import "testing"

func AddTwoNumber(a int, b int) int {
	return a + b
}

func TestAddTwoNummber(t *testing.T){
	var (
		a = 1
		b = 2
		output = 10
	)
	result:= AddTwoNumber(a,b)
	if result != output {
		t.Errorf("AddTwoNumber(%d,%d) = %d; want %d", a, b, result, output)
	}
}
 
// B1: cd vào thư mục chứa file main_test.go
// B2: chạy lệnh: go test or go test -v
// B3: Kết quả trả về: