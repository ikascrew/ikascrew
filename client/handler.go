package client

import (
	"fmt"
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

	go func() {
		http.ListenAndServe(":5555", nil)
	}()
}

func (ika *IkascrewClient) switchHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	t := r.FormValue("type")

	err := ika.callSwitch(name, t)
	if err != nil {
		fmt.Println(err)
	}
}
