package common

import (
	"fmt"
	"github.com/satori/uuid"
)

func UUID() string {
	if tmp, err := uuid.NewV4(); err == nil {
		return fmt.Sprintf("%s", tmp)
	} else {
		panic("生成UUID时发生错误")
	}
}
