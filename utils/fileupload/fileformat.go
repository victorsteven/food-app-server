package fileupload

import (
	"github.com/twinj/uuid"
	"path"
)

func FormatFile(fn string) string {

	ext := path.Ext(fn)
	u := uuid.NewV4()

	newFileName := u.String() + ext

	return newFileName
}
