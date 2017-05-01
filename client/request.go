package client

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/secondarykey/ikascrew/pb"
)

func (i *IkascrewClient) syncServer() error {

	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	r, err := c.Sync(context.Background(), &pb.SyncRequest{})
	if err != nil {
		return err
	}

	fmt.Printf("Success[%s]\n", r)

	return nil
}

func (i *IkascrewClient) callSwitch(f string, t string, e string) error {
	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewIkascrewClient(conn)

	r, err := c.Effect(context.Background(), &pb.EffectRequest{
		Name:   f,
		Type:   t,
		Effect: e,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Success[%s]\n", r)
	return nil
}
