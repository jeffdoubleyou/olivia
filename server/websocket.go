package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/jeffdoubleyou/olivia/locales"

	"github.com/jeffdoubleyou/olivia/modules/start"

	"github.com/gookit/color"
	"github.com/gorilla/websocket"
	"github.com/jeffdoubleyou/olivia/analysis"
	"github.com/jeffdoubleyou/olivia/user"
	"github.com/jeffdoubleyou/olivia/util"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RequestMessage is the structure that uses entry connections to chat with the websocket
type RequestMessage struct {
	Type        int              `json:"type"` // 0 for handshakes and 1 for messages
	Content     string           `json:"content"`
	Token       string           `json:"user_token"`
	Locale      string           `json:"locale"`
	Language    string           `json:"language"`
	Information user.Information `json:"information"`
	Context     string           `json:"context"`
}

// ResponseMessage is the structure used to reply to the user through the websocket
type ResponseMessage struct {
	Content     string                 `json:"content"`
	Tag         string                 `json:"tag"`
	Information user.Information       `json:"information"`
	Data        map[string]interface{} `json:"data"`
}

// SocketHandle manages the entry connections and reply with the neural network
func SocketHandle(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	fmt.Println(color.FgGreen.Render("A new connection has been opened"))

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Unmarshal the json content of the message
		var request RequestMessage
		if err = json.Unmarshal(msg, &request); err != nil {
			continue
		}

		// Set the information from the client into the cache
		if reflect.DeepEqual(user.GetUserInformation(request.Token), user.Information{}) {
			user.SetUserInformation(request.Token, request.Information)
		}

		j, _ := json.MarshalIndent(request, "", "\t")
		fmt.Printf("REQUEST: %s\n", j)
		//request.Locale = "1212"
		// If the type of requests is a handshake then execute the start modules
		if request.Type == 0 {
			start.ExecuteModules(request.Token, request.Locale)

			message := start.GetMessage()
			if message != "" {
				// Generate the response to send to the user
				response := ResponseMessage{
					Content:     message,
					Tag:         "start module",
					Information: user.GetUserInformation(request.Token),
				}

				bytes, err := json.Marshal(response)
				if err != nil {
					panic(err)
				}

				if err = conn.WriteMessage(msgType, bytes); err != nil {
					continue
				}
			}

			continue
		}

		// Write message back to browser
		response := Reply(request)
		r, _ := json.MarshalIndent(response, "", "\t")
		fmt.Printf("RESPONSE: %s\n", r)

		if err = conn.WriteMessage(msgType, response); err != nil {
			continue
		}
	}
}

// Reply takes the entry message and returns an array of bytes for the answer
func Reply(request RequestMessage) []byte {
	var responseSentence, responseTag string
	var intent *analysis.Intent
	// Send a message from res/datasets/messages.json if it is too long
	if len(request.Content) > 500 {
		responseTag = "too long"
		responseSentence = util.GetMessage(request.Locale, responseTag)
	} else {
		// If the given language is not supported yet, set english
		locale := request.Locale
		language := request.Language
		if !locales.Exists(language) {
			language = locale
		}

		if request.Context != "" {
			responseTag, responseSentence, intent = analysis.NewSentence(
				language, request.Content,
			).Calculate(*cache, neuralNetworks[locale], request.Token, request.Context)
		} else {
			responseTag, responseSentence, intent = analysis.NewSentence(
				language, request.Content,
			).Calculate(*cache, neuralNetworks[locale], request.Token)
		}
	}

	// Marshall the response in json
	response := ResponseMessage{
		Content:     responseSentence,
		Tag:         responseTag,
		Information: user.GetUserInformation(request.Token),
	}

	if intent != nil {
		response.Data = intent.Data
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Res: %s\n", bytes)
	return bytes
}
