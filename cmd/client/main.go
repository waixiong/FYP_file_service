package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "getitqec.com/server/file/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Printf("Starting file service client...\n")
	client, conn, err := GetFileClient()
	if err != nil {
		fmt.Printf("get client fail : %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	ctx := context.TODO()
	file, err := os.Open("./cmd/client/choir.png")
	if err != nil {
		fmt.Printf("read file fail : %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	chunkSize := 32 * 1024
	buf := make([]byte, chunkSize)
	fileByte := []byte{}
	writing := true
	if err != nil {
		fmt.Printf("err send name : %v\n", err)
		os.Exit(1)
	}
	for writing {
		// put as many bytes as `chunkSize` into the
		// buf array.
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}

			if err != nil {
				fmt.Printf("errored while copying from file to buf : %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("\tsize %v\n", len(buf[:n]))
		fileByte = append(fileByte, buf[:n]...)
	}
	ack, err := client.UploadImage(ctx, &pb.UploadImageRequest{
		Bucket: "testBuc",
		Name:   "testwearefamily.png",
		Image:  fileByte,
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("id: %s\n", ack.Ack)
	}

	// // Upload Image
	// stream, err := client.UploadImage(ctx)
	// if err != nil {
	// 	fmt.Printf("init upload fail : %v\n", err)
	// 	os.Exit(1)
	// }

	// file, err := os.Open("./cmd/client/choir.png")
	// if err != nil {
	// 	fmt.Printf("read file fail : %v\n", err)
	// 	os.Exit(1)
	// }
	// defer file.Close()
	// chunkSize := 32 * 1024
	// buf := make([]byte, chunkSize)

	// writing := true
	// err = stream.Send(&pb.UploadFileRequest{
	// 	Data: &pb.UploadImageRequest_Name{
	// 		Name: "testing_img_from_golang",
	// 	},
	// })
	// if err != nil {
	// 	fmt.Printf("err send name : %v\n", err)
	// 	os.Exit(1)
	// }
	// for writing {
	// 	// put as many bytes as `chunkSize` into the
	// 	// buf array.
	// 	n, err := file.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			writing = false
	// 			err = nil
	// 			continue
	// 		}

	// 		if err != nil {
	// 			fmt.Printf("errored while copying from file to buf : %v\n", err)
	// 			os.Exit(1)
	// 		}
	// 	}

	// 	fmt.Printf("\tsize %v\n", len(buf[:n]))
	// 	err = stream.Send(&pb.UploadImageRequest{
	// 		Data: &pb.UploadImageRequest_Image{
	// 			Image: buf[:n],
	// 		},
	// 	})
	// 	if err != nil {
	// 		fmt.Printf("failed to send chunk via stream : %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }
	// ack, err := stream.CloseAndRecv()
	// if err != nil {
	// 	fmt.Printf("failed being acknowledged : %v\n", err)
	// } else {
	// 	fmt.Printf("Ack : %v\n", ack.Ack)
	// }

	// Download Image
	// stream, err := client.DownloadFile(ctx, &pb.DownloadImageRequest{Name: "testing_img_from_golang"})
	// if err != nil {
	// 	fmt.Printf("err init download : %v\n", err)
	// 	os.Exit(1)
	// }
	// // bytes := []byte{}
	// var bytes []byte
	// fmt.Printf("lllll : %v\n", len(bytes))
	// for {
	// 	response, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("err receive image : %v", err)
	// 	}
	// 	fmt.Printf("byte receive : %v\n", len(response.Image))
	// 	bytes = append(bytes, response.Image...)
	// }
	// file, err := os.Create("./cmd/client/receiveImg")
	// if err != nil {
	// 	log.Fatalf("err open file : %v", err)
	// }
	// defer file.Close()
	// file.Write(bytes)
}

func GetFileClient() (pb.FileServiceClient, *grpc.ClientConn, error) {
	certFilePath, _ := filepath.Abs(os.Getenv("CRT_PATH"))
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

	a := os.Getenv("SERVICE_HOST") + ":" + os.Getenv("GRPC_SERVICE_PORT")
	if len(a) == 0 {
		a = "0.0.0.0:8800"
	}
	conn, err := grpc.Dial(a, grpc.WithTransportCredentials(creds))
	return pb.NewFileServiceClient(conn), conn, err
}
