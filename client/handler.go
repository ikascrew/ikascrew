package client

import (
	"fmt"
	"net/http"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/effect"
	"github.com/secondarykey/ikascrew/video"
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

	var v ikascrew.Video
	var err error

	file := ikascrew.ProjectName() + "/" + name

	switch t {
	case "file":
		v, err = video.NewFile(file)
	case "image":
		v, err = video.NewImage(file)
	case "mic":
		v, err = video.NewMicrophone()
	default:
		err = fmt.Errorf("Not Support Type[%s]", t)
	}
	if err != nil {
		fmt.Println("Error createVideo:", err)
	}

	e, err := effect.NewNormal(v)

	ika.window.SetEffect(e)
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
