package rpc

import (
	"jk-common/util"
	"jk-user/model"
	pb "jk-user/proto"
	"strconv"

	"github.com/astaxie/beego"

	"golang.org/x/net/context"
)

type Rpc_user struct {
}

func (this *Rpc_user) List(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	user_id := int(req.UserId)
	list, total_page, total := model.User_List(req.Key, int(req.Start), int(req.Length), req.OrderField, req.Desc, user_id)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	resp.Total = int64(total)
	resp.PageTotal = int64(total_page)
	return resp, nil
}

func (this *Rpc_user) Add(ctx context.Context, request *pb.UserReq) (*pb.UserResp, error) {
	msg := request.GetMessage()
	obj := model.Users{}
	util.StringToObj(msg, &obj)
	_, err := model.User_Add(&obj)
	return &pb.UserResp{}, err
}

func (this *Rpc_user) Update(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	msg := request.GetMessage()
	obj := model.Users{}
	util.StringToObj(msg, &obj)
	e := model.User_Update(&obj)
	return &pb.UserResp{}, e
}

func (this *Rpc_user) Del(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	id, _ := strconv.Atoi(request.GetMessage())
	err = model.User_Del(id)
	return &pb.UserResp{}, err

}
func (this *Rpc_user) Login(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	obj := model.Users{}
	util.StringToObj(request.Message, &obj)
	//	util.PrintLog(request.Message, obj.Password)
	info, e := model.Login(obj.Email, obj.Password)
	return &pb.UserResp{Message: util.ObjToString(info)}, e
}

func (this *Rpc_user) Get_Mgr_UserList(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	d := `{"success":true}`
	return &pb.UserResp{Message: d}, err
}

func (this *Rpc_user) Login_AZure(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	d, e := model.User_Login_AZure(request.Message)
	//util.PrintLog(util.ObjToString(d))
	return &pb.UserResp{Message: util.ObjToString(d)}, e
}

func (this *Rpc_user) GetUserByKey(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	d, e := model.User_GetByApiKey(request.Message)
	//util.PrintLog(util.ObjToString(d))
	return &pb.UserResp{Message: util.ObjToString(d)}, e
}
func (this *Rpc_user) GetUserById(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	d, e := model.User_GetById(request.Message)
	//util.PrintLog(util.ObjToString(d))
	return &pb.UserResp{Message: util.ObjToString(d)}, e
}

func (this *Rpc_user) AccessStatistics(ctx context.Context,
	req *pb.UserReq) (resp *pb.DashboardAccessResp, err error) {
	access_list := model.AccessStatistics()
	resp = &pb.DashboardAccessResp{}
	resp.AccessList = util.ObjToString(access_list)
	return resp, nil
}

func (this *Rpc_user) GetUserList(ctx context.Context,
	req *pb.PagerReq) (resp *pb.PagerResp, err error) {
	//	user_id := int(req.UserId)
	total, list := model.Get_User_List(int(req.Start), int(req.Length), req.Key, req.Id, req.Status, req.OrgId, req.Message, req.QueryUserId)
	resp = &pb.PagerResp{}
	resp.Data = util.ObjToString(list)
	resp.Total = int64(total)
	return resp, nil
}

func (this *Rpc_user) AddUser(ctx context.Context, request *pb.UserReq) (*pb.UserResp, error) {
	msg := request.GetMessage()
	obj := model.User_Edit{}
	util.StringToObj(msg, &obj)
	_, err := model.Add_User(&obj)
	return &pb.UserResp{}, err
}

func (this *Rpc_user) UpdateUser(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	msg := request.GetMessage()
	obj := model.User_Edit{}
	util.StringToObj(msg, &obj)
	e := model.Update_User(&obj)
	return &pb.UserResp{}, e
}

func (this *Rpc_user) DelUser(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	msg := request.GetMessage()
	Obj := model.User_Del_Edit{}
	util.StringToObj(msg, &Obj)
	err = model.Del_User(&Obj)
	return &pb.UserResp{}, err

}

func (this *Rpc_user) EditUserStatus(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	msg := request.GetMessage()
	Obj := model.User_Del_Edit{}
	util.StringToObj(msg, &Obj)
	err = model.Edit_User_Status(&Obj)
	return &pb.UserResp{}, err

}

func (this *Rpc_user) UserLogin(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	obj := model.Users{}
	util.StringToObj(request.Message, &obj)
	//	util.PrintLog(request.Message, obj.Password)
	info, e := model.User_Login(obj.Email, obj.Password, request.ClientIp)
	return &pb.UserResp{Message: util.ObjToString(info)}, e
}

func (this *Rpc_user) User_Login_AZure(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	d, e := model.User_Login_With_AZure(request.Message, request.ClientIp)
	//util.PrintLog(util.ObjToString(d))
	return &pb.UserResp{Message: util.ObjToString(d)}, e
}

func (this *Rpc_user) GetUserDetailByKey(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	d, e := model.Get_User_Detail_By_Key(request.Message)
	//util.PrintLog(util.ObjToString(d))
	return &pb.UserResp{Message: util.ObjToString(d)}, e
}

func (this *Rpc_user) GetUserIdWithParentId(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	beego.Debug(request.Message)
	d, e := model.Get_UserId_With_ParentId(request.Message)
	return &pb.UserResp{Message: d}, e
}

func (this *Rpc_user) UpdateNewUser(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	e := model.Update_New_User(request.Email, request.ApiKey)
	return &pb.UserResp{}, e
}

func (this *Rpc_user) UpdateAllUser(ctx context.Context,
	request *pb.UserReq) (response *pb.UserResp, err error) {
	e := model.Update_All_User(request.Message)
	return &pb.UserResp{}, e
}
