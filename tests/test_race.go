package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {

	//req()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go req()
	}

	wg.Wait()
}

func req() {
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Int63n(999999999-111111111) + 111111111
	ref := fmt.Sprintf("09%d", rnd)
	var jsonData = []byte("")

	fmt.Printf("ref id: %s -> time %s -> timestamp %d\n", ref, time.Now().Format("2006-01-02 15:04:05.000Z"), time.Now().UnixNano())

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:8882/user/%s/credit/code/abc", ref), bytes.NewBuffer(jsonData))
	req.Close = true
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("ref id %d: %s \n", ref, string(body))

	fmt.Println("-----------------")

	wg.Done()
}
