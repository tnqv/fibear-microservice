package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"time"
	"github.com/dgrijalva/jwt-go"
	"crypto/rsa"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"errors"
	"github.com/appleboy/gorush/rpc/proto"
	pb "app/pb"
	common "app/common"
)

type bearServer struct{
		jwtPublicKey *rsa.PublicKey
}

var (
	address = "127.0.0.1:9000"
)

func NewBearServer(rsaPublicKey string) (*bearServer,error){
	data,err := ioutil.ReadFile(rsaPublicKey)
	if err != nil {
			return nil, fmt.Errorf("Error reading the jwt public key: %v ",err)
	}

	publicKey,err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		 return nil, fmt.Errorf("Error parsing the jwt public key: %v",err)
	}
	return &bearServer{publicKey},nil
}

func (bs *bearServer) GetBearBlocks(ctx context.Context, request *pb.BearBlockRequest) (*pb.BearBlockResponse,error){
				_,err := checkToken(ctx,bs.jwtPublicKey)
				if err != nil {
					return nil,err
				}

				bearBlocks := &pb.BearBlockResponse{}

				fmt.Println(request.Date)
				dateTimestamp := request.Date
				if dateTimestamp == 0 {
						dateTimestamp = uint64(time.Now().Unix())
				}
				fmt.Println(dateTimestamp)
				t := time.Unix(int64(dateTimestamp),0)
				fmt.Println(t)
				if request.UserId <= 0 {
					return nil,grpc.Errorf(codes.InvalidArgument,"No user id specified")
				}
				bearBlocks.UserBlockDates,err = getBearBlock(t.Format("2006-01-02"),request.BearId,request.UserId)
				if err != nil {
					return nil,err
				}
				return bearBlocks,nil
}

func (bs *bearServer) GetBlocksToAssign(ctx context.Context,request *pb.GetBlockRegisterRequest)(*pb.GetBlockRegisterResponse,error){
	_,err := checkToken(ctx,bs.jwtPublicKey)
	if err != nil {
		return nil,err
	}
	dateTimestamp := request.Date
	if dateTimestamp == 0 {
			dateTimestamp = uint64(time.Now().Unix())
	}
	t := time.Unix(int64(dateTimestamp),0)
	blockRegisterResponse := &pb.GetBlockRegisterResponse{}
	blockRegisterResponse.UserBlockDates,err = getBlockForAssign(request.UserId,t.Format("2006-01-02"))
	if err != nil {
		return nil,err
	}

	return blockRegisterResponse,nil
}

func (bs *bearServer) GetBearDetail(ctx context.Context, request *pb.BearRequest) (*pb.BearResponse,error){
		_,err := checkToken(ctx,bs.jwtPublicKey)
		if err != nil {
				return nil,err
		}

		bear := &pb.BearResponse{}
		bear.User,bear.Reviews,err = getBear(request.UserId)
		if err != nil{
			return nil,err
		}
		return bear,nil
}

func (bs *bearServer) GetListBear(ctx context.Context, request *pb.ListBearRequest)(*pb.ListBearResponse,error){
		_,err := checkToken(ctx,bs.jwtPublicKey)
		if err != nil {
				return nil,err
		}
		listBearResponse := &pb.ListBearResponse{}
		listBearResponse.Users,err = getListBear(request.City)
		if err != nil {
			return nil ,err
		}

		return listBearResponse,nil
}

func (bs *bearServer) CreateBlocks(ctx context.Context,request *pb.CreateBlockRequest)(*pb.CreateBlockResponse,error){
		_,err := checkToken(ctx,bs.jwtPublicKey)
		if err != nil {
				return nil,err
		}
		createBlockResponse := &pb.CreateBlockResponse{}
		createBlockResponse.Block,err = createBlock(request.BearId,request.UserBlockDate)
		if err != nil {
			return nil, err
		}
		createBlockResponse.Status = "success"
		return createBlockResponse,nil
}

func (bs *bearServer) PostReview(ctx context.Context,request *pb.PostReviewRequest)(*pb.PostReviewResponse,error){
	 _,err := checkToken(ctx,bs.jwtPublicKey)
	 if err != nil {
		 return nil,err
	 }
	 postReviewResponse := &pb.PostReviewResponse{}
	 postReviewResponse.Status , postReviewResponse.Message,err = postReview(request)
	 if err != nil {
		 return nil,err
	 }

	 return postReviewResponse,nil

}

func postReview(request *pb.PostReviewRequest)(string,string,error){
		db,err := common.MysqlConnection()
		if err != nil {
			return "Failed","Error open database", err.Err
		}
		defer db.Close()

		var review pb.Review

		db.Where("user_reviewed = ? AND user_id = ?",request.UserId,request.BearId).First(&review)
		if review.Id > 0 {
			return "Failed","Duplicated review",grpc.Errorf(codes.AlreadyExists,"Duplicated review")
		}
		review.UserReviewed.Id = request.UserId
		review.UserId = uint64(request.BearId)
		review.Rate = uint64(request.Rate)
		review.Description = request.Description


		db.Raw("INSERT INTO reviews(user_reviewed,user_id,rate,description) VALUES (?,?,?,?)",review.UserReviewed.Id,review.UserId,review.Rate,review.Description).Row()


		var bear pb.User
		db.Where("id = ?",review.UserId).First(&bear)

		sendNotification(bear.DeviceToken,"You received new Review","New review","You have received new review from user","Review Received")
		return "Success","Your review is created",nil
}

func getBlockForAssign(bearID uint32, date string)([]*pb.UserBlockDate,error){
	db,err := common.MysqlConnection()
	if err != nil {
		return nil,err.Err
	}
	defer db.Close()
	var listBlocks []*pb.Block
	var userBlockDates []*pb.UserBlockDate

	db.Find(&listBlocks)

	rows,errQuery := db.Table("user_block_dates as ubd").Select("ubd.id,ubd.price,ubd.description,ubd.status,ubd.block_date,ubd.block_id").Where("user_id = ? AND block_date = ?",bearID,date).Rows()
	if errQuery != nil {
		return nil,errQuery
	}

	for rows.Next() {
		userBlockDate := &pb.UserBlockDate{}
		// var blockId uint32

		err := rows.Scan(&userBlockDate.Id,&userBlockDate.Price,&userBlockDate.Description,&userBlockDate.Status,&userBlockDate.BlockDate,&userBlockDate.BlockId)
		if err != nil {
				return nil,err
		}
		userBlockDates = append(userBlockDates,userBlockDate)
	}

	var userBlockDatesReturn []*pb.UserBlockDate

	for  _,block := range listBlocks {
				blockExisted := false
				for _,value := range userBlockDates {

						if block.Id == value.BlockId {
								value.Block = block
								userBlockDatesReturn = append(userBlockDatesReturn,value)
								blockExisted = true
						}
				}
				temp := &pb.UserBlockDate{Block : block}
				if blockExisted == false {
						userBlockDatesReturn = append(userBlockDatesReturn,temp)
				}
	}
	// for _,value := range listBlocks {
	// 	for _,ubd := range userBlockDates {
	// 				ubd.Block = value
	// 	}


	// 	// if blockId == value.Id {

	// 	// }else {
	// 	// 		userBlockDate.Block = value
	// 	// 		userBlockDates = append(userBlockDates,userBlockDate)
	// 	// }
	// }
	return userBlockDatesReturn,nil

}

func createBlock(bearID uint32,params []*pb.UserBlockDate)([]*pb.UserBlockDate,error){
	db,err := common.MysqlConnection()
	if err != nil {
		return nil,err.Err
	}
	defer db.Close()
	var listUserBlockDate []*pb.UserBlockDate
	if len(params) <= 0 {
		return nil , grpc.Errorf(codes.InvalidArgument,"No block defined")
	}

	for _,value := range params {
			if value.Block.Id >= 0 && value.Block.Id < 6 && value.BlockDate != ""{
				userBlockDate := &pb.UserBlockDate{}

				db.Where("user_id = ? AND block_id = ? AND block_date = ?",bearID,value.Block.Id,value.BlockDate).First(&userBlockDate)

				if userBlockDate.Id > 0 {

				}else{
					db.Raw("INSERT INTO user_block_dates(user_id,block_id,block_date,description,status,price) VALUES (?,?,?,?,?,?)",bearID,value.Block.Id,value.BlockDate,value.Description,"FREE",value.Price).Row()
					userBlockDate.BlockId = value.Block.Id
					userBlockDate.Description = value.Description
					userBlockDate.Status = "FREE"
					userBlockDate.Price = value.Price
					listUserBlockDate = append(listUserBlockDate,userBlockDate)
				}

			}
	}
	return listUserBlockDate,nil
}
func getBearBlock(date string,bearID uint32,userID uint32)([]*pb.UserBlockDate,error){
	db,err := common.MysqlConnection()
	if err != nil{
		return nil,err.Err
	}
	defer db.Close()

	if bearID == 0 {
		return nil, grpc.Errorf(codes.NotFound,"Invalid user id")
	}
	var blocksOfUser []*pb.UserBlockDate
	rows,errQuery := db.Table("user_block_dates as ubd").Select("ubd.id,ubd.description,ubd.status,ubd.price,ubd.block_date,blocks.description,blocks.name,blocks.hour_start,blocks.hour_end").Joins("LEFT JOIN blocks ON blocks.id = ubd.block_id").Where("user_id = ? AND block_date = ?",bearID,date).Rows()
	if errQuery != nil {
		return nil,errQuery
	}



	for rows.Next() {
		var userBlockDate pb.UserBlockDate
		var block pb.Block

		err := rows.Scan(&userBlockDate.Id,&userBlockDate.Description,&userBlockDate.Status,&userBlockDate.Price,&userBlockDate.BlockDate,&block.Description,&block.Name,&block.HourStart,&block.HourEnd)

		if err != nil {
				return nil,err
		}

		var order pb.BlockOrder
		db.Where("user_id = ? AND user_block_date_id = ?",userID,userBlockDate.Id).First(&order)
		if order.Id > 0 {
				userBlockDate.IsOrdered = true
		}
		userBlockDate.Block = &block
		blocksOfUser = append(blocksOfUser,&userBlockDate)

	}
	return blocksOfUser,nil
}

func getBear(userID uint32)(*pb.User,[]*pb.Review,error){
	db,err := common.MysqlConnection()
	if err != nil {
		return nil,nil,err.Err
	}
	defer db.Close()
	var user []*pb.User

	if err := db.Where("id = ?",userID).First(&user).Error; err != nil{
			return nil,nil,err
	}

	if len(user) < 1 {
		return nil,nil, errors.New("User not found")
	}

	userDetail := user[0]
	userDetail.PasswordHash = ""
	profile := &pb.Profile{}

	if err := db.Where("user_id = ?",userID).First(&profile).Error; err != nil{
		return nil,nil,err
	}

	userDetail.Profile = profile

	// userDetail.Profile = profile
	// profile := userDetail.Profiles
	// if len(profile) <= 0 {
	// 	profile = []*pb.Profile{}
	// }else {

	// }

	var reviews []*pb.Review
	rows,errSelect := db.Table("reviews").Select("reviews.rate,reviews.description,users.username,users.email").Joins("LEFT JOIN users ON reviews.user_reviewed = users.id").Where("reviews.user_id = ?",userDetail.Id).Rows()

	if errSelect != nil {
		 return nil , nil , errSelect
	}

	for rows.Next() {
			var review pb.Review
			var user pb.User

			err := rows.Scan(&review.Rate,&review.Description,&user.Username,&user.Email)
			if err != nil {
				return nil, nil , err
			}
			review.UserReviewed = &user
			reviews = append(reviews,&review)
	}

	return user[0],reviews,nil
}

func getListBear(city uint32) ([]*pb.User, error){
			db,err := common.MysqlConnection()
			if err != nil {
					return nil, err.Err
			}
			defer db.Close()

			var users []*pb.User
		  // if rows,err :=	db.Joins("LEFT JOIN profiles on users.id = profiles.user_id ").Where("profiles.province_id = ?",city).Find(&user).Error;err != nil {
			// 		log.Fatal(err)
			// }
			rows,errSelect := db.Table("users").Select("users.id,users.email,users.username,profiles.avatar,profiles.birthdate,profiles.firstname,profiles.lastname,profiles.sex,profiles.description,profiles.star_rate").Joins("LEFT JOIN profiles ON profiles.user_id = users.id").Where("users.role_id = 2").Rows()

			if errSelect != nil {
				 return nil , errSelect
			}

			for rows.Next() {
					var user pb.User
					profile := &pb.Profile{}

					err := rows.Scan(&user.Id,&user.Email,&user.Username,&profile.Avatar,&profile.Birthdate,&profile.Firstname,&profile.Lastname,&profile.Sex,&profile.Description,&profile.StarRate)
					if err != nil {
						fmt.Println(err)
					}
					user.Profile = profile
					users = append(users,&user)
			}

			return users,nil

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