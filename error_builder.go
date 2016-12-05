package common

import (
	"fmt"
	"strings"
)

type ErrorBuilder []string

func (ec *ErrorBuilder) Append(e error) { *ec = append(*ec, e.Error()) }

func (ec *ErrorBuilder) ToError() error {
	var ret error
	if len(*ec) > 0 {
		ret = fmt.Errorf(strings.Join(*ec, "\n"))
	}

	return ret
}
