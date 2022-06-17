package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	orderSearchSchema "github.com/hellofresh/schema-registry-go/service/customer/order/v1beta1"
	subscriptionSchema "github.com/hellofresh/schema-registry-go/service/customer/subscription/v1beta1"
	customerSearchSchema "github.com/hellofresh/schema-registry-go/service/customer/v1beta1"
	searchSchema "github.com/hellofresh/schema-registry-go/service/customer/v1beta1"
	schema "github.com/hellofresh/schema-registry-go/service/shipping/tracking/v1beta2"
	v1 "github.com/hellofresh/schema-registry-go/shared/v1"
)

var (
	address = "owl-search-service-grpc.ahoy-k8s.hellofresh.io:443"

	// following doesn't work yet!, have to port forward to test the server
	// however internally it works (cluster-cluster call)
	// address     = "owl-search-service-rpc.staging-k8s.hellofresh.io:80"
)

func main() {
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	println(address)
	// Set up a connection to the server.

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn.Connect()

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// searchAll(conn)
	// searchOrders(conn)
	searchCustomers(conn)
	//searchSubscriptions(conn)

	//track(conn)

}

func track(conn *grpc.ClientConn) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := schema.NewOdlBoxServiceClient(conn)

	r, err := c.FindTrackingDetails(ctx, &schema.FindTrackingDetailsRequest{
		OdlBoxId:   "isdf",
		RegionCode: "us",
	})

	if err != nil {
		log.Fatalf("could not track box %v", err)
	}

	log.Printf("response: %v", r.Items)

}

func searchAll(conn *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := searchSchema.NewSearchServiceClient(conn)

	r, err := c.SearchAll(ctx, &searchSchema.SearchAllRequest{
		BusinessDivision: &v1.BusinessDivision{
			Brand:      0,
			RegionCode: "us",
		},
		Query: "hello",
	})

	if err != nil {
		// this should trigger as not implemented
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r.GetCustomerIds())
}

func searchOrders(conn *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	email := "test"
	size := int32(10)

	oc := orderSearchSchema.NewCustomerOrderSearchServiceClient(conn)
	res, err := oc.SearchOrders(ctx, &orderSearchSchema.SearchOrdersRequest{
		BusinessDivision: &v1.BusinessDivision{
			RegionCode: "us",
		},
		Email:    &email,
		PageSize: &size,
	})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	// we should get some result
	log.Println("total orders:", res.GetTotal())
	log.Println("IDs count:", len(res.GetIds()))
}

func searchCustomers(conn *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	email := "dasfasd"
	size := int32(10)

	cs := customerSearchSchema.NewCustomerSearchServiceClient(conn)
	res, err := cs.SearchCustomers(ctx, &customerSearchSchema.SearchCustomersRequest{
		BusinessDivision: &v1.BusinessDivision{
			RegionCode: "us",
		},
		Email:    &email,
		PageSize: &size,
	})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	// we should get some result
	log.Println("total customers:", res.GetTotal())
	log.Println("IDs count:", len(res.GetIds()))
}

func searchSubscriptions(conn *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	email := "lupugabrielsp+55@gmail.com"
	size := int32(10)

	ss := subscriptionSchema.NewCustomerSubscriptionSearchServiceClient(conn)
	res, err := ss.SearchSubscriptions(ctx, &subscriptionSchema.SearchSubscriptionsRequest{
		BusinessDivision: &v1.BusinessDivision{
			RegionCode: "us",
		},
		Email:    &email,
		PageSize: &size,
	})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	// we should get some result
	log.Println("total subscriptions:", res.GetTotal())
	log.Println("IDs count:", len(res.GetIds()))
}
