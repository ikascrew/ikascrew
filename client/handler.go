package client

import (
	"fmt"
	"net/http"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/effect"
	"github.com/secondarykey/ikascrew/video"
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

	name := r.FormValue("name")
	t := r.FormValue("type")

	v, err := video.Get(video.Type(t), name)
	if err != nil {
		fmt.Println(err)
		return
	}

	e, err := effect.NewNormal(v)
	if err != nil {
		fmt.Println(err)
		return
	}

	wk := ika.window.GetEffect()
	ika.window.SetEffect(e)

	wk.Release()

}

// switch
func (ika *IkascrewClient) switchHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	t := r.FormValue("type")
	e := r.FormValue("effect")

	err := ika.callSwitch(name, t, e)
	if err != nil {
		fmt.Println(err)
	}
}
