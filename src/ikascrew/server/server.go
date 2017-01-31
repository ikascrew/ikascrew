package server

import (
	"fmt"
	"ikascrew"
	"net/http"
	"time"
)

type Sync struct {
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	Frame int    `json:"frame"`
}

type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Information struct {
	Success bool   `json:"success"`
	Video1  string `json:"video1"`
	Video2  string `json:"video2"`
}

type IkascrewServer struct {
	q *ikascrew.Queue
}

const ADDRESS = "localhost:5555"

func init() {
}

func Address() string {
	return ADDRESS
}

func Start(d, f string) error {

	project := d
	err := ikascrew.Loading(project)
	if err != nil {
		return fmt.Errorf("Error Loading directory:%s", err)
	}

	ikascrew.PrintVideos()

	q, err := ikascrew.NewQueue(f)
	if err != nil {
		return fmt.Errorf("Error New Queue:%s", err)
	}

	ika := &IkascrewServer{
		q: q,
	}

	win := ikascrew.NewWindow("ikascrew")
	time.Sleep(300 * time.Millisecond)

	go func() {
		win.Play(q)
	}()

	http.HandleFunc("/sync", ika.syncHandler)
	http.HandleFunc("/push", ika.pushHandler)
	http.HandleFunc("/switch", ika.switchHandler)
	http.HandleFunc("/effect", ika.effectHandler)
	http.HandleFunc("/info", ika.informationHandler)
	return http.ListenAndServe(ADDRESS, nil)
}

func (i *IkascrewServer) syncHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("# Request Sync")

	v1, v2 := i.q.Name()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{")
	fmt.Fprintf(w, "\"v1\" : \"%s\",", v1)
	fmt.Fprintf(w, "\"v2\" : \"%s\",", v2)
	fmt.Fprintf(w, "\"frame\" : %d", i.q.Current())
	fmt.Fprintf(w, "}")
}

func (i *IkascrewServer) effectHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("# Request Effect")

	success := "true"
	effect := r.FormValue("effect")

	v, err := ikascrew.GetVideo(effect)
	if err != nil {
		success = "false"
	} else {
		i.q.Effect(v)
		err = fmt.Errorf("Effect Video:%v" + effect)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{")
	fmt.Fprintf(w, "\"success\" : %s,", success)
	fmt.Fprintf(w, "\"message\" : \"%s\"", err)
	fmt.Fprintf(w, "}")
}

func (i *IkascrewServer) pushHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("# Request Push")

	success := "true"
	next := r.FormValue("next")

	fmt.Printf("[%s]\n", next)

	v, err := ikascrew.GetVideo(next)
	if err != nil {
		success = "false"
	} else {
		i.q.Sub(v)
		err = fmt.Errorf("Next Video:%v" + next)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{")
	fmt.Fprintf(w, "\"success\" : %s,", success)
	fmt.Fprintf(w, "\"message\" : \"%s\"", err)
	fmt.Fprintf(w, "}")
}

func (i *IkascrewServer) switchHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("# Request Switch")

	success := "true"
	w.Header().Set("Content-Type", "application/json")

	dur := 200
	i.q.Switch(dur)

	fmt.Fprintf(w, "{")
	fmt.Fprintf(w, "\"success\" : %s,", success)
	fmt.Fprintf(w, "\"message\" : \"%s %d\"", "done!", dur)
	fmt.Fprintf(w, "}")
}

func (i *IkascrewServer) informationHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("# Request Infomation")
	success := "true"
	w.Header().Set("Content-Type", "application/json")

	v1, v2 := i.q.Name()

	fmt.Fprintf(w, "{")
	fmt.Fprintf(w, "\"success\" : %s,", success)
	fmt.Fprintf(w, "\"video1\" : \"%s\",", v1)
	fmt.Fprintf(w, "\"video2\" : \"%s\"", v2)
	fmt.Fprintf(w, "}")
}
