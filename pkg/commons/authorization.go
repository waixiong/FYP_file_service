package commons

var (
	services = make(map[string]bool)

	grpcMethod = make(map[string]bool)
)

func InitServicesAuthorization() {
	services["auth_service"] = true
	services["catalog_service"] = true
	// services["delivery_service"] = true
	services["file_service"] = true
	services["mailnotification_service"] = true
	services["map_service"] = true
	services["order_service"] = true
	// services["payment_service"] = true
	services["rider_service"] = true
	services["store_service"] = true
	services["user_service"] = true

	services["sudoadmin"] = true

	// base on service
	grpcMethod["/userproto.UserService/signIn"] = true
	grpcMethod["/userproto.UserService/setStock"] = true
	grpcMethod["/userproto.UserService/queryStock"] = true
}

func IsGetItService(name string) bool {
	// a, ok := services[name]
	if _, ok := services[name]; ok {
		return true
	}
	return false
}

func IsAuthException(path string) bool {
	// a, ok := services[name]
	if _, ok := grpcMethod[path]; ok {
		return true
	}
	return false
}
