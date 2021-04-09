package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"

	courseProto "github.com/fahimsGit/basic-microservice/proto/course"
	pb "github.com/fahimsGit/basic-microservice/proto/student"
	_ "github.com/fahimsGit/basic-microservice/statik"
	"github.com/fahimsGit/basic-microservice/student/handler"
	"github.com/fahimsGit/basic-microservice/student/repository"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/soheilhy/cmux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ms struct {
	addr string
	port int
	srv  *grpc.Server
	mux  *http.ServeMux
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	log.Printf("preflight request for %s", r.URL.Path)
}

func (s *ms) Run() {
	// Create the main listener.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux object.
	tcpm := cmux.New(l)

	// Declare the match for different services required.
	httpl := tcpm.Match(cmux.HTTP1Fast())
	grpcl := tcpm.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	http2 := tcpm.Match(cmux.HTTP2())

	// Link the endpoint to the handler function.
	http.Handle("/", allowCORS(s.mux))

	// Initialize the servers by passing in the custom listeners (sub-listeners).
	go s.ServeGRPC(grpcl)
	go s.ServeHTTP(httpl)
	go s.ServeHTTP(http2)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a C, handle it
			l.Close()
			os.Exit(1)
		}
	}()

	log.Println("GRPC server started.")
	log.Println("HTTP server started.")
	log.Printf("Server listening on port %d \n", s.port)

	// Start cmux serving.
	if err := tcpm.Serve(); !strings.Contains(err.Error(),
		"use of closed network connection") {
		log.Fatal(err)
	}
}

func (s *ms) ServeGRPC(l net.Listener) {
	if err := s.srv.Serve(l); err != nil {
		log.Fatalf("could not start GRPC sever: %v", err)
	}
}

func (s *ms) ServeHTTP(l net.Listener) {
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("could not start HTTP server: %v", err)
	}
}

func main() {
	// Init services

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Client not started %v", err)
	}
	defer conn.Close()
	var client courseProto.CourseServiceClient
	client = courseProto.NewCourseServiceClient(conn)

	ms := ms{
		addr: "localhost:10000",
		port: 10000,
	}

	dbClient, err := getDBClient()
	if err != nil {
		log.Fatalf("cold not get database client: %v", err)
	}
	repo := repository.NewMongoRepository(dbClient)
	svcImpl := handler.NewService(repo, client)

	ms.srv = grpc.NewServer()

	// Register GRPC
	pb.RegisterStudentServiceServer(ms.srv, svcImpl)
	reflection.Register(ms.srv)

	// Register GRPC gateway
	ms.mux = http.NewServeMux()
	gmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	ms.mux.Handle("/", gmux)

	err = pb.RegisterStudentServiceHandlerFromEndpoint(context.Background(), gmux, ms.addr, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("could not register grpc gateway: %v", err)
	}

	// Enable swagger
	err = enableSwagger(ms.mux, "/docs/", "./proto/student/student.swagger.json")
	if err != nil {
		log.Fatalf("could not enable swagger: %v", err)
	}

	ms.Run()
}

func enableSwagger(mux *http.ServeMux, prefix string, jsonFilePath string) error {
	if mux == nil {
		return errors.New("no http mux found: server was not created properly")
	}

	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	saticHandler := http.FileServer(statikFS)
	mux.Handle(prefix, http.StripPrefix(prefix, saticHandler))

	mux.HandleFunc(fmt.Sprintf("%vswagger.json", prefix), func(w http.ResponseWriter, req *http.Request) {
		source, err := os.Open(jsonFilePath)
		if err != nil {
			return
		}
		defer source.Close()
		io.Copy(w, source)
	})
	return nil
}

func getDBClient() (*mongo.Client, error) {
	fmt.Println("==========================")
	fmt.Println("==========================")
	fmt.Println("Connecting db client.")
	fmt.Println("==========================")
	fmt.Println("==========================")

	cs := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cs))
	if err != nil {
		return nil, err
	}
	// defer client.Disconnect(ctx)

	return client, nil
}
