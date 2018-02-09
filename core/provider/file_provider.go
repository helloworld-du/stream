package provider

import (
	"sync"
	"os"
	"errors"
	"bufio"
	"io"
	"strings"
)

const BUFFER_SIZE = 50

type FileProvider struct {
	LineSeparator byte
	filePath	string
	file		*os.File
	lock	*sync.Mutex
	buf		chan string
	lastErr error
}

func NewFileProvider(filePath  string, lineSeparator byte) (*FileProvider) {
	return &FileProvider{lineSeparator, filePath, nil, &sync.Mutex{}, nil, nil}
}


func (fp *FileProvider) init() (err error) {

	if len(fp.filePath) == 0 {
		return errors.New("[FileProvider]no file path")
	}
	fp.lock.Lock()
	defer fp.lock.Unlock()
	if fp.file != nil {
		return nil
	}
	fp.file, err = os.Open(fp.filePath);
	if  err != nil {
		return err
	}
	fp.file.Seek(0, 0)
	fp.buf = make(chan string, BUFFER_SIZE)

	go func() {
		br := bufio.NewReader(fp.file)
		ls := string(fp.LineSeparator)
		for{
			//每次读取一行
			line, err := br.ReadString(fp.LineSeparator)
			line = strings.TrimRight(line, ls)

			if len(line) > 0 || (len(line) == 0 && err == nil) {
				//有一行数据，不管是否有错
				//或者这一行是空 即line长度是0而没有报错
				fp.buf <- line
			}
			if err != nil {
				//有报错，推出
				fp.lastErr = err
				break
			}
		}
		//文件读取完成了就关闭了chan，数据存在buffer中
		close(fp.buf)
	}()
	return nil


}


func (fp *FileProvider) Read() (msg interface{}, err error, hasNext bool) {
	if err = fp.init(); err != nil {
		return "", err, false
	}
	msg, hasNext = <- fp.buf
	if !hasNext {
		//数据如读取完成了，就关闭文件了
		fp.stop()
		if fp.lastErr != io.EOF {
			err = fp.lastErr
		}
	}
	return
}

func (fp *FileProvider)stop() {
	defer func() {
		fp.lock.Unlock()
		recover()
	}()
	fp.lock.Lock()
	fp.file.Close()
	fp.file = nil
}
