package console

import (
	"bufio"
	"fmt"
	"os"
	"time"
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
	fmt.Println("please entry q or Q to quit the progame!")
	reader := bufio.NewReader(os.Stdin)
	for {

		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "q" || command == "Q" {
			break
		}
		fmt.Println("please entry q or Q to quit the progame!")
		time.Sleep(1 * time.Second)
	}
}
