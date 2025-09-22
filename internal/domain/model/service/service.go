package service

//go:generate go-enum -f=$GOFILE --marshal --names

// Type indicates service type that sends request to saf.
/*
ENUM(

	OfferService
	RideLifeCycleService
	PromotionService

)
*/
// nolint: recvcheck
type Type int
