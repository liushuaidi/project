package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// AssertsManageCC 实现一个智能合约来管理不良资产
type AssertsManageCC struct {
}

// Customer 客户名下所有信息汇总
type Customer struct {
	Name           string           `json:"Name"`           //客户名称
	CustomerInfo   CustomerInfo     `json:"customerInfo"`   //客户信息
	CollateralInfo []CollateralInfo `json:"collateralInfo"` //押品信息
	ProjectInfo    ProjectInfo      `json:"projectInfo"`    //项目信息
}

// CustomerInfo 客户信息
type CustomerInfo struct {
	ID           string `json:"id"`           //客户编号
	Code         string `json:"code"`         //统一社会信用代码
	Type         string `json:"type"`         //类型
	Money        string `json:"money"`        //注册资本
	Person       string `json:"person"`       //法人代表
	Date         string `json:"date"`         //成立日期
	BusinessDate string `json:"businessDate"` //营业期限
	ApprovalDate string `json:"approvalDate"` //核准日期
	Trade        string `json:"trade"`        //所属行业
}

// CollateralInfo 押品信息
type CollateralInfo struct {
	CollateralID   string `json:"collateralId"`   //押品编号
	CollateralName string `json:"collateralName"` //押品名称
}

// ProjectInfo 项目信息
type ProjectInfo struct {
	ProjectName        string `json:"projectName"`        //项目名称
	ProjectID          string `json:"projectId"`          //项目编号
	ProjectType        string `json:"projectType"`        //业务类型
	ProjectTrade       string `json:"projectTrade"`       //所属行业
	ProjectDate        string `json:"projectDate"`        //批复下达日
	ProjectApprove     string `json:"projectApprove"`     //审批是否通过
	ProjectPart        string `json:"projectPart"`        //是否成立有限合伙人
	ProjectInvest      string `json:"projectInvest"`      //是否有自有资金投资
	ProjectMoney       string `json:"projectMoney"`       //持有债券金额
	ProjectCompanyType string `json:"projectCompanyType"` //被投资企业类型
}

// HistoryProjectInfo 项目历史信息
type HistoryProjectInfo struct {
	TxID        string      `json:"txid"`        //交易id
	Time        string      `json:"time"`        //交易时间戳
	ProjectInfo ProjectInfo `json:"ProjectInfo"` //项目信息
}

// HistoryCustomerInfo 客户历史信息
type HistoryCustomerInfo struct {
	TxID         string       `json:"txid"`
	Time         string       `json:"time"`
	CustomerInfo CustomerInfo `json:"customerInfo"`
}

// HistoryCollateralInfo 押品变更历史信息
type HistoryCollateralInfo struct {
	TxID            string           `json:"txid"`
	Time            string           `json:"time"`
	CollateralInfos []CollateralInfo `json:"collateralInfos"`
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (a *AssertsManageCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (a *AssertsManageCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "addCustomerInfo" {
		return a.addCustomerInfo(stub, args)
	} else if fn == "addCollateralInfo" {
		return a.addCollateralInfo(stub, args)
	} else if fn == "addProjectInfo" {
		return a.addProjectInfo(stub, args)
	} else if fn == "getCustomerInfo" {
		return a.getCustomerInfo(stub, args)
	} else if fn == "getHistoryProjectInfo" {
		return a.getHistoryProjectInfo(stub, args)
	} else if fn == "getHistoryCustomerInfo" {
		return a.getHistoryCustomerInfo(stub, args)
	} else if fn == "getHistoryCollateralInfo" {
		return a.getHistoryCollateralInfo(stub, args)
	}
	return shim.Error("Recevied unkown function invocation")
}

// 添加客户信息
func (a *AssertsManageCC) addCustomerInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// var err error
	var CustomerInfo CustomerInfo

	// args[参数列表]
	// 1.检查参数的个数
	if len(args) != 10 {
		// 用shim.Error返回错误信息
		return shim.Error("Incorrect number of arguments.")
	}

	// 2.验证参数的正确性
	// 从参数中获取客户名称
	Name := args[0]
	if Name == "" {
		return shim.Error("CustomerName can not be empty.")
	}

	CustomerInfo.ID = args[1]
	CustomerInfo.Code = args[2]
	CustomerInfo.Type = args[3]
	CustomerInfo.Money = args[4]
	CustomerInfo.Person = args[5]
	CustomerInfo.Date = args[6]
	CustomerInfo.BusinessDate = args[7]
	CustomerInfo.ApprovalDate = args[8]
	CustomerInfo.Trade = args[9]

	// // 3.验证数据是否存在 [应该存在or不应该存在]
	// // 验证需要读取 stateDB，需要 shim 包中的 GetState 方法
	// // stateDB（KV类型）需要定义组合键的方法来区分用户和资产

	// if customerBytes, err := stub.GetState(CustomerInfo.ID); err == nil && len(customerBytes) != 0 {
	// 	return shim.Error("Customer already exist")
	// }

	// 4.状态写入
	// 序列化对象 CustomerInfo
	JSONasBytes, err := json.Marshal(CustomerInfo)
	if err != nil {
		// return shim.Error(err.Error())
		return shim.Error(fmt.Sprintf("marshal error, %s", err))
	}
	// 用 stub.PutState 写入 stateDB（KV类型数据库）
	// PutState 方法， 如果数据不存在，就新增；如果数据存在，就修改。
	if err := stub.PutState(Name+"CustomerInfo", JSONasBytes); err != nil {
		// return shim.Error(err.Error())
		return shim.Error(fmt.Sprintf("put stateDB error, %s", err))
	}

	// 成功返回
	return shim.Success(nil)
}

// 添加押品信息
func (a *AssertsManageCC) addCollateralInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 声明客户和押品结构体
	// var Customer Customer
	var CollateralInfos []CollateralInfo
	var CollateralInfo CollateralInfo

	// 1.检查参数的个数
	// 判断押品信息是否成对出现
	if (len(args)-1)%2 != 0 || len(args) == 1 {
		return shim.Error("Incorrect number of arguments")
	}

	// 2.验证参数的正确性

	// 获取客户姓名
	Name := args[0]
	// Customer.Name = Name

	// 3.验证数据是否存在
	if customerBytes, err := stub.GetState(Name); err != nil && len(customerBytes) == 0 {
		return shim.Error("Customer not found")
	}
	// 批量添加
	for i := 1; i < len(args); i = i+2 {
		CollateralInfo.CollateralID = args[i]
		CollateralInfo.CollateralName = args[i+1]
		CollateralInfos = append(CollateralInfos, CollateralInfo)
	}


	// 4.状态写入
	// 序列化对象
	JSONasBytes, err := json.Marshal(CollateralInfos)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal collateral error, %s", err))
	}
	// 提交，写入 stateDB，此时相当于把用户信息也修改了。
	if err := stub.PutState(Name+"CollateralsInfo", JSONasBytes); err != nil {
		return shim.Error(fmt.Sprintf("put stateDB error, %s", err))
	}
	// 成功返回
	return shim.Success(nil)
}

// 添加项目信息
func (a *AssertsManageCC) addProjectInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// var Customer Customer
	var ProjectInfo ProjectInfo

	// 1.检查参数的个数
	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments.")
	}

	// 2.验证参数的正确性
	Name := args[0]
	if Name == "" {
		return shim.Error("CustomerName can not be empty.")
	}
	// 3.验证数据是否存在 
	if customerBytes, err := stub.GetState(Name); err != nil && len(customerBytes) == 0 {
		return shim.Error("Customer not found")
	}

	ProjectInfo.ProjectName = args[1]
	ProjectInfo.ProjectID = args[2]
	ProjectInfo.ProjectType = args[3]
	ProjectInfo.ProjectTrade = args[4]
	ProjectInfo.ProjectDate = args[5]
	ProjectInfo.ProjectApprove = args[6]
	ProjectInfo.ProjectPart = args[7]
	ProjectInfo.ProjectInvest = args[8]
	ProjectInfo.ProjectMoney = args[9]
	ProjectInfo.ProjectCompanyType = args[10]

	// 4.状态写入
	// 序列化对象
	JSONasBytes, err := json.Marshal(ProjectInfo)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal project error, %s", err))
	}
	// 提交，写入 stateDB
	if err := stub.PutState(Name+"ProjectInfo", JSONasBytes); err != nil {
		return shim.Error(fmt.Sprintf("put stateDB error, %s", err))
	}
	// 成功返回
	return shim.Success(nil)
}

// 获取客户信息
func (a *AssertsManageCC) getCustomerInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 1.检查参数的个数
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	// 2.验证参数的正确性
	Name := args[0]
	if Name == "" {
		return shim.Error("CustomerName can not be empty.")
	}

	// 3.验证数据是否存在 
	if customerBytes, err := stub.GetState(Name); err != nil && len(customerBytes) == 0 {
		return shim.Error("Customer not found")
	}

	CustomerInfoAsBytes, err := stub.GetState(Name + "CustomerInfo")
	if err != nil {
		return shim.Error(err.Error())
	}
	CollateralInfosAsBytes, err := stub.GetState(Name + "CollateralsInfo")
	if err != nil {
		return shim.Error(err.Error())
	}
	ProjectInfoAsBytes, err := stub.GetState(Name + "ProjectInfo")
	if err != nil {
		return shim.Error(err.Error())
	}

	var CustomerInfo CustomerInfo
	var CollateralInfos []CollateralInfo
	var ProjectInfo ProjectInfo
	json.Unmarshal(CustomerInfoAsBytes, &CustomerInfo)
	json.Unmarshal(CollateralInfosAsBytes, &CollateralInfos)
	json.Unmarshal(ProjectInfoAsBytes, &ProjectInfo)
	// 
	// if err != nil {
	// 	return shim.Error(err.Error())
	// 	// return shim.Error(fmt.Sprintf("unmarshal project error, %s", err))
	// }

	var Customer Customer
	Customer.Name = Name
	Customer.CustomerInfo = CustomerInfo
	Customer.CollateralInfo = CollateralInfos
	Customer.ProjectInfo = ProjectInfo

	CustomerAsBytes, err := json.Marshal(Customer)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal project error, %s", err))
		// return shim.Error(err.Error())
	}

	return shim.Success(CustomerAsBytes)
}

// 获取项目的历史数据
func (a *AssertsManageCC) getHistoryProjectInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	Name := args[0]
	
	// 用 GetHistoryForKey 返回一个可以迭代的对象
	resultsIterator, err := stub.GetHistoryForKey(Name + "ProjectInfo")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var HistoryProjectInfos []HistoryProjectInfo
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var HistoryProjectInfo HistoryProjectInfo
		HistoryProjectInfo.TxID = response.TxId
		txtimestamp := response.Timestamp
		// tm := time.Unix(txtimestamp.Seconds+28800, 0)
		// loc, _ := time.LoadLocation("Asia/Beijing")
		// datestr := tm.In(loc).Format("2006-01-02 03:04:05 PM")
		tm := time.Unix(txtimestamp.Seconds, 0)
		datestr := tm.Format("2006-01-02 03:04:05 PM")
		HistoryProjectInfo.Time = datestr
		json.Unmarshal(response.Value, &HistoryProjectInfo.ProjectInfo)
		HistoryProjectInfos = append(HistoryProjectInfos, HistoryProjectInfo)
	}
	jsonsAsBytes, err := json.Marshal(HistoryProjectInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

// 获取客户历史数据
func (a *AssertsManageCC) getHistoryCustomerInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	Name := args[0]
	resultsIterator, err := stub.GetHistoryForKey(Name + "CustomerInfo")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var HistoryCustomerInfos []HistoryCustomerInfo
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var HistoryCustomerInfo HistoryCustomerInfo
		HistoryCustomerInfo.TxID = response.TxId
		txtimestamp := response.Timestamp
		tm := time.Unix(txtimestamp.Seconds, 0)
		datestr := tm.Format("2006-01-02 03:04:05 PM")
		HistoryCustomerInfo.Time = datestr
		json.Unmarshal(response.Value, &HistoryCustomerInfo.CustomerInfo)
		HistoryCustomerInfos = append(HistoryCustomerInfos, HistoryCustomerInfo)
	}
	jsonsAsBytes, err := json.Marshal(HistoryCustomerInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

// 获取押品历史数据
func (a *AssertsManageCC) getHistoryCollateralInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	Name := args[0]
	
	resultsIterator, err := stub.GetHistoryForKey(Name + "CollateralsInfo")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var HistoryCollateralInfos []HistoryCollateralInfo
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var HistoryCollateralInfo HistoryCollateralInfo
		HistoryCollateralInfo.TxID = response.TxId
		txtimestamp := response.Timestamp
		tm := time.Unix(txtimestamp.Seconds, 0)
		datestr := tm.Format("2006-01-02 03:04:05 PM")
		HistoryCollateralInfo.Time = datestr
		var CollateralInfos []CollateralInfo
		json.Unmarshal(response.Value, &CollateralInfos)
		HistoryCollateralInfo.CollateralInfos = CollateralInfos
		HistoryCollateralInfos = append(HistoryCollateralInfos, HistoryCollateralInfo)
	}
	jsonsAsBytes, err := json.Marshal(HistoryCollateralInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

func main() {
	err := shim.Start(new(AssertsManageCC))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s ", err)
	}
}
