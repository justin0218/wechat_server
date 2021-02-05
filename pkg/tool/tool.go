package tool

import (
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func IntInArray(item int, arr []int) (dx int) {
	for idx, v := range arr {
		if v == item {
			dx = idx
			return
		}
	}
	return -1
}

func StringInArray(item string, arr []string) (dx int) {
	for idx, v := range arr {
		if v == item {
			dx = idx
			return
		}
	}
	return -1
}

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func Tracefile(content string) (err error) {
	fd, _ := os.OpenFile("aun.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd_content := content + "\n"
	buf := []byte(fd_content)
	_, err = fd.Write(buf)
	if err != nil {
		return
	}
	err = fd.Close()
	if err != nil {
		return
	}
	return
}
