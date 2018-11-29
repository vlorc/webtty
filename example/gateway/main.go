package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
)

func forward(table sync.Map, message []byte) {
	iter := jsoniter.ConfigFastest.BorrowIterator(message)
	defer jsoniter.ConfigFastest.ReturnIterator(iter)

	var dest string
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, field string) bool {
		if "dest" == field {
			dest = iter.ReadString()
			return false
		}
		iter.Skip()
		return true
	})

	it, ok := table.Load(dest)
	if ok {
		log.Printf("forward: [%s]: %s\n", dest, message)
		it.(*websocket.Conn).WriteMessage(websocket.TextMessage, message)
	}
}

func main() {
	var table sync.Map
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	addr := flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()
	log.SetFlags(0)
	log.Fatal(http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[1:]
		if _, ok := table.Load(id); ok {
			w.WriteHeader(403)
			log.Printf("exist [%s]\n", id)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer conn.Close()

		table.Store(id, conn)
		log.Printf("join [%s]\n", id)
		defer table.Delete(id)
		defer log.Printf("leave [%s]\n", id)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read: [%s]: %s\n", id, err.Error())
				break
			}
			forward(table, message)
		}
	})))
}
