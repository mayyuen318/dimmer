package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/tarm/goserial"
)

type Page struct {
	Error string
}

var OpMap = map[string][]byte{
	"On":  []byte("0"),
	"Off": []byte("1"),
	"Inc": []byte("+"),
	"Dec": []byte("-"),
}

var port = flag.Int("port", 8080, "Port number of web server")
var comport = flag.String("serial-port", "", "Serial port to Zigbee")

func dimmerControlHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{}
	op := r.FormValue("op")

	if op != "" {
		if OpMap[op] != nil {
			c := &serial.Config{Name: *comport, Baud: 9600}
			s, err := serial.OpenPort(c)
			if err != nil {
				p.Error = err.Error()
			} else {
				_, err := s.Write(OpMap[op])
				if err != nil {
					p.Error = err.Error()
				}
				p.Error = "OK"
			}
		} else {
			p.Error = "Unknown Op"
			fmt.Printf("OP: %s", op)
		}
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func main() {
	flag.Parse()

	if *comport == "" {
		fmt.Printf("Serial Port needed\n")
		return
	}

	http.HandleFunc("/", dimmerControlHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
