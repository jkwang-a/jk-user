package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"

	"golang.org/x/net/context"
)

type Rpc_log struct {
}

func (this *Rpc_log) List(ctx context.Context,
	request *pb.LogFindReq) (response *pb.LogResp, err error) {
	total, list := model.Log_Query(int(request.Start), int(request.Length), "", "", request.StartTime, request.EndTime, request.Key)
	return &pb.LogResp{Message: util.ObjToString(list), Total: int64(total)}, nil
}

func (this *Rpc_log) Save(ctx context.Context, req *pb.LogReq) (*pb.LogResp, error) {

	obj := model.Log{}
	obj.Action = req.Action
	obj.Client_ip = req.ClientIp
	obj.Content = req.Content
	obj.User_id = int(req.UserId)
	obj.Create_time = util.GetUtcTimeStr()
	model.Log_Save(obj)
	return &pb.LogResp{}, nil
}

func (this *Rpc_log) Del(ctx context.Context,
	request *pb.LogReq) (response *pb.LogResp, err error) {
	return &pb.LogResp{}, nil

}
