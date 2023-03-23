package lark

import (
	"Surfing/util"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendMessageByGroupBot(text string, webhookURL string) {

	method := "POST"
	payload := strings.NewReader(`{"msg_type": "text","content": {"text": "` + text + `"}}`)
	//text为你要发送的信息

	client := &http.Client{}
	req, err := http.NewRequest(method, webhookURL, payload)
	if err != nil {
		util.PrintLog(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		util.PrintLog(err.Error())
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		util.PrintLog(err.Error())
		return
	}

	util.PrintLog(string(body))

}
