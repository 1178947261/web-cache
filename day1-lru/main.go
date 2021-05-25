package main

import (
	"fmt"
	"strconv"
	"strings"
)

type People interface {
	Speak(string) string
}

type Student struct{}

func (stu *Student) Speak(think string) (talk string) {
	if think == "sb" {
		talk = "你是个大帅比"
	} else {
		talk = "您好"
	}
	return
}

func main() {
	var peo People = &Student{}
	think := "bitch"
	fmt.Println(peo.Speak(think))

	s := "hello world hello world hello world hello world hello worldhello worldhello worldhello worldhello worldhello worldhello world"
	//str := "wo"

	//以str为分隔符，将s切分成多个子串，结果中**不包含**str本身。如果str为空则将s切分成Unicode字符列表。
	//如果s中没有str子串，则将整个s作为[]string的第一个元素返回。
	//参数n表示最多切分出几个子串，超出的部分将不再切分，最后一个n包含了所有剩下的不切分。
	//如果n为0，则返回nil；如果n小于0，则不限制切分个数，全部切分
	index := strings.SplitN(s, " ", 2)
	fmt.Println(index) //2

	fmt.Println(strconv.Itoa(1) + "2")
	println(int(1))
}
