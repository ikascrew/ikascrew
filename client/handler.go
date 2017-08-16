package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ikascrew/ikascrew"
)

const (
	WEBDIR = "/.public"
)

func init() {
}

func (ika *IkascrewClient) startHTTP() {

	http.HandleFunc("/switch", ika.switchHandler)
	http.Handle("/", http.FileServer(http.Dir(ikascrew.ProjectName()+"/.public/")))

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (ika *IkascrewClient) switchHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	t := r.FormValue("type")

	err := ika.callEffect(name, t)
	if err != nil {
		fmt.Println(err)
	}
}
