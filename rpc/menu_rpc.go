package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"golang.org/x/net/context"
)

type Rpc_menu struct {
}

func (this *Rpc_menu) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {

	list := model.Modules_List("")
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)

	return resp, nil
}

func (this *Rpc_menu) Add(ctx context.Context, request *pb.Req) (*pb.Resp, error) {
	msg := request.GetMessage()
	Obj := model.Modules{}
	util.StringToObj(msg, &Obj)
	_, err := model.Modules_Add_Menu(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_menu) Update(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.Modules{}
	util.StringToObj(msg, &Obj)
	err = model.Modules_Update(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_menu) Del(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.Modules_Del(id)
	return &pb.Resp{}, nil
}

func (this *Rpc_menu) GetAllMenu(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {

	list := model.All_Modules_List("")
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_menu) GetUserMenu(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id, _ := strconv.Atoi(req.Id)
	list := model.GetUserMenu(user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_menu) GetUserAction(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id, _ := strconv.Atoi(req.Id)
	list := model.GetUserAction(user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_menu) GetMenu(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id, _ := strconv.Atoi(req.Id)
	list := model.Get_Menu(user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_menu) AddMenu(ctx context.Context, request *pb.Req) (*pb.Resp, error) {
	msg := request.GetMessage()
	Obj := model.Menu{}
	util.StringToObj(msg, &Obj)
	_, err := model.Add_Menu(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_menu) DelMenu(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.Del_Menu(id)
	return &pb.Resp{}, nil
}

func (this *Rpc_menu) EditMenu(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.Menu{}
	util.StringToObj(msg, &Obj)
	err = model.Edit_Menu(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_menu) StartStopMenu(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.Menu{}
	util.StringToObj(msg, &Obj)
	err = model.Start_Stop_Menu(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_menu) GetMenuWithAuth(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id, _ := strconv.Atoi(req.Id)
	list := model.Get_Menu_With_Auth(user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	return resp, nil
}

func (this *Rpc_menu) SortMenu(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	Obj := model.Menu_Sort{}
	util.StringToObj(msg, &Obj)
	err = model.Sort_Menu(&Obj)
	return &pb.Resp{}, err
}
