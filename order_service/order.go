package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"github.com/dgrijalva/jwt-go"
	"crypto/rsa"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"github.com/appleboy/gorush/rpc/proto"
	"time"
	pb "app/pb"
	common "app/common"
)

type orderServer struct{
		jwtPublicKey *rsa.PublicKey
}
const (
	 address = "127.0.0.1:9000"
)

func NewOrderServer(rsaPublicKey string) (*orderServer,error){
	data,err := ioutil.ReadFile(rsaPublicKey)
	if err != nil {
			return nil, fmt.Errorf("Error reading the jwt public key: %v ",err)
	}

	publicKey,err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		 return nil, fmt.Errorf("Error parsing the jwt public key: %v",err)
	}
	return &orderServer{publicKey},nil
}

func (os *orderServer) GetListOrders(ctx context.Context, request *pb.ListOrderRequest) (*pb.ListOrderResponse,error){
	_,err := checkToken(ctx,os.jwtPublicKey)
	if err != nil {
		return nil,err
	}

	dateTimestamp := request.Date

	if dateTimestamp == 0 {
			dateTimestamp = uint64(time.Now().Unix())
	}

	t := time.Unix(int64(dateTimestamp),0)

	listOrderResponse := &pb.ListOrderResponse{}
	listOrderResponse.Orders,err = getOrders(t.Format("2006-01-02"),request.UserId)
	if err != nil {
		return nil,err
	}
	return listOrderResponse,nil
}

func (os *orderServer) OrderBlock(ctx context.Context,request *pb.OrderBlockRequest) (*pb.OrderBlockResponse,error){
	_,err := checkToken(ctx,os.jwtPublicKey)
	if err != nil {
		return nil,err
	}
	if request.UserId == 0 {
		return nil,grpc.Errorf(codes.InvalidArgument,"Invalid argument")

	}

	if request.BlockId == 0 {
		return nil,grpc.Errorf(codes.InvalidArgument,"Invalid argument")
	}
	inserted,err := orderBlockFromUser(request.BlockId,request.UserId)
	if err != nil {
		return nil,err
	}

	var orderResponse pb.OrderBlockResponse

	if inserted == true {
			orderResponse = pb.OrderBlockResponse{Message : "Order sucessful", Status : "OK"}
	}else {
		  orderResponse = pb.OrderBlockResponse{Message : "Order failed, something happened", Status : "Failed"}
	}
	return &orderResponse,nil
}

func (os *orderServer) ConfirmOrder(ctx context.Context,request *pb.ConfirmOrderRequest) (*pb.ConfirmOrderResponse,error){
	_,err := checkToken(ctx,os.jwtPublicKey)
	if err != nil {
		return nil,err
	}
	if request.OrderId == 0 {
		return nil,grpc.Errorf(codes.InvalidArgument,"Invalid argument")
	}
	confirmOrderResp := &pb.ConfirmOrderResponse{}
	var status bool
	confirmOrderResp.BlockOrder,confirmOrderResp.Message,status,err = confirmOrder(request.OrderId)
	if err != nil {
			return nil, err
	}
	if status == true {
		confirmOrderResp.Status = "OK"
	}else {
		confirmOrderResp.Status = "Failed"
	}

	return confirmOrderResp,nil

}


func (os *orderServer) FinishService(ctx context.Context,request *pb.FinishServiceRequest)(*pb.FinishServiceResponse,error){
	_,err := checkToken(ctx,os.jwtPublicKey)
	if err != nil {
			return nil, err
	}
	if request.OrderBlockId == 0 {
		return nil,grpc.Errorf(codes.InvalidArgument,"Invalid argument")
	}

	finishServiceRes := &pb.FinishServiceResponse{}
	finishServiceRes.Message,finishServiceRes.Status,err = finishOrder(request.OrderBlockId)
	if err != nil {
		return nil,err
	}
	return finishServiceRes,nil
}

func finishOrder(orderId uint32)(string,string,error){
	db,err := common.MysqlConnection()
	if err != nil {
		return "Error open database","Failed",err.Err
	}
	defer db.Close()
	var orderBlock pb.BlockOrder
	var userID uint32
	var userBlockDateID uint32

	db.Raw("Select * from block_orders WHERE id IN (?) AND status = 'CONFIRM'",orderId).Row().Scan(&orderBlock.Id,&userID,&userBlockDateID,&orderBlock.Status)
	if orderBlock.Id == 0 {
		return "Error processing","Failed",grpc.Errorf(codes.NotFound,"This block not found")
	}

	tx := db.Begin()
	defer func(){
		if r := recover();r != nil{
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return "Error processing","Failed",err.Err
	}

	db.Raw("UPDATE block_orders SET status = 'FINISH' WHERE id IN (?)", orderId).Row()
	commitErr := tx.Commit().Error
	if err != nil {
		return "Error processing","Failed",commitErr
	}

	var userTokenString string
	db.Table("block_orders AS bo").Select("users.device_token").Joins("INNER JOIN user_block_dates AS ubd ON ubd.id = bo.user_block_date_id").Joins("INNER JOIN users ON users.id = ubd.user_id").Where("bo.id IN (?)",orderId).Row().Scan(&userTokenString)

	sendNotification(userTokenString,"You order have finished new order",fmt.Sprintf("Your %d have finished",orderId),fmt.Sprintf("%d order fnished",orderId),"finish order")

	return "Your service is finished","OK",nil
}

func confirmOrder(orderId uint32) (*pb.BlockOrder,string,bool,error){
	db,err := common.MysqlConnection()
	if err != nil {
		return nil,"Error processing",false,err.Err
	}

	defer db.Close()

	var orderBlock pb.BlockOrder
	var userID uint32
	var userBlockDateID uint32

	db.Raw("Select * from block_orders WHERE id IN (?) AND status = 'WAITING'",orderId).Row().Scan(&orderBlock.Id,&userID,&userBlockDateID,&orderBlock.Status)
	if orderBlock.Id == 0 {
		return nil,"Error processing",false,grpc.Errorf(codes.NotFound,"This block not found")
	}

	if userID == 0 || userBlockDateID == 0 {
		  return nil,"Error processing",false,grpc.Errorf(codes.NotFound,"Error argument , no user or block specified ")
	}



	tx := db.Begin()
	defer func(){
		if r := recover();r != nil{
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil,"Error processing",false,err.Err
	}

	db.Raw("UPDATE block_orders SET status = 'CANCEL' WHERE user_block_date_id IN (?)",userBlockDateID)

	db.Raw("UPDATE block_orders SET status = 'CONFIRM' WHERE id IN (?)", orderId).Row()

	if err := db.Exec("UPDATE user_block_dates SET status = 'BUSY' WHERE id IN (?)",userBlockDateID).Error; err != nil{
		tx.Rollback()
		return nil,"Error processing",false,err
	}


	commitErr := tx.Commit().Error
	if err != nil {
		return nil,"Error processing",false,commitErr
	}
	var userPush pb.User
	db.Where("id = ?",userID).First(&userPush)
	sendNotification(userPush.DeviceToken,"You order have finished new order",fmt.Sprintf("Your %d have finished",orderId),fmt.Sprintf("%d order fnished",orderId),"finish order")
	return &orderBlock,"Your order confirm is success",true,nil

}

func orderBlockFromUser(blockId uint32,userId uint32)(bool,error){
		db,err := common.MysqlConnection()
		if err != nil {
			return false,err.Err
		}

		defer db.Close()
		// order := &pb.BlockOrder{}

		var user pb.User
		db.First(&user)
		if user.Id <= 0 {
			return false, grpc.Errorf(codes.NotFound,"User not found")
		}
		var userBlockDate pb.UserBlockDate
		db.First(&userBlockDate,blockId)
		if userBlockDate.Id == 0 {
				return false,grpc.Errorf(codes.NotFound,"This block not found")
		}
		var orderCheck uint32
		row := db.Table("block_orders").Select("user_id").Where("user_id = ? AND user_block_date_id = ?",userId,blockId).Row()
		row.Scan(&orderCheck)

		if orderCheck != 0 {
			return false,grpc.Errorf(codes.AlreadyExists,"You already ordered this block")
		}

		tx := db.Begin()
		defer func(){
			if r := recover();r != nil{
				tx.Rollback()
			}
		}()

		if tx.Error != nil {
			return false,err.Err
		}

		if err := db.Exec("INSERT INTO block_orders(user_id,user_block_date_id,status) VALUES (?,?,'WAITING')",userId,blockId).Error; err != nil{
			tx.Rollback()
			return false,err
		}

		commitErr := tx.Commit().Error
		if err != nil {
			return false,commitErr
		}

		var userNotif pb.User
		db.Where("user_id = ?",&userNotif)

		sendNotification(user.DeviceToken,"You have new order","New order","You have new order from user","New order")
		return true,nil
}
//219 xvnt
func sendNotification(deviceToken string,opts...string){
			// deviceToken := "dtTS_CSTYW4:APA91bF7kFqGD6E2MKFQfz7aE3wcxPiewv28ElW6NXIPQAZ_gWTpg7Bg0qOecPRy0mOfGUGTq5mh3Jq_SCWd6S-eJB5WzbL6ZFlJA8TPh5t2eR2BekYFpPB3S74CA2lkOVJQ9QMD4ya5"
			conn,err := grpc.Dial(address,grpc.WithInsecure())
			if err != nil {
				  log.Printf("cannot connect to push notification service: %v",err)
			}
			defer conn.Close()
			c := proto.NewGorushClient(conn)
			r,err := c.Send(context.Background(),&proto.NotificationRequest{
						Platform: 2,
						Tokens: []string{deviceToken},
						Message: opts[0],
						Alert: &proto.Alert{
							Title:    opts[1],
							Body:     opts[2],
							Subtitle: opts[3],
						},
			})

			if err != nil {
					log.Printf("cannot push %v",err)
			}else{
				log.Printf("Success: %t\n", r.Success)
				log.Printf("Count: %d\n", r.Counts)
			}
}
func getOrders(date string,userId uint32)([]*pb.BlockOrder,error){
		db,err := common.MysqlConnection()
		if err != nil {
			return nil ,err.Err
		}

		defer db.Close()

		var orders []*pb.BlockOrder

		rows,errQuery := db.Table("block_orders as bo").Select("bo.id,users.username,user_block_dates.id,blocks.name,blocks.description").Joins("LEFT JOIN user_block_dates ON user_block_dates.id = bo.user_block_date_id").Joins("LEFT JOIN users ON users.id = bo.user_id").Joins("LEFT JOIN blocks ON blocks.id = user_block_dates.block_id").Where("user_block_dates.user_id = ? AND user_block_dates.block_date = ? AND bo.status = 'WAITING'",userId,date).Rows()
		if errQuery != nil {
			return nil,errQuery
		}

		for rows.Next() {
			var order pb.BlockOrder
			var user pb.User
			var ubd pb.UserBlockDate
			var block pb.Block

			err := rows.Scan(&order.Id,&user.Username,&ubd.Id,&block.Name,&block.Description)
			if err != nil {
					return nil,err
			}

			ubd.Block = &block
			order.UserBlockDate = &ubd
			order.User = &user

			orders = append(orders,&order)

		}
		// sendNotification()
		return orders,nil
}

func checkToken(ctx context.Context, jwtPublicKey *rsa.PublicKey) (bool, error){

		// var token *jwt.Token
		// var err error

		md,ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return false,grpc.Errorf(codes.PermissionDenied,"Missing token")
		}

		jwtToken,ok := md["authorization"]
		if !ok {
				return false, grpc.Errorf(codes.PermissionDenied,"Missing token")
		}

		_ ,err := validateToken(jwtToken[0],jwtPublicKey)
		if err != nil {
				return false, grpc.Errorf(codes.PermissionDenied,"Invalid token")
		}

		return true,nil

}

func validateToken(token string,publickey *rsa.PublicKey) (*jwt.Token,error){
		jwtToken, err := jwt.Parse(token,func(t *jwt.Token)(interface{},error){
				if _,ok := t.Method.(*jwt.SigningMethodRSA); !ok{
						log.Printf("Wrong signing method : %v",t.Header["alg"])
						return nil,fmt.Errorf("Invalid token")
				}
				return publickey,nil
		})

		if err == nil && jwtToken.Valid {
			 return jwtToken,nil
		}
		return nil,err
}

