package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"golang.org/x/net/context"
)

type Rpc_user_group struct {
}

func (this *Rpc_user_group) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	//	user_id := int(req.UserId)
	list := model.User_Group_List(req.Status)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_user_group) Add(ctx context.Context, request *pb.Req) (*pb.Resp, error) {
	msg := request.GetMessage()
	Obj := model.UserGroup{}
	util.StringToObj(msg, &Obj)
	_, err := model.User_Group_Add(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_user_group) Update(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.UserGroup{}
	util.StringToObj(msg, &Obj)
	err = model.User_Group_Update(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_user_group) Del(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.User_Group_Del(id)
	return &pb.Resp{}, nil
}

func (this *Rpc_user_group) EditStatus(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.UpdateStatus{}
	util.StringToObj(msg, &Obj)
	err = model.Edit_User_Group_Status(&Obj)
	return &pb.Resp{}, err
}
