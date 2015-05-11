package console

import (
	"bufio"
	"fmt"
	"os"
)

func Print(v ...interface{}) {
	fmt.Print(v)
}

func Println(v ...interface{}) {
	fmt.Println(v)
}

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v)
}

//检查输入信号量
func CheckInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("please entry q or Q to quit the progame!")
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "q" || command == "Q" {
			break
		}

	}

}
