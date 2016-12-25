package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"
	"time"

	"ikascrew"
	"ikascrew/server"
)

func init() {
	fmt.Println("########################## Starting ikascrew Client")
}

type ikascrewClient struct {
	q *ikascrew.Queue
}

func Start(d string) error {

	err := ikascrew.Loading(d)
	if err != nil {
		return err
	}

	ika := &ikascrewClient{}
	s, err := ika.getSync()
	if err != nil {
		return err
	}

	q, err := ikascrew.NewSourceQueue(s.V1, s.Frame+500)
	if err != nil {
		return err
	}

	v2, err := ikascrew.GetSource(s.V2)
	if err != nil {
		fmt.Println(err)
	} else {
		q.V2 = v2
	}

	ika.q = q

	win := ikascrew.NewWindow("ikascrew client")
	time.Sleep(300 * time.Millisecond)
	go func() {
		win.Play(q)
	}()

	return ika.input()
}

func (c *ikascrewClient) info() error {

	values := url.Values{}
	resp, err := http.PostForm("http://"+server.Address()+"/info", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	var i server.Information
	err = json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	if !i.Success {
		return fmt.Errorf("Error")
	}

	fmt.Println("Video1=" + i.Video1)
	fmt.Println("Video2=" + i.Video2)

	return nil
}

func (c *ikascrewClient) remote(cmd []string) error {

	command := cmd[0]

	values := url.Values{}
	if command == "effect" {
		if len(cmd) < 2 {
			return fmt.Errorf("remote effect arg 2")
		}
		values["effect"] = cmd[1:]
	}
	resp, err := http.PostForm("http://"+server.Address()+"/"+command, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	var m server.Message
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if !m.Success {
		return fmt.Errorf(m.Message)
	}

	return nil
}

func (c *ikascrewClient) push(f string) error {

	values := url.Values{}
	values.Add("next", f)

	resp, err := http.PostForm("http://"+server.Address()+"/push", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	var m server.Message
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	if !m.Success {
		return fmt.Errorf(m.Message)
	}

	v, err := ikascrew.GetVideo(f)
	if err != nil {
		return err
	}
	c.q.Sub(v)

	return nil
}

func (c *ikascrewClient) getSync() (*server.Sync, error) {
	resp, err := http.Get("http://" + server.Address() + "/sync")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	var s server.Sync
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *ikascrewClient) ls() error {
	d := ikascrew.List()
	fmt.Println("##########################")
	for _, elm := range d {
		fmt.Println(elm)
	}
	fmt.Println("##########################")
	return nil
}

func (c *ikascrewClient) input() error {

	var sc = bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("ikascrew> ")
		sc.Scan()
		cmd := sc.Text()
		cmds := strings.Split(cmd, " ")

		switch cmds[0] {
		case "info":
			err := c.info()
			if err != nil {
				fmt.Println("Error Information:", err)
			}
		case "sync":
			err := c.sync()
			if err != nil {
				fmt.Println("Error Sync:", err)
			}
		case "ls":
			err := c.ls()
			if err != nil {
				fmt.Println("Error Ls:", err)
			}
		case "effect":
			if len(cmds) < 2 {
				fmt.Println("Error Effect Arg length 2")
			} else {
				err := c.effect(cmds[1])
				if err != nil {
					fmt.Println("Error Effect:", err)
				}
			}
		case "switch":
			name := ""
			if len(cmds) >= 2 {
				name = cmds[1]
			}

			err := c.switchVideo(name)
			if err != nil {
				fmt.Println("Error Switch:", err)
			}
		case "push":
			if len(cmds) < 2 {
				fmt.Println("Error Push Arg length 2")
			} else {
				err := c.push(cmds[1])
				if err != nil {
					fmt.Println("Error Push:", err)
				}
			}
		case "load":
			if len(cmds) < 2 {
				fmt.Println("Error Load Arg length 2")
			} else {
				err := c.load(cmds[1])
				if err != nil {
					fmt.Println("Error Load:", err)
				}
			}
		case "remote":
			if len(cmds) < 2 {
				fmt.Println("Error Remote Arg length 2")
			} else {
				err := c.remote(cmds[1:])
				if err != nil {
					fmt.Println("Error remote:", err)
				}
			}
		case "q":
			in := c.quit()
			if in == "Y" {
				fmt.Println("Bye!")
				return nil
			}
		case "":
		default:
			fmt.Println("Command Not Found:" + cmds[0])
		}

	}
	return nil
}

func (c *ikascrewClient) switchVideo(e string) error {

	if e != "" {
		v, err := ikascrew.GetVideo(e)
		if err != nil {
			return err
		}
		return c.q.EffectSwitch(v)
	}
	return c.q.Switch(200)
}

func (c *ikascrewClient) effect(f string) error {
	v, err := ikascrew.GetVideo(f)
	if err != nil {
		return err
	}

	c.q.Effect(v)
	return nil
}

func (c *ikascrewClient) sync() error {

	s, err := c.getSync()
	if err != nil {
		return err
	}

	v, err := ikascrew.GetSource(s.V1)
	if err != nil {
		return err
	}

	c.q.Set(v, s.Frame+200)

	v2, err := ikascrew.GetSource(s.V2)
	if err != nil {
		return err
	}
	c.q.Sub(v2)

	return nil
}

func (c *ikascrewClient) load(f string) error {

	v, err := ikascrew.GetVideo(f)
	if err != nil {
		return err
	}

	c.q.Set(v, 0)
	return nil
}

func (c *ikascrewClient) quit() string {
	var sc = bufio.NewScanner(os.Stdin)
	fmt.Println("\nByeBye?[Y/n]")
	sc.Scan()
	yN := sc.Text()
	return yN
}
