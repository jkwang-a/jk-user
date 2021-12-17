package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"golang.org/x/net/context"
)

type Rpc_org struct {
}

func (this *Rpc_org) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {

	list := model.Org_List("")
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)

	return resp, nil
}

func (this *Rpc_org) Add(ctx context.Context, request *pb.OrgReq) (*pb.OrgResp, error) {
	msg := request.GetMessage()
	Org := model.Org{}
	util.StringToObj(msg, &Org)
	_, err := model.Org_Add(&Org)
	return &pb.OrgResp{}, err
}

func (this *Rpc_org) Update(ctx context.Context,
	request *pb.OrgReq) (response *pb.OrgResp, err error) {
	msg := request.GetMessage()
	Org := model.Org{}
	util.StringToObj(msg, &Org)
	err = model.Org_Update(&Org)
	return &pb.OrgResp{}, err
}

func (this *Rpc_org) Del(ctx context.Context,
	request *pb.OrgReq) (response *pb.OrgResp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.Org_Del(id)
	return &pb.OrgResp{}, nil
}

func (this *Rpc_org) GetOrgList(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	list := model.Get_Org_List()
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_org) EditOrgStatus(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.BatchDel{}
	util.StringToObj(msg, &Obj)
	err = model.Edit_Org_Status(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_org) BatchDelOrg(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.BatchDel{}
	util.StringToObj(msg, &Obj)
	err = model.Batch_Del_Org(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_org) SortOrg(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.Sort{}
	util.StringToObj(msg, &Obj)
	err = model.Sort_Org(&Obj)
	return &pb.Resp{}, err
}
