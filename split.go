package main

import (
    "fmt"
    "strings"
)

func main() {
    s := strings.Split("127.0.0.1|{'mail': 'gdumoulin@me.com', 'newsletter': 'true'}", "|")
	ip:= s[0]
	fmt.Println(ip, "ok")
	if (len(s) > 1) {
		fmt.Println(s[1], "json")
	}
    
}