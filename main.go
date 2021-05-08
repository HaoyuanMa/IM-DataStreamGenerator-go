package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type loginRequest struct {
	Username	string `json:"username"`
	Password	string `json:"password"`
}

type loginResponse struct {
	Message string `json:"message"`
	Status int32 `json:"status"`
	Token string `json:"token"`
}

type ClientCall struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

func main() {
	token := ""
	for{
		if token != ""{
			break
		}
		token = Login()
	}
	ws,err := websocket.Dial("ws://localhost:5202/Socket/BuildConnection?token=" + token,"","http://localhost")
	if err != nil {
		fmt.Println(err)
	}

	oprator := ""
	for{
		if oprator == "n"{
			break
		}
		fmt.Println("Please input count of data")
		var count int
		fmt.Scanf("%d",&count)
		fmt.Println("Please input speed of data generating(ms)")
		var delay int
		fmt.Scanf("%d",&delay)
		for i := 0; i < count; i++ {
			call := ClientCall{
				Method:"UploadStream",
				Params: map[string]string{
					"content":strconv.Itoa(i),
				},
			}
			data,_ := json.Marshal(call)
			fmt.Println("send : " + strconv.Itoa(i))
			_, err := ws.Write(data)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Duration(int64(time.Millisecond) * int64(delay)))
		}

		fmt.Println("Continue to send(y/n)")
		fmt.Scanf("%s",&oprator)



	}
}

func Login() string {
	username := ""
	password := ""
	fmt.Println("Please input username")
	fmt.Scanf("%s",&username)
	fmt.Println("Please input password")
	fmt.Scanf("%s",&password)
	client := http.Client{}
	data := loginRequest{
		Username: username,
		Password: password,
	}
	jsonStr,_ := json.Marshal(data)
	response,err := client.Post("http://localhost:5202/Account/Login","application/json",bytes.NewBuffer(jsonStr))
	if err!=nil{
		return ""
	}
	result,_ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(result))
	var res loginResponse
	if err := json.Unmarshal(result, &res); err != nil {
		return ""
	}
	return res.Token
}
