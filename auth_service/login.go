package main

import (
	"crypto/rsa"
	"time"
	"fmt"
	"io/ioutil"
	pb "app/pb"

	jwt "github.com/dgrijalva/jwt-go"
	common	"app/common"

	// "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct {
	jwtPrivatekey *rsa.PrivateKey
}

func NewAuthServer(rsaPrivateKey string) (*authServer, error){
	key, err := ioutil.ReadFile(rsaPrivateKey)
	if err != nil {
			return nil, fmt.Errorf("Error reading the jwt private key",err)
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil{
			return nil, fmt.Errorf("Error parsing the jwt private key")
	}
	return &authServer{parsedKey},nil
}

func HashPassword(password string)(string , error){
		bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
		return string(bytes),err
}

func (as *authServer) RefreshToken(ctx context.Context,request *pb.RefreshTokenRequest)(*pb.RefreshTokenResponse,error){
			if request.Token == "" {
				return nil,grpc.Errorf(codes.InvalidArgument,"No Token specified")
			}

			if request.UserId == 0 {
				return nil,grpc.Errorf(codes.InvalidArgument,"No user id specified")
			}
			response := &pb.RefreshTokenResponse{}
			var err error
			response.Message,response.Status,err = updateToken(request.Token,request.UserId)
			if err != nil {
				return nil,err
			}

			return response,nil
}

func (as *authServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse,error){

			user,err := getUser(request.Username)
			if err != nil {
				return nil, grpc.Errorf(codes.InvalidArgument,"Username not found")
			}
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(request.Password))

			if err != nil {
				 return nil,  grpc.Errorf(codes.InvalidArgument,"Wrong username or password")
			}

			token := jwt.New(jwt.SigningMethodRS256)
			claims := make(jwt.MapClaims)

			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
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

func updateToken(token string,userID uint32) (string,string,error){
				db, errDB := common.MysqlConnection()
				if errDB != nil {
						return "","",errDB.Err
				}
				defer db.Close()
				var user pb.User
				db.First(&user,userID)
				if user.Id <= 0 {
					 return "","",grpc.Errorf(codes.NotFound,"User not found")
				}
				user.DeviceToken = token
				if err := db.Save(&user).Error; err != nil {
					return "","",err
				}

				return "Updated device token","success",nil

}

func getUser(username string) (*pb.User,error){
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