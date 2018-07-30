package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"strings"

	auth "app/pb/auth"
	bear "app/pb/bear"
	order "app/pb/order"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	authEndPoint = flag.String("auth_endpoint", "auth_service:7800", "authenticate of AuthService")
	bearEndPoint = flag.String("bear_endpoint","bear_service:7801","bear service")
	orderEndPoint = flag.String("order_endpoint","order_service:7802","order service")
  // adminEndPoint = flag.String("admin_endpoint","127.0.0.1:7807","admin service")
)

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept","Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}

func newGateway(ctx context.Context,opts ...runtime.ServeMuxOption)(http.Handler,error){
		mux := runtime.NewServeMux(opts...)
		dialOpts := []grpc.DialOption{grpc.WithInsecure()}
		err := auth.RegisterAuthHandlerFromEndpoint(ctx, mux, *authEndPoint, dialOpts)
		if err != nil {
			return nil,err
		}
		err = bear.RegisterBearHandlerFromEndpoint(ctx,mux,*bearEndPoint,dialOpts)
		if err != nil{
			return nil,err
		}
		err = order.RegisterOrderHandlerFromEndpoint(ctx,mux,*orderEndPoint,dialOpts)
		if err != nil {
			return nil,err
		}

		// err = admin.RegisterAdminHandlerFromEndpoint(ctx,mux,*adminEndPoint,dialOpts)
		// if err != nil {
		// 	return nil,err
		// }

		return mux,nil
}

func RunEndPoint(address string, opts ...runtime.ServeMuxOption) error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	gw ,err := newGateway(ctx,opts...)
	if err != nil {
		return err
	}

	mux.Handle("/",gw)

	http.ListenAndServe(address, allowCORS(mux))
	return nil
}

func main(){
	flag.Parse()
	defer glog.Flush()

	if err := RunEndPoint(":8080"); err != nil {
		glog.Fatal(err)
	}
}
