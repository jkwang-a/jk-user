package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"golang.org/x/net/context"
)

type Rpc_flow_chart struct {
}

func (this *Rpc_flow_chart) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id := int(req.UserId)
	list, total_page, total := model.Flow_Chart_List(req.Key, int(req.Start), int(req.Length), req.OrderField, req.Desc, user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	resp.Total = int64(total)
	resp.PageTotal = int64(total_page)
	return resp, nil
}

func (this *Rpc_flow_chart) Add(ctx context.Context, request *pb.Req) (*pb.Resp, error) {
	user_id := int(request.UserId)
	msg := request.GetMessage()
	Obj := model.Flow_Chart{}
	Obj.User_id = user_id
	util.StringToObj(msg, &Obj)
	_, err := model.Flow_chart_Add(&Obj)
	return &pb.Resp{}, err
}

func (this *Rpc_flow_chart) Update(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	msg := request.GetMessage()
	obj := model.Flow_Chart{}
	util.StringToObj(msg, &obj)
	e := model.Flow_Chart_Update(&obj)
	return &pb.Resp{}, e
}

func (this *Rpc_flow_chart) Del(ctx context.Context,
	request *pb.Req) (response *pb.Resp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.Flow_Chart_Del(id)
	return &pb.Resp{}, err

}
