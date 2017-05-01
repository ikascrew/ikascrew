package server

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/effect"
	"github.com/secondarykey/ikascrew/pb"
	"github.com/secondarykey/ikascrew/video"
)

func init() {
}

func (i *IkascrewServer) startRPC() error {

	lis, err := net.Listen("tcp", ":55555")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterIkascrewServer(s, i)

	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}

func (i *IkascrewServer) Sync(ctx context.Context, r *pb.SyncRequest) (*pb.SyncReply, error) {
	i.window.FullScreen()

	rep := &pb.SyncReply{
		Source:  "logo.png",
		Type:    "image",
		Project: ikascrew.ProjectName(),
	}

	return rep, nil
}

func (i *IkascrewServer) Effect(ctx context.Context, r *pb.EffectRequest) (*pb.EffectReply, error) {

	rep := &pb.EffectReply{
		Success: false,
	}
	name := ikascrew.ProjectName() + "/" + r.Name

	var v ikascrew.Video
	var err error

	switch r.Type {
	case "file":
		v, err = video.NewFile(name)
	case "image":
		v, err = video.NewImage(name)
	case "mic":
		v, err = video.NewMicrophone()
	default:
		err = fmt.Errorf("Not Support Type[%s]", r.Type)
	}

	if err != nil {
		return rep, err
	}

	var e ikascrew.Effect
	switch r.Effect {
	case "switch":
		e, err = effect.NewSwitch(v, i.window.GetEffect())
	case "mate":
		now := i.window.GetEffect()
		switch now.(type) {
		case *effect.Mate:
			err = fmt.Errorf("MateEffect Can not be used continuously.")
		default:
			e, err = effect.NewMate(v, now)
		}
	default:
		err = fmt.Errorf("Not Support Effect[%s]", r.Effect)
	}

	if err != nil {
		return rep, err
	}

	i.window.SetEffect(e)

	rep.Success = true
	return rep, nil
}
