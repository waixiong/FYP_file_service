package commons

import (
	"context"
	"fmt"

	//pb "./proto"

	"getitqec.com/server/file/pkg/logger"
	"google.golang.org/grpc/metadata"
	// "google.golang.org/grpc/metadata"
)

// logger is to mock a sophisticated logging system. To simplify the example, we just print out the content.
// func logger(format string, a ...interface{}) {
// 	fmt.Printf("LOG:\t"+format+"\n", a...)
// }

var (
//port = flag.Int("port", 50051, "the port to serve on")

)

func GetUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	// fmt.Printf("\tAuth Func Get metadata...\n")
	if !ok {
		logger.Log.Error(fmt.Sprintf("Get User ID from metadata fail"))
		return "", ErrMissingMetadata
	}
	users := md.Get("User")
	if len(users) != 1 {
		logger.Log.Error(fmt.Sprintf("Get User ID from metadata fail"))
		return "", ErrMissingMetadata
	}
	return users[0], nil
}

const (
	CatalogTable        = "Catalog"
	ItemCollection      = "Item"
	InventoryCollection = "Inventory"
)
