package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	pb "./pb"

	jwt "github.com/dgrijalva/jwt-go"
	common	"./common"

	// "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
	"time"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type adminServer struct {
	jwtPublicKey *rsa.PublicKey
	jwtPrivatekey *rsa.PrivateKey
}
func NewAdminServer(rsaPublicKey string,rsaPrivateKey string) (*adminServer,error){
	data,err := ioutil.ReadFile(rsaPublicKey)
	if err != nil {
			return nil, fmt.Errorf("Error reading the jwt public key: %v ",err)
	}

	publicKey,err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		 return nil, fmt.Errorf("Error parsing the jwt public key: %v",err)
	}

	key, err := ioutil.ReadFile(rsaPrivateKey)
	if err != nil {
			return nil, fmt.Errorf("Error reading the jwt private key",err)
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil{
			return nil, fmt.Errorf("Error parsing the jwt private key")
	}

	return &adminServer{publicKey,parsedKey},nil
}

func (as *adminServer) GetListUser(ctx context.Context,msg *empty.Empty) (*pb.ListUsersResponse,error){
	_,err := checkToken(ctx,as.jwtPublicKey)
	if err != nil {
		return nil,err
	}
	userListResponse := &pb.ListUsersResponse{}

	userListResponse.Users,userListResponse.Bears,err = getUser()
	if err != nil {
		 return nil,err
	}
	return userListResponse,nil
}

func (as *adminServer) GetListBlocks(ctx context.Context,msg *empty.Empty) (*pb.ListBlocksResponse,error){
	_,err := checkToken(ctx,as.jwtPublicKey)
	if err != nil {
		return nil,err
	}
	blockListResponse := &pb.ListBlocksResponse{}
	blockListResponse.Blocks,err = getBlocks()
	if err != nil {
		return nil,err
	}
	return blockListResponse,nil
}

func (as *adminServer) GetListReviews(ctx context.Context,msg *empty.Empty) (*pb.ListReviewsResponse,error){
	_,err := checkToken(ctx,as.jwtPublicKey)
	if err != nil {
		return nil,err
	}
	reviewListResponse := &pb.ListReviewsResponse{}
	reviewListResponse.Reviews,err = getReviews()
	if err != nil {
		return nil,err
	}
	return reviewListResponse,nil
}

func getReviews()([]*pb.Review,error){
		db,err := common.MysqlConnection()
		if err != nil {
			return nil,err.Err
		}
		defer db.Close()
		var reviewList []*pb.Review
		rows,errSelect := db.Table("reviews r").Select("r.id,u.username,u.email,i.username,i.email,r.rate,r.description").Joins("LEFT JOIN users u ON u.id = r.user_reviewed ").Joins("LEFT JOIN users i ON i.id = r.user_id").Rows()
		if errSelect != nil {
			return nil, errSelect
		}

		for rows.Next(){
				var review pb.Review
				var user pb.User
				var bear pb.User
				err := rows.Scan(&review.Id,&user.Username,&user.Email,&bear.Username,&bear.Email,&review.Rate,&review.Description)
				if err != nil {
					fmt.Println(err)
				}
				review.User = &bear
				review.UserReviewed = &user
				reviewList = append(reviewList,&review)
		}
		return reviewList,nil
}

func getBlocks()([]*pb.Block,error){
		db,err := common.MysqlConnection()
		if err != nil {
			 return nil,err.Err
		}
		defer db.Close()
		var blockList []*pb.Block
		rows,errSelect := db.Table("blocks").Select("id,description,name,hour_start,hour_end").Rows()
		if errSelect != nil {
			return nil,errSelect
		}
		for rows.Next(){
				var block pb.Block

				err := rows.Scan(&block.Id,&block.Description,&block.Name,&block.HourStart,&block.HourEnd)
				if err != nil {
					fmt.Println(err)
				}
				blockList = append(blockList,&block)
		}

		return blockList,nil
}

func getUser()([]*pb.User,[]*pb.User,error){
	db,err := common.MysqlConnection()
	if err != nil {
		return nil,nil,err.Err
	}

	defer db.Close()
	var userList []*pb.User
	var bearList []*pb.User
	rows,errSelect := db.Table("users").Select("users.id,users.password_hash,users.created_at,users.updated_at,users.phone,users.role_id,users.email,users.username,profiles.avatar,profiles.birthdate,profiles.firstname,profiles.lastname,profiles.sex,profiles.description,profiles.star_rate").Joins("LEFT JOIN profiles ON profiles.user_id = users.id").Rows()
	if errSelect != nil {
		return nil,nil , errSelect
 	}

 for rows.Next() {
		 var user pb.User
		 profile := &pb.Profile{}

		 err := rows.Scan(&user.Id,&user.PasswordHash,&user.CreatedAt,&user.UpdatedAt,&user.Phone,&user.RoleId,&user.Email,&user.Username,&profile.Avatar,&profile.Birthdate,&profile.Firstname,&profile.Lastname,&profile.Sex,&profile.Description,&profile.StarRate)
		 if err != nil {
			 fmt.Println(err)
		 }
		 user.Profile = profile
		 if user.RoleId == 1 {
			  userList = append(userList,&user)
		 }else {
				bearList = append(bearList,&user)
		 }
 }

 return userList,bearList,nil
}


func (as *adminServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse,error){

	user,err := getAdmin(request.Username)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument,"Username not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(request.Password))

	if err != nil {
		 return nil,  grpc.Errorf(codes.InvalidArgument,"Wrong username or password")
	}

	if user.RoleId != 3 {
		return nil, grpc.Errorf(codes.PermissionDenied,"No permission")
	}
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["admin"] = true
	claims["iss"] = "auth.service"
	claims["iat"] = time.Now().Unix()
	claims["email"] = user.Email
	claims["sub"] = user.Username
	token.Claims = claims
	tokenString,err := token.SignedString(as.jwtPrivatekey)
	if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated,err.Error())
	}
	user.PasswordHash = ""
	return &pb.LoginResponse{tokenString,user},nil
}

func getAdmin(username string) (*pb.User,error){
		// user := &pb.User{}
		var user []*pb.User
		// profile := &pb.Profile{}
		db, errDB := common.MysqlConnection()
		if errDB != nil {
				return nil,errDB.Err
		}
		defer db.Close()

		if err := db.Where("username = ?",username).First(&user).Error; err != nil{
				return nil,err
		}

		if len(user) < 1 {
			return nil, grpc.Errorf(codes.NotFound,"User not found")
		}

		userDetail := user[0]

		profile := &pb.Profile{}

		if err := db.Where("user_id = ?",userDetail.Id).First(&profile).Error; err != nil{
			return nil,err
		}

		userDetail.Profile = profile
		// if err := db.Where("user_id = ?",user.Id).First(&profile).Error; err != nil{
		// 	return nil ,err
		// }
		// user.Profile = profile
		return userDetail,nil
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
			return false, err
	}

	return true,nil

}

func validateToken(token string,publickey *rsa.PublicKey) (*jwt.Token,error){
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token,&claims,func(t *jwt.Token)(interface{},error){
			if _,ok := t.Method.(*jwt.SigningMethodRSA); !ok{
					log.Printf("Wrong signing method : %v",t.Header["alg"])
					return nil,fmt.Errorf("Invalid token")
			}
			return publickey,nil
	})


	if err == nil && jwtToken.Valid {
			permission := claims["admin"]
			if permission != nil {
				return jwtToken,nil
			}else{
				return nil,fmt.Errorf("Invalid permission")
			}
	}

	return nil,err
}
