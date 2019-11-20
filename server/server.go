package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/lazydog/grpc2wayConfirmationChat/protos"
	myutils "github.com/lazydog/grpc2wayConfirmationChat/utils"
)


type genericMultiSessionGroupChat struct {
	pb.UnimplementedMultiSessionChatServerServer
	registeredUsers map[string]string //map of username:usersecret
	userRecvStream map[string]pb.MultiSessionChatServer_RecvrStreamServer
	userSendStream map[string]pb.MultiSessionChatServer_SenderStreamServer
	sessions map[string][]string//session_secret:usernames list
	sessionMessages map[string][]
}


func (g *genericMultiSessionGroupChat) requestServer(ctx context.Context, req *pb.UserRequest) (*pb.UserReqResponse, error){
	switch req.Type {
	case pb.UserRequestType_REGISTER_USER:
		if _, ok := g.registeredUsers[req.Username]; ok {
			return &pb.UserReqResponse{Response: "Username already exists!!", Status: false}, nil
		} else {
			hash := myutils.NewSHA1Hash()
			g.registeredUsers[req.Username] = hash
			return &pb.UserReqResponse{Response: hash, Status: true}, nil
		}
	case pb.UserRequestType_START_SESSION:
		session_secret := myutils.NewSHA1Hash()
		session_users := make([]string,0)
		for _, username := range req.SessionUsernames{
			if _,ok:= g.registeredUsers[username]; !ok{
				return &pb.UserReqResponse{Status:false, Response:fmt.Sprintf("Username %s does not exists..", username)}, nil
			}
			session_users = append(session_users, username)
		}
		g.sessions[session_secret] = session_users
		return &pb.UserReqResponse{Response:session_secret, Status:true}, nil
	}
	return nil,nil
}