package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
)

type Orders struct {
	Stream   string `json:"stream"`
	Data   struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	} `json:"data"`
}

var (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url := "wss://stream.binance.com:9443/stream?streams=ethbtc@depth5"
	log.Printf("connecting to %s", url)

	c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/stream?streams=btcusdt@depth20", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			//var orders Orders
			//json.Unmarshal(message, &orders)
			//fmt.Printf("|%-9.2f|%-9.2f|\n", Green + "Price" + Reset, Green + "Amount" + Reset)
			//for _, v := range orders.Data.Asks {
			//	fmt.Printf("|%-9.2f|%-9.2f|\n", v[0], v[1])
			//}
			//fmt.Printf("|%-9.2f|%-9.2f|\n", Red + "Price" + Reset, Green + "Amount" + Reset)
			//for _, v := range orders.Data.Bids {
			//	fmt.Printf("|%-9.2f|%-9.2f|\n", v[0], v[1])
			//}
			//fmt.Println("*********************")
			//fmt.Println("----------------------")
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}