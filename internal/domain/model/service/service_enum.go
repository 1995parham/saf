// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package service

import (
	"fmt"
	"strings"
)

const (
	// TypeOfferService is a Type of type OfferService.
	TypeOfferService Type = iota
	// TypeRideLifeCycleService is a Type of type RideLifeCycleService.
	TypeRideLifeCycleService
	// TypePromotionService is a Type of type PromotionService.
	TypePromotionService
)

var ErrInvalidType = fmt.Errorf("not a valid Type, try [%s]", strings.Join(_TypeNames, ", "))

const _TypeName = "OfferServiceRideLifeCycleServicePromotionService"

var _TypeNames = []string{
	_TypeName[0:12],
	_TypeName[12:32],
	_TypeName[32:48],
}

// TypeNames returns a list of possible string values of Type.
func TypeNames() []string {
	tmp := make([]string, len(_TypeNames))
	copy(tmp, _TypeNames)
	return tmp
}

var _TypeMap = map[Type]string{
	TypeOfferService:         _TypeName[0:12],
	TypeRideLifeCycleService: _TypeName[12:32],
	TypePromotionService:     _TypeName[32:48],
}

// String implements the Stringer interface.
func (x Type) String() string {
	if str, ok := _TypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Type(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Type) IsValid() bool {
	_, ok := _TypeMap[x]
	return ok
}

var _TypeValue = map[string]Type{
	_TypeName[0:12]:  TypeOfferService,
	_TypeName[12:32]: TypeRideLifeCycleService,
	_TypeName[32:48]: TypePromotionService,
}

// ParseType attempts to convert a string to a Type.
func ParseType(name string) (Type, error) {
	if x, ok := _TypeValue[name]; ok {
		return x, nil
	}
	return Type(0), fmt.Errorf("%s is %w", name, ErrInvalidType)
}

// MarshalText implements the text marshaller method.
func (x Type) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Type) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
