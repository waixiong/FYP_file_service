package grpcClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"

	authv1 "getitqec.com/server/file/pkg/api/clients/auth/v1"
	"getitqec.com/server/file/pkg/commons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetAuthClient() (authv1.AuthenticationServiceClient, *grpc.ClientConn, error) {
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

	a := commons.ENVVariable("AUTH_SERVER_ADDR")
	if len(a) == 0 {
		a = "0.0.0.0:8080"
	}
	conn, err := grpc.Dial(a, grpc.WithTransportCredentials(creds))
	return authv1.NewAuthenticationServiceClient(conn), conn, err
}
