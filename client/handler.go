package client

import (
	"fmt"
	"net/http"

	"github.com/secondarykey/ikascrew"
	//"github.com/secondarykey/ikascrew/video"
)

const (
	WEBDIR = "/.public"
)

func init() {
}

func (ika *IkascrewClient) startHTTP() {

	http.HandleFunc("/load", ika.loadHandler)
	http.HandleFunc("/switch", ika.switchHandler)
	http.Handle("/", http.FileServer(http.Dir(ikascrew.ProjectName()+"/.public/")))

	go func() {
		http.ListenAndServe(":5555", nil)
	}()
}

// load
func (ika *IkascrewClient) loadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	//name := r.FormValue("name")
	//t := r.FormValue("type")

	//v, err := video.Get(video.Type(t), name)
	//if err != nil {
	//fmt.Println(err)
	//return
	//}

	//TODO 変更用の処理を作成？
	//ika.window.SetEffect(e)
}

// switch
func (ika *IkascrewClient) switchHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	t := r.FormValue("type")

	err := ika.callSwitch(name, t)
	if err != nil {
		fmt.Println(err)
	}
}
