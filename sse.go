package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

const (
	DEBUG bool = true
)

//Following is all the SSE implemntation ... need to break out into own file
type RequestModel struct {
	lastEventId string
	topic       string
	message     string
}

type SSE struct {
	parentApp      *App
	channels       map[string]*Channel
	incomingClient chan *RequestModel
	closingClient  chan *Client
	config         *ini.File
}

func createSSE(a *App) (sse *SSE) {
	sse = &SSE{
		parentApp:      a,
		channels:       make(map[string]*Channel),
		incomingClient: make(chan *RequestModel),
		closingClient:  make(chan *Client),
		config:         a.cfg,
	}

	go sse.dispatch()

	return
}

func (sse *SSE) ServeSSE(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Flusher is not supported", http.StatusInternalServerError)
		return
	}

	topic := req.URL.Path[len("/events/"):]

	// In every case we should receive topic
	//topic is anything after "/events/" prefix
	if topic == "" {
		http.Error(rw, "Topic is missing", http.StatusBadRequest)
		return
	}
	_, error := sse.config.GetSection(topic)
	if error != nil {
		http.Error(rw, "Topic not available", http.StatusBadRequest)
		return
	}

	//sets up content for a topic
	if req.Method == "POST" {
		b, _ := ioutil.ReadAll(req.Body)
		// update channel in db with request body as value; clients will query via api on instanciation
		// use an upsert so that the valid channel record is created if not there.  Making ini entry becomes method of creation
		var c channeldef
		c.name.String = topic
		c.name.Valid = true
		c.value.String = string(b)
		c.value.Valid = true

		//c.name = null.NewString(topic, topic != "")

		// call router
		if err := c.upsertChannel(sse.parentApp.DB); err != nil {
			if DEBUG {
				fmt.Println("Channel upsert failed on %s with error %s", c.name.String, err)
			}
		}

		// broadcast value
		go sse.addToChannel(&RequestModel{
			topic:   topic,
			message: string(b), //ingest all of the body as the message; dispatcher places all of message as data
		})
		rw.WriteHeader(204) // got it, no content in reply.  Alter for Pentaho integration

	} else if req.Method == "GET" {
		// if topic in list of ini channels
		if len(sse.config.Section(topic).Name()) > 1 {

			// Since we can listen non existing channel, we have to create new one
			channel := sse.createChannelIfNotExist(topic)
			// Every client channel receives messages
			clientChan := make(chan *Message)
			client := CreateClient(topic, clientChan)
			channel.clients[clientChan] = true

			if DEBUG {
				fmt.Printf("New client has been registered in %s channel. Total: %d\n\n", channel.name, len(channel.clients))
			}

			rw.Header().Set("Content-Type", "text/event-stream")
			rw.Header().Set("Cache-Control", "no-cache")
			rw.Header().Set("Connection", "keep-alive")
			rw.Header().Set("Access-Control-Allow-Origin", "*")

			close := rw.(http.CloseNotifier).CloseNotify()

			// After we started SSE connection between server and client
			// We have to set timeout for client
			seconds := sse.config.Section(topic).Key("expires").MustInt(0)
			//initialize to 4 hours in case seconds is 0
			evtime := time.NewTimer(time.Hour * 4)
			//set expiration based on time allotted
			if seconds > 0 {
				evtime = time.NewTimer(time.Second * time.Duration(seconds)) // TIMEOUT)
			}
			//get switch
			for {
				select {
				case <-evtime.C: // waits here
					fmt.Fprintf(rw, "event: %s\n", "timeout")
					fmt.Fprintf(rw, "data: %ds\n\n", seconds)
					flusher.Flush()
					sse.closingClient <- client
					return
				case <-close:
					sse.closingClient <- client
					return
				case msg := <-clientChan:
					if msg == nil {
						return
					}

					//This is the standard format for SSE.  JSON can be placed in "data: " position for ingestion by receivers.
					fmt.Fprintf(rw, "id: %d\n", msg.id)
					//	fmt.Fprintf(rw, "event: %s\n", msg.event) //Don't set event when using EventSource to catch messages with onmessage handler.
					fmt.Fprintf(rw, "data: %s\n\n", strings.Replace(string(msg.data), "\n", "\ndata: ", -1))
					flusher.Flush()
				}
			}
		} else {
			fmt.Fprintf(rw, "No such channel available")
		}
	}
}

func (sse *SSE) addToChannel(m *RequestModel) {
	sse.createChannelIfNotExist(m.topic)
	sse.incomingClient <- m
}

func (sse *SSE) createChannelIfNotExist(chanName string) *Channel {
	if !sse.doesChannelExist(chanName) {
		sse.channels[chanName] = CreateChannel(chanName)
	}
	return sse.channels[chanName]
}

func (sse *SSE) dispatch() {
	for {
		select {
		case s := <-sse.incomingClient:
			message := NewMessage("msg", s.message, s.topic) //set event, data and channel respectively
			go addMessage(sse.channels[s.topic], message)
		case c := <-sse.closingClient:
			c.Close(sse.channels[c.Channel()].clients)
			if DEBUG {
				fmt.Printf("Removed client from %s channel. Clients left: %d\n\n", c.Channel(), sse.channels[c.Channel()].ClientsCount())
			}
		}
	}
}

func addMessage(channel *Channel, message *Message) {
	// Add message for every client in this channel
	for client := range channel.clients {
		client <- message
	}
}

func (s *SSE) doesChannelExist(name string) bool {
	_, ok := s.channels[name]
	return ok
}
