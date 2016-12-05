package common

import (
	"fmt"
	"strings"
)

type ErrorBuilder []string

func (ec *ErrorBuilder) Append(e error) { *ec = append(*ec, e.Error()) }

func (ec *ErrorBuilder) ToError() error { return fmt.Errorf(strings.Join(*ec, "\n")) }
