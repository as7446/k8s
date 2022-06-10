### 文件读取操作
* 按行读取
```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("example.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
		fmt.Println(string(buf))
	}
}

```
* 带缓冲区的方式读取到终端
```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("example.txt")
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Print(str)
	}
}

```

* 读取文件到显示终端（一次性读取数据到内存中）,不是适合操作大文件。
```go
package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("go.mod")
	if err != nil {
		fmt.Errorf("读取文件错误：%v", err)
		return
	}
	fmt.Println(string(data))
}

```
