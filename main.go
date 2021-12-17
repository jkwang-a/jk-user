package main

import (
	"jk-user/model"
	rpc_pb "jk-user/proto"
	"jk-user/rpc"

	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis-protocol/server/grpc"

	//	_ "github.com/go-chassis/go-chassis-protocol/server/highway"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
)

func main() {

	user := rpc.Rpc_user{}
	chassis.RegisterSchema("grpc", &user,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterUserServer
		}),
	)

	log := rpc.Rpc_log{}
	chassis.RegisterSchema("grpc", &log,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterLogServer
		}),
	)
	role := rpc.Rpc_role{}
	chassis.RegisterSchema("grpc", &role,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterRoleServer
		}),
	)

	org := rpc.Rpc_org{}
	chassis.RegisterSchema("grpc", &org,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterOrgServer
		}),
	)

	menu := rpc.Rpc_menu{}
	chassis.RegisterSchema("grpc", &menu,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterMenuServer
		}),
	)

	flowChart := rpc.Rpc_flow_chart{}
	chassis.RegisterSchema("grpc", &flowChart,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterFlowChartServer
		}),
	)

	userGroup := rpc.Rpc_user_group{}
	chassis.RegisterSchema("grpc", &userGroup,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterUserGroupServer
		}),
	)

	rule := rpc.Rpc_rule{}
	chassis.RegisterSchema("grpc", &rule,
		server.RegisterOption(func(o *server.RegisterOptions) {
			o.RPCSvcDesc = rpc_pb.RegisterRuleServer
		}),
	)
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.")
		return
	}
	model.Init()
	chassis.Run()
}
