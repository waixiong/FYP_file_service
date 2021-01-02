package rest

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	// "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "getitqec.com/server/file/pkg/api/v1"
	// pb "../api/v1"

	"getitqec.com/server/file/pkg/logger"
	"getitqec.com/server/file/pkg/protocol/grpcClient"
	"getitqec.com/server/file/pkg/protocol/rest/middleware"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "Token":
		return "token", true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, grpcPort, httpPort, certFilePath string, keyFilePath string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// mux := runtime.NewServeMux()
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(CustomMatcher),
		runtime.WithErrorHandler(DefaultHTTPProtoErrorHandler),
		// runtime.WithProtoErrorHandler(DefaultHTTPProtoErrorHandler),
	)
	opts := []grpc.DialOption{}
	if certFilePath != "" && keyFilePath != "" {
		// creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		// creds, err := credentials.NewClientTLSFromFile(certFilePath, "CheeTest")
		// if err != nil {
		// 	log.Fatalf("Failed to generate credentials %v", err)
		// }

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

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	if err := pb.RegisterFileServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		// log.Fatalf("failed to start HTTP gateway: %v", err)
		logger.Log.Fatal("failed to start HTTP gateway: %v", zap.String("reason", err.Error()))
	}
	fmt.Println("REST : gRPC client up")

	// cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		Addr: ":" + httpPort,
		// Handler: mux,
		Handler: middleware.AddRequestID(middleware.AddLogger(logger.Log, mux)),
		// TLSConfig: tlsConfig,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	fmt.Println("starting HTTP/REST gateway...")
	logger.Log.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}

func RunCustomServer(ctx context.Context, grpcPort, httpPort, certFilePath string, keyFilePath string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{}
	if certFilePath != "" && keyFilePath != "" {
		// creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		// creds, err := credentials.NewClientTLSFromFile(certFilePath, "CheeTest")
		// if err != nil {
		// 	log.Fatalf("Failed to generate credentials %v", err)
		// }

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

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Println(certFilePath)
		log.Println(err)
		return err
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/api/file/image", imageHandle)

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		Addr: ":" + httpPort,
		// Handler: mux,
		Handler:   middleware.AddRequestID(middleware.AddLogger(logger.Log, handler)),
		TLSConfig: tlsConfig,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	fmt.Println("starting HTTP/REST gateway...")
	logger.Log.Info("starting HTTP/REST gateway...")
	// return srv.ListenAndServeTLS(certFilePath, keyFilePath)
	return srv.ListenAndServe()
}

func imageHandle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	fileClient, conn, err := grpcClient.GetFileClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v - Problem getting file service", http.StatusInternalServerError)))
		return
	}
	defer conn.Close()
	fmt.Printf("Protocol : %v\n", r.ProtoMajor)
	if r.Method == "GET" {
		names, ok := r.URL.Query()["name"]
		if !ok || len(names[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v - Please specific name", http.StatusBadRequest)))
			return
		}
		stream, err := fileClient.DownloadFile(r.Context(), &pb.DownloadImageRequest{Name: names[0]})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v - Problem downloading file", http.StatusInternalServerError)))
			return
		}
		// var bytes []byte
		// fmt.Printf("lllll : %v\n", len(bytes))
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("%v - err getting image", http.StatusInternalServerError)))
				log.Fatalf("err receive image : %v", err)
				return
			}
			// fmt.Printf("byte receive : %v\n", len(response.Image))
			// bytes = append(bytes, response.Image...)
			// if r.ProtoMajor == 2 {
			w.Write(response.Image)
			// }
		}
		// if r.ProtoMajor == 1 {
		// 	w.Write(bytes)
		// }
		return
	} else if r.Method == "POST" {
		fmt.Printf("Header : %v\n", r.Header)

		// Get form byte, modified content-d to Content-D
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v - fail read body", http.StatusInternalServerError)))
			// log.Fatalf("err ParseMultipartForm : %v", err)
			return
		}
		s := string(b)
		s = strings.Replace(s, "content-disposition", "Content-Disposition", -1)

		v := r.Header.Get("Content-Type")
		if v == "" {
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%v - err no content type : %v", http.StatusInternalServerError, err)))
				fmt.Printf("err ParseMultipartForm : %v\n", err)
				return
			}
		}
		d, params, err := mime.ParseMediaType(v)
		if err != nil || !(d == "multipart/form-data" /*|| d == "multipart/mixed"*/) {
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%v - err not multipart : %v", http.StatusInternalServerError, err)))
				fmt.Printf("err ParseMultipartForm : %v\n", err)
				return
			}
		}
		boundary, ok := params["boundary"]
		if !ok {
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%v - err no boundry : %v", http.StatusInternalServerError, err)))
				fmt.Printf("err ParseMultipartForm : %v\n", err)
				return
			}
		}
		multipartReader := multipart.NewReader(bytes.NewBuffer([]byte(s)), boundary)
		// multipartReader := multipart.NewReader(r.Body, params["boundary"])
		defer r.Body.Close()

		// bytes := []byte{}
		fileClient, conn, err := grpcClient.GetFileClient()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v - Problem getting file service", http.StatusInternalServerError)))
			return
		}
		defer conn.Close()
		stream, err := fileClient.UploadFile(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v - err open upload stteam", http.StatusInternalServerError)))
			log.Fatalf("err open uploaded stream : %v", err)
			return
		}
		// reading
		for {
			// fmt.Printf("Reading\n")
			part, err := multipartReader.NextPart()
			if err == io.EOF {
				// fmt.Printf("\tEOF\n")
				break
			}
			if err != nil {
				http.Error(w, "unexpected error when retrieving a part of the message", http.StatusInternalServerError)
				return
			}
			defer part.Close()
			fileBytes, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "failed to read content of the part", http.StatusInternalServerError)
				return
			}
			fmt.Printf("\t\tcase %v\n", part.Header.Get("Content-ID"))
			switch part.Header.Get("Content-ID") {
			case "metadata":
				fmt.Printf("metadata %v\n", string(fileBytes))

			case "media":
				fmt.Printf("filesize = %d\n", len(fileBytes))
			default:
				if part.FormName() == "name" {
					// fmt.Printf("name : %v\n", string(fileBytes))
					stream.Send(&pb.UploadImageRequestStream{
						Data: &pb.UploadImageRequestStream_Name{Name: string(fileBytes)},
					})

				} else if part.FormName() == "image" {
					// fmt.Printf("size = %d\n", len(fileBytes))
					// bytes = append(bytes, fileBytes...)
					err = stream.Send(&pb.UploadImageRequestStream{
						Data: &pb.UploadImageRequestStream_Image{Image: fileBytes},
					})
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(fmt.Sprintf("%v - internal network error", http.StatusInternalServerError)))
						log.Printf("internal network error : %v\n", err)
						return
					}
				}
			}
		}

		ack, err := stream.CloseAndRecv()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%v - err upload image", http.StatusInternalServerError)))
			log.Fatalf("err receive uploaded image : %v", err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(ack.Ack)))
		return
	}
	// Send a message back to the client
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(fmt.Sprintf("%v method is not implemented", r.Method)))
	return
}

func createFormData(fields map[string]string, boundary string) (string, []byte, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	if boundary != "" {
		mw.SetBoundary(boundary)
	}
	for k, v := range fields {
		mw.WriteField(k, v)
	}
	if err := mw.Close(); err != nil {
		return "", nil, err
	}
	return mw.FormDataContentType(), buf.Bytes(), nil
}
