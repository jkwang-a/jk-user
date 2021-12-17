package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"golang.org/x/net/context"
)

type Rpc_role struct {
}

func (this *Rpc_role) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id := int(req.UserId)
	list, total_page, total := model.Role_List(req.Key, int(req.Start), int(req.Length), req.OrderField, req.Desc, user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	resp.Total = int64(total)
	resp.PageTotal = int64(total_page)
	return resp, nil
}

func (this *Rpc_role) Add(ctx context.Context, request *pb.RoleReq) (*pb.RoleResp, error) {
	msg := request.GetMessage()
	Role := model.Role{}
	util.StringToObj(msg, &Role)
	_, err := model.Role_Add(&Role)
	return &pb.RoleResp{}, err
}

func (this *Rpc_role) Del(ctx context.Context,
	request *pb.RoleReq) (response *pb.RoleResp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.Role_Del(id)
	return &pb.RoleResp{}, err

}
func (this *Rpc_role) Update(ctx context.Context,
	request *pb.RoleReq) (response *pb.RoleResp, err error) {
	msg := request.GetMessage()
	obj := model.Role{}
	util.StringToObj(msg, &obj)
	err = model.Role_Update(&obj)
	return &pb.RoleResp{}, err
}
