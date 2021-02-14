package slack

import (
	"aws-example/config"
	errors "aws-example/error"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type messageRequest struct {
	Channel string `json:"channel"`
	Message string `json:"text"`
}

type messageResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type Client interface {
	SendMessage(message string) error
}

type ClientImpl struct {
	Client http.Client
}

func NewClient(client http.Client) Client {
	return ClientImpl{Client: client}
}

func (c ClientImpl) SendMessage(message string) error {
	requestBody := messageRequest{
		Message: message,
		Channel: config.SLACK_CHANNEL,
	}
	data, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	client := c.Client
	request, err := http.NewRequest("POST", config.SLACK_SEND_MESSAGE_URL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	enrichRequest(request)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	response, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		return errors.InternalError{Message: er.Error()}
	}
	var responseBody messageResponse
	er = json.Unmarshal(response, &responseBody)
	if er != nil {
		return errors.InternalError{Message: er.Error()}
	}
	if !responseBody.Ok {
		return errors.InternalError{Message: responseBody.Error}
	}
	return nil
}

func enrichRequest(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	var bearer = "Bearer " + config.SLACK_TOKEN
	request.Header.Set("Authorization", bearer)
}
