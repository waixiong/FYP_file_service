package grpcClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"

	storev1 "getitqec.com/server/file/pkg/api/clients/store/v1"
	"getitqec.com/server/file/pkg/commons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetStoreClient() (storev1.StoreServiceClient, *grpc.ClientConn, error) {
	b, _ := ioutil.ReadFile(CertFilePath)
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		log.Fatalf("fail to dial: %v", errors.New("credentials: failed to append certificates"))
	}
	config := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            cp,
	}
	creds := credentials.NewTLS(config)

	a := commons.ENVVariable("STORE_SERVER_ADDR")
	if len(a) == 0 {
		a = "0.0.0.0:9090"
	}
	conn, err := grpc.Dial(a, grpc.WithTransportCredentials(creds))
	return storev1.NewStoreServiceClient(conn), conn, err
}
