package client

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/ikascrew/ikascrew/pb"
)

func (i *IkascrewClient) syncServer() (*pb.SyncReply, error) {

	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	r, err := c.Sync(context.Background(), &pb.SyncRequest{})
	if err != nil {
		return nil, err
	}

	fmt.Println(r)

	return r, nil
}

func (i *IkascrewClient) callEffect(f string, t string) error {

	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	r, err := c.Effect(context.Background(), &pb.EffectRequest{
		Name: f,
		Type: t,
	})

	if err != nil {
		return err
	}

	fmt.Println(r)
	return nil
}

func (i *IkascrewClient) callNext() error {
	return i.callSwitch("next")
}

func (i *IkascrewClient) callPrev() error {
	return i.callSwitch("prev")
}

func (i *IkascrewClient) callSwitch(t string) error {

	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	r, err := c.Switch(context.Background(), &pb.SwitchRequest{
		Type: t,
	})

	if err != nil {
		return err
	}

	fmt.Println(r)
	return nil
}
