package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"getitqec.com/server/file/pkg/commons"
	"getitqec.com/server/file/pkg/logger"

	// "getitqec.com/server/user/pkg/handlers"
	"getitqec.com/server/file/pkg/model"
	"getitqec.com/server/file/pkg/protocol/grpc"
	"getitqec.com/server/file/pkg/protocol/grpcClient"
	"getitqec.com/server/file/pkg/protocol/rest"
	service "getitqec.com/server/file/pkg/service/v1"
	// "getitqec.com/server/catalog/pkg/protocol/rest"
)

var (
	tls      = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert", "configs/key/certs/mycert.pem", "The TLS cert file")
	keyFile  = flag.String("key", "configs/key/private/mykey.pem", "The TLS key file")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 8090, "The server port")
	//svc  = &dynamodb.DynamoDB{}
)

// Config is configuration for Server
type Config struct {
	GRPCPort string
	HTTPPort string
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	// cfg := &Config{GRPCPort: "8090", HTTPPort: "8091", LogLevel: -1}
	cfg := &Config{GRPCPort: commons.ENVVariable("GRPC_SERVICE_PORT"), HTTPPort: commons.ENVVariable("REST_SERVICE_PORT"), LogLevel: -1}

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// initialize model
	// dynamo := commons.GetDB()
	// create table if empty
	// e := dynamo.CreateTable(
	// 	"User",
	// 	[]*dynamodb.AttributeDefinition{
	// 		{
	// 			AttributeName: aws.String("UserID"),
	// 			AttributeType: aws.String("S"),
	// 		},
	// 	},
	// 	[]*dynamodb.KeySchemaElement{
	// 		{
	// 			AttributeName: aws.String("UserID"),
	// 			KeyType:       aws.String("HASH"),
	// 		},
	// 	},
	// 	10, 10,
	// )
	// if e != nil {
	// 	fmt.Printf("Create Table Issue : %v", e)
	// } else {
	// 	fmt.Printf("Table Created")
	// }
	s3Client := commons.InitCustomS3Client(ctx, os.Getenv("BACKBLAZE_ID"), os.Getenv("BACKBLAZE_KEY"))

	mongoDB, err := commons.InitMongoDB(ctx)
	if err != nil {
		return fmt.Errorf("error getting connect mongo client: %v", err)
	}
	defer mongoDB.Disconnect(ctx)
	model := model.InitModel(mongoDB, s3Client)
	// model.Init()

	// // initialize handlers
	// handler := handlers.NewHandler(model, authClient)
	service := service.NewServer(model)

	// SSL Key
	certFilePath, _ := filepath.Abs(commons.ENVVariable("CRT_PATH"))
	keyFilePath, _ := filepath.Abs(commons.ENVVariable("KEY_PATH"))
	grpcClient.CertFilePath = certFilePath
	grpcClient.ServerAddr = commons.ENVVariable("AUTH_SERVER_ADDR")

	// run HTTP gateway
	go func() {
		fmt.Println("Run REST")
		// _ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort, certFilePath, keyFilePath)
		// _ = rest.RunCustomServer(ctx, cfg.GRPCPort, cfg.HTTPPort, certFilePath, keyFilePath)
		_ = rest.RunCustomServer(ctx, cfg.GRPCPort, cfg.HTTPPort, "", "")
	}()
	// fmt.Println("Run REST")
	// _ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort, certFilePath, keyFilePath)
	fmt.Println("Run gRPC")
	return grpc.RunServer(ctx, service, cfg.GRPCPort, certFilePath, keyFilePath)
}
