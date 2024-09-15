package inplace

import "os"

func Save(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer func() {
		fp.Close()
		if err != nil {
			os.Remove(path)
		}
	}()

	_, err = fp.Write(data)
	if err != nil {
		return err
	}
	return nil
}
