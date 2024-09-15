package atomic_update

import (
	"fmt"
	"os"
	"time"
)

func Save(path string, data []byte) error {
	tmp := fmt.Sprintf("%s.tmp.%d", path, time.Now().UnixNano())
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return err
	}
	defer func() {
		fp.Close()
		if err != nil {
			os.Remove(tmp)
		}
	}()

	_, err = fp.Write(data)
	if err != nil {
		return err
	}
	err = fp.Sync()
	if err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
