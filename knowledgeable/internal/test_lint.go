package internal

import "fmt"

func TestLintFail() {
	fmt.Errorf("this error is ignored") // errcheck vil fejle her
}
