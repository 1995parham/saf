package service

//go:generate go-enum -f=$GOFILE --marshal --names

/*
ENUM(

	OfferService
	RideLifeCycleService
	PromotionService

)
*/
// Type indicates service type that sends request to saf.
// nolint: recvcheck
type Type int
