// filesystem_user_test
package gfs

import (
	"fmt"
	"testing"
)

func TestUMask(t *testing.T) {
	fm := FileMask("077")
	fmt.Println(fm.GetAfterMasked())
}
