package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"golang.org/x/net/context"
)

type Rpc_rule struct {
}

func (this *Rpc_rule) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	//	user_id := int(req.UserId)
	list := model.Rule_List()
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_rule) Add(ctx context.Context, request *pb.Req) (*pb.Resp, error) {
	msg := request.GetMessage()
	Obj := model.Rule{}
	util.StringToObj(msg, &Obj)
	_, err := model.Rule_Add(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_rule) Update(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.Rule{}
	util.StringToObj(msg, &Obj)
	err = model.Rule_Update(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_rule) Del(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.Rule_Del(id)
	return &pb.Resp{}, nil
}

func (this *Rpc_rule) BatchDel(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.BatchDel{}
	util.StringToObj(msg, &Obj)
	err = model.Batch_Del(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_rule) EditRuleStatus(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.BatchDel{}
	util.StringToObj(msg, &Obj)
	err = model.Edit_Rule_Status(&Obj)
	return &pb.Resp{}, err
}
