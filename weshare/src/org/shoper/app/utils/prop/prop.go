package prop

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
)

var (
	COMMENT   byte = '#' //注释符号
	SEPARATOR byte = '=' //key-value分隔符
)

type Properties struct {
	datas map[string]string
}

//Load specify file
func Load(filePath string) (p *Properties, e error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	//用于确定读取到第几行内容出错，出错提示
	lineNumber := 0
	p = new(Properties)
	p.datas = make(map[string]string)
	for {
		//每次读一行内容
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		lineNumber++
		//trim
		line = bytes.TrimSpace(line)
		//判断是否读取到空白行,读取到空白行那么继续下次循环
		if bytes.Equal(line, []byte("")) {
			continue
		}
		//检查是否该行存在#注释语句
		commentIndex := bytes.IndexByte(line, COMMENT)
		separatorIndex := bytes.IndexByte(line, SEPARATOR)
		//如果index为0，那么该行应该是注释语句，直接跳过
		propIndex := 0
		if commentIndex == 0 {
			continue
		} else if (commentIndex == -1 && separatorIndex <= 0) || commentIndex-1 == separatorIndex {
			//如果#下标为-1并且＝下标小于0说明该错误，因为不存在 =sd 这种无key有value的语法，如果注释在＝之后也就是 s=#
			return p, errors.New("解析错误,LineNumber:" + strconv.Itoa(lineNumber))
		} else {
			//正常情况下
			if commentIndex-1 > separatorIndex {
				//有注释的情况，注释代码至少要在separatorIndex+1的后面把注释内容去掉
				propIndex = commentIndex
			}
		}

		if propIndex == 0 {
			line = line[0:]
		} else {
			line = line[0:propIndex]
		}
		//这里的newLine得到是每行去掉注释的内容,现在来分离=两边的内容
		key := line[0:separatorIndex]
		value := line[separatorIndex+1:]
		p.put(string(key[0:]), string(value[0:]))
	}
	return p, nil
}
func (p *Properties) put(key string, value string) {
	p.datas[key] = value
}
func (p *Properties) GetData(str string) string {
	return p.datas[str]
}
