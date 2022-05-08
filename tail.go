package tail

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const BUFSIZE = 64

func SeekLineNFromEnd(file *os.File, n_lines uint) (err error) {
	if err != nil {
		return err
	}
	if n_lines == 0 {
		file.Close()
		return errors.New("TailNLines: n_lines is 0, no read required")
	}

	buf := make([]byte, BUFSIZE)

	line_end := byte('\n')

	stat, err := os.Stat(file.Name())
	if err != nil {
		file.Close()
		return err
	}
	var start_pos int64
	pos := stat.Size() - BUFSIZE

	for {

		if pos <= 0 {
			start_pos = 0
			break
		}

		if n_lines == 0 {
			break
		}

		n, err := file.ReadAt(buf, pos)

		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		nl := bytes.LastIndexByte(buf, line_end)
		if nl == -1 {
			start_pos = pos
			pos = pos - BUFSIZE
		} else {
			start_pos = pos + int64(nl) + 1
			pos = pos + int64(nl) - BUFSIZE
			n_lines -= 1
		}

	}

	_, err = file.Seek(start_pos, 0)
	if err != nil {
		return err
	}

	return
}

// Tail last n_lines of file with path
func TailFile(path string, n_lines uint) (content string, err error) {
	file, err := os.Open(path)
	content = ""
	if err != nil {
		return content, err
	}

	err = SeekLineNFromEnd(file, n_lines)
	if err != nil {
		return content, nil
	}

	data := make([]byte, 1000)
	for {
		n, err := file.Read(data)
		if err != nil {
			break
		}
		content += string(data[:n])
	}

	return content, nil
}
