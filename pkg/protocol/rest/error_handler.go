package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"code": 13, "message": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType(fallback))
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

type errorBody struct {
	Err  string `json:"error,omitempty"`
	Code int    `json:"code,omitempty"`
}

// (context.Context, *ServeMux, Marshaler, http.ResponseWriter, *http.Request, error)
func DefaultHTTPProtoErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	// return Internal when Marshal failed
	const fallback = `{"code": 13, "message": "failed to marshal error message"}`

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	w.Header().Del("Trailer")

	contentType := marshaler.ContentType(fallback)
	// Check marshaler on run time in order to keep backwards compatibility
	// An interface param needs to be added to the ContentType() function on
	// the Marshal interface to be able to remove this check
	// if typeMarshaler, ok := marshaler.(runtime.contentTypeMarshaler); ok {
	// pb := s.Proto()
	// contentType = typeMarshaler.ContentTypeFromMessage(pb)
	// }
	w.Header().Set("Content-Type", contentType)

	buf, merr := marshaler.Marshal(s.Proto())
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s.Proto(), merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	// md, ok := runtime.ServerMetadataFromContext(ctx)
	// if !ok {
	// 	grpclog.Infof("Failed to extract ServerMetadata from context")
	// }

	// runtime.handleForwardResponseServerMetadata(w, mux, md)
	// runtime.handleForwardResponseTrailerHeader(w, md)
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	// handleForwardResponseTrailer(w, md)
}
