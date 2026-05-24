package clients

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	listingpb "github.com/Kost0/internship-exchange/proto/listing"
)

type ListingClient struct {
	client listingpb.ListingServiceClient
}

func NewListingClient(addr string) (*ListingClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ListingClient{client: listingpb.NewListingServiceClient(conn)}, nil
}

func (c *ListingClient) SyncCompany(ctx context.Context, userID, name, logoURL, industry, city string) {
	_, err := c.client.SyncCompany(ctx, &listingpb.SyncCompanyRequest{
		UserId:   userID,
		Name:     name,
		LogoUrl:  logoURL,
		Industry: industry,
		City:     city,
	})
	if err != nil {
		log.Printf("SyncCompany error: %v", err)
	}
}
