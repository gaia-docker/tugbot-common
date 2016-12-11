package common

import (
	"github.com/stretchr/testify/assert"

	"errors"
	"testing"
)

func TestErrorBuilder(t *testing.T) {
	var ec ErrorBuilder
	ec.Append(errors.New("uno"))
	ec.Append(errors.New("dos"))
	assert.Equal(t, "uno\ndos", ec.ToError().Error())
}

func TestErrorBuilderNoErrors(t *testing.T) {
	var ec ErrorBuilder
	assert.Nil(t, ec.ToError())
}
