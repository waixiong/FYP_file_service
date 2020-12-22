package grpcClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"

	pb "getitqec.com/server/file/pkg/api/v1"
	"getitqec.com/server/file/pkg/commons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	CertFilePath = ""
	ServerAddr   = ""
)

func GetFileClient() (pb.FileServiceClient, *grpc.ClientConn, error) {
	certFilePath, _ := filepath.Abs(commons.ENVVariable("CRT_PATH"))
	b, _ := ioutil.ReadFile(certFilePath)
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		log.Fatalf("fail to dial: %v", errors.New("credentials: failed to append certificates"))
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            cp,
	}
	creds := credentials.NewTLS(config)

	a := commons.ENVVariable("SERVICE_HOST") + ":" + commons.ENVVariable("GRPC_SERVICE_PORT")
	if len(a) == 0 {
		a = "0.0.0.0:8800"
	}
	conn, err := grpc.Dial(a, grpc.WithTransportCredentials(creds))
	return pb.NewFileServiceClient(conn), conn, err
}
