package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {

	engine := gin.Default()
	{
		engine.GET("/getChainInfo", queryBlockchainInfo)              //查询区块链信息
		engine.POST("/addCustomerInfo", addCustomer)                  //添加客户信息
		engine.POST("/addCollateralInfo", addCollateral)              //添加押品
		engine.POST("/addProjectInfo", addProject)                    //添加项目
		engine.GET("/getCustomerInfo", queryCustomerInfo)             //查询客户信息
		engine.GET("/getHistoryCustomerInfo", getHistoryCustomer)     //客户历史信息查询
		engine.GET("/getHistoryCollateralInfo", getHistoryCollateral) //押品变更历史查询
		engine.GET("/getHistoryProjectInfo", getHistoryProject)       //项目历史信息查询

		// engine.GET("/blockchaininfo", queryBlockchainInfo)            //查询区块链信息
		// engine.POST("/customer", addCustomer)                         //添加客户信息
		// engine.POST("/collateral", addCollateral)                     //添加押品
		// engine.POST("/project", addProject)                           //添加项目
		// engine.GET("/customer", queryCustomerInfo)              //查询客户信息
		// engine.GET("/history/customer", getHistoryCustomer)     //客户历史信息查询
		// engine.GET("/history/collateral", getHistoryCollateral) //押品变更历史查询
		// engine.GET("/history/project", getHistoryProject)       //项目历史信息查询
	}

	engine.Run() // listen and serve on 0.0.0.0:8080，默认端口 8080
	// engine.Run(":8086") 通过此方式修改监听的端口，此时监听 8086 端口
}

func queryBlockchainInfo(ctx *gin.Context) {
	resp := queryBlockchain()

	ctx.JSON(http.StatusOK, resp)
}

// Customer 客户关键字
type Customer struct {
	Name         string `form:"name" binding:"required"`         //客户姓名
	ID           string `form:"id" binding:"required"`           //客户编号
	Code         string `form:"code" binding:"required"`         //统一社会信用代码
	Type         string `form:"type" binding:"required"`         //类型
	Money        string `form:"money" binding:"required"`        //注册资本
	Person       string `form:"person" binding:"required"`       //法人代表
	Date         string `form:"date" binding:"required"`         //成立日期
	BusinessDate string `form:"businessDate" binding:"required"` //营业期限
	ApprovalDate string `form:"approvalDate" binding:"required"` //核准日期
	Trade        string `form:"trade" binding:"required"`        //所属行业
}

// 添加客户信息
func addCustomer(ctx *gin.Context) {
	// 参数处理
	req := new(Customer)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	// 区块链交互
	resp, err := channelExecute("addCustomerInfo", [][]byte{
		[]byte(req.Name),
		[]byte(req.ID),
		[]byte(req.Code),
		[]byte(req.Type),
		[]byte(req.Money),
		[]byte(req.Person),
		[]byte(req.Date),
		[]byte(req.BusinessDate),
		[]byte(req.ApprovalDate),
		[]byte(req.Trade),
	})

	// 因为 postman 对于 非200-300 直接的错误，会直接返回错误编号，而不显示错误内容
	// 所以此处通过 200 直接返回，并显示错误内容
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 查询客户信息
func queryCustomerInfo(ctx *gin.Context) {
	// 若参数在 path 中，用 Param() 方法来提取参数
	// user := ctx.Param("name")
	user := ctx.Query("name")

	resp, err := channelQuery("getCustomerInfo", [][]byte{
		[]byte(user),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	//ctx.JSON(http.StatusOK, resp)
	// 返回 Payload（对于接收者有用的数据） 方便查看
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

// 查询客户历史信息
func getHistoryCustomer(ctx *gin.Context) {
	// 若参数在 path 中，用 Param() 方法来提取参数
	// user := ctx.Param("name")
	user := ctx.Query("name")

	resp, err := channelQuery("getHistoryCustomerInfo", [][]byte{
		[]byte(user),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	//ctx.JSON(http.StatusOK, resp)
	// 返回 Payload（对于接收者有用的数据） 方便查看
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

// Collateral 押品信息
type Collateral struct {
	Name           string `form:"name" binding:"required"`           //客户姓名
	CollateralID   string `form:"collateralId" binding:"required"`   //押品编号
	CollateralName string `form:"collateralName" binding:"required"` //押品名称
}

// 资产登记
func addCollateral(ctx *gin.Context) {
	req := new(Collateral)
	// 参数在 form 表单中，用 ShouldBind() 方法来提取参数
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	resp, err := channelExecute("addCollateralInfo", [][]byte{
		[]byte(req.Name),
		[]byte(req.CollateralID),
		[]byte(req.CollateralName),
	})
	fmt.Println(resp)

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 押品变更历史查询
func getHistoryCollateral(ctx *gin.Context) {
	// 若参数在 path 中，用 Param() 方法来提取参数
	// user := ctx.Param("name")
	user := ctx.Query("name")
	resp, err := channelQuery("getHistoryCollateralInfo", [][]byte{
		[]byte(user),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	//ctx.JSON(http.StatusOK, resp)
	// 返回 Payload（对于接收者有用的数据） 方便查看
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

// Project 项目信息
type Project struct {
	Name               string `form:"name" binding:"required"`               //客户名称
	ProjectName        string `form:"projectName" binding:"required"`        //项目名称
	ProjectID          string `form:"projectId" binding:"required"`          //项目编号
	ProjectType        string `form:"projectType" binding:"required"`        //业务类型
	ProjectTrade       string `form:"projectTrade" binding:"required"`       //所属行业
	ProjectDate        string `form:"projectDate" binding:"required"`        //批复下达日
	ProjectApprove     string `form:"projectApprove" binding:"required"`     //审批是否通过
	ProjectPart        string `form:"projectPart" binding:"required"`        //是否成立有限合伙人
	ProjectInvest      string `form:"projectInvest" binding:"required"`      //是否有自有资金投资
	ProjectMoney       string `form:"projectMoney" binding:"required"`       //持有债券金额
	ProjectCompanyType string `form:"projectCompanyType" binding:"required"` //被投资企业类型
}

//
func addProject(ctx *gin.Context) {
	req := new(Project)
	// 参数在 form 表单中，用 ShouldBind() 方法来提取参数

	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	resp, err := channelExecute("addProjectInfo", [][]byte{
		[]byte(req.Name),
		[]byte(req.ProjectName),
		[]byte(req.ProjectID),
		[]byte(req.ProjectType),
		[]byte(req.ProjectTrade),
		[]byte(req.ProjectDate),
		[]byte(req.ProjectApprove),
		[]byte(req.ProjectPart),
		[]byte(req.ProjectInvest),
		[]byte(req.ProjectMoney),
		[]byte(req.ProjectCompanyType),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 项目变更历史查询
func getHistoryProject(ctx *gin.Context) {
	// 若参数在 path 中，用 Param() 方法来提取参数
	// user := ctx.Param("name")
	user := ctx.Query("name")
	resp, err := channelQuery("getHistoryProjectInfo", [][]byte{
		[]byte(user),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

// 代码中用到的名字都是 yaml文件中的 key 而不是 value
var (
	sdk           *fabsdk.FabricSDK
	channelName   = "mychannel"
	chaincodeName = "assetscc"
	org           = "org1" // 对应了 configtx.yaml 文件的160行
	user          = "Admin"
	configPath    = "./config.yaml"
)

// 初始化 SDK，需要用到 配置文件：config.yaml
func init() {
	var err error
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}

// 区块链管理
func manageBlockchain() {
	// 表明身份
	ctx := sdk.Context(fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := resmgmt.New(ctx) // resource management 资源管理包 resmgmt
	if err != nil {
		panic(err)
	}

	// 具体操作
	// SaveChannel 创建or更新通道
	cli.SaveChannel(resmgmt.SaveChannelRequest{}, resmgmt.WithOrdererEndpoint("orderer.example.com"),
		resmgmt.WithTargetEndpoints())
}

// 区块链查询  账本查询
func queryBlockchain() *fab.BlockchainInfoResponse {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := ledger.New(ctx) // 实例化一个账本客户端
	if err != nil {
		panic(err)
	}

	resp, err := cli.QueryInfo(ledger.WithTargetEndpoints("peer0.org1.example.com"))
	if err != nil {
		panic(err)
	}

	// fmt.Println(resp)

	// // 1 QueryBlockByHash queries the ledger for block by block hash.
	// cli.QueryBlockByHash(resp.BCI.CurrentBlockHash)
	// // height, currentBlockHash, previousBlockHash := resp.BCI.Height, resp.BCI.CurrentBlockHash, resp.BCI.PreviousBlockHash

	// // 2 从第0个块开始遍历
	// for i := uint64(0); i <= resp.BCI.Height; i++ {
	// 	cli.QueryBlock(i)
	// }
	// // 以上两种方式都可以实现区块链浏览器读取区块的功能

	// high := resp.BCI.Height
	// return height, currentBlockHash, previousBlockHash
	return resp
}

// 区块链交互
func channelExecute(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 状态更新，insert/update/delete
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.example.com"))
	if err != nil {
		return channel.Response{}, err
	}

	// 链码事件监听
	go func() {
		// channel
		reg, ccevt, err := cli.RegisterChaincodeEvent(chaincodeName, "eventname")
		if err != nil {
			return
		}
		defer cli.UnregisterChaincodeEvent(reg)

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-ccevt:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("event timeout, exit!")
				return
			}
		}

		// event
		// eventcli, err := event.New(ctx)
		// if err != nil {
		// 	return
		// }
		// eventcli.RegisterChaincodeEvent(chaincodeName, "eventname")
		// ... same as channel moudle
		//
	}()

	// 交易状态事件监听
	go func() {
		eventcli, err := event.New(ctx)
		if err != nil {
			return
		}

		reg, status, err := eventcli.RegisterTxStatusEvent(string(resp.TransactionID))
		defer eventcli.Unregister(reg) // 注册必有注销

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-status:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("event timeout, exit!")
				return
			}
		}
	}()

	return resp, nil
}

func channelQuery(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 状态的查询，select
	return cli.Query(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.example.com"))
}

// 事件监听
func eventHandle() {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := event.New(ctx)
	if err != nil {
		panic(err)
	}

	// 交易状态事件		--区块链交互模块
	// 链码事件 业务事件	--区块链交互模块

	// 区块事件
	reg, blkevent, err := cli.RegisterBlockEvent()
	if err != nil {
		panic(err)
	}
	defer cli.Unregister(reg)

	timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		select {
		case evt := <-blkevent:
			fmt.Printf("received a block, %+v", evt)
		case <-timeoutctx.Done():
			fmt.Println("event timeout, exit!")
			return
		}
	}
}
