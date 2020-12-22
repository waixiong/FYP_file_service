package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	pb "getitqec.com/server/file/pkg/api/v1"
	"getitqec.com/server/file/pkg/handlers"
	"getitqec.com/server/file/pkg/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//pb "./proto"
)

var httpClient = &http.Client{}

// logger is to mock a sophisticated logging system. To simplify the example, we just print out the content.
// func logger(format string, a ...interface{}) {
// 	fmt.Printf("LOG:\t"+format+"\n", a...)
// }

// var (
// 	//port = flag.Int("port", 50051, "the port to serve on")

// 	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
// 	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
// )

// Server class
type Server struct {
	model model.FileModelI

	pb.UnimplementedFileServiceServer
}

// Register request response
// waiting implementation

// func (s *Server) CreateItem(ctx context.Context, req *pb.Item) (*empty.Empty, error) {
// 	// md, ok := metadata.FromIncomingContext(ctx)
// 	// fmt.Printf("\tGet metadata...\n")
// 	// if !ok {
// 	// 	return nil, commons.ErrMissingMetadata
// 	// }
// 	// fmt.Printf("\t%v\n", md)

// 	handler := &handlers.CreateItemHandler{Model: s.model}
// 	return &empty.Empty{}, handler.CreateItem(ctx, req)
// }

func (s *Server) UploadFile(srv pb.FileService_UploadFileServer) (err error) {
	// while there are messages coming
	var bytes []byte
	filename := ""
	// fmt.Printf("Uploading...\n")
	for {
		req, e := srv.Recv()
		// fmt.Printf("\tReceive\n")
		err = e
		if err != nil {
			if err == io.EOF {
				// fmt.Printf("\tDone\n")
				goto END
			}
			fmt.Printf("\terror\n")
			err = errors.Wrapf(err, "failed unexpectadely while reading chunks from stream")
			return
		}
		switch u := req.Data.(type) {
		case *pb.UploadImageRequestStream_Image:
			bytes = append(bytes, u.Image...)
			// fmt.Printf("\t\tbyte\n")
		case *pb.UploadImageRequestStream_Name:
			filename = u.Name
			// fmt.Printf("\t\tname : %v\n", u.Name)
		}
	}

END:
	// once the transmission finished, send the
	// confirmation if nothing went wrong
	handler := &handlers.UploadFileHandler{Model: s.model}
	handler.UploadFile(srv.Context(), filename, bytes)
	err = srv.SendAndClose(&pb.Acknowledgement{
		Ack: "OK",
	})
	// ...

	return
	// return status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}

const chunkSize = 32 * 1024

func (s *Server) DownloadFile(req *pb.DownloadImageRequest, srv pb.FileService_DownloadFileServer) error {
	handler := &handlers.DownloadFileHandler{Model: s.model}
	data, err := handler.DownloadFile(srv.Context(), req.Name)
	if err != nil {
		return err
	}
	response := &pb.DownloadImageReponseStream{}
	for currentByte := 0; currentByte < len(data); currentByte += chunkSize {
		if currentByte+chunkSize > len(data) {
			response.Image = data[currentByte:len(data)]
		} else {
			response.Image = data[currentByte : currentByte+chunkSize]
		}
		if err := srv.Send(response); err != nil {
			return err
		}
	}
	return nil
	// return status.Errorf(codes.Unimplemented, "method DownloadImage not implemented")
}

func (s *Server) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.Acknowledgement, error) {
	handler := handlers.UploadImageHandler{Model: s.model}
	id, err := handler.UploadImage(ctx, req.Bucket, req.Name, req.Image)
	return &pb.Acknowledgement{Ack: id}, err
	// return nil, status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}
func (s *Server) DownloadImage(ctx context.Context, req *pb.DownloadImageRequest) (*pb.DownloadImageReponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadImage not implemented")
}

// NewServer return new auth server service
func NewServer(model model.FileModelI) *Server {
	server := &Server{}
	server.model = model
	return server
}
