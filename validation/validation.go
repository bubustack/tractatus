// Package validation exposes BubuStack transport contract validation helpers.
package validation

import (
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

// Validate checks a generated tractatus protobuf message against its
// Protovalidate annotations.
func Validate(msg proto.Message) error {
	return protovalidate.Validate(msg)
}
