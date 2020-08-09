package mti_test

import (
	"testing"

	"github.com/jattento/go-iso8583/pkg/mti"
	"github.com/stretchr/testify/assert"
)

func TestMTI_Origin(t *testing.T) {
	assert.Equal(t, mti.OriginAcquirer, mti.MTI(0).Origin())
	assert.Equal(t, mti.OriginReservedByISO9, mti.MTI(9).Origin())
}

func TestMTI_Function(t *testing.T) {
	assert.Equal(t, mti.FunctionRequest, mti.MTI(0).Function())
	assert.Equal(t, mti.FunctionReservedByISO9, mti.MTI(90).Function())

}

func TestMTI_Class(t *testing.T) {
	assert.Equal(t, mti.ClassReservedByISO000, mti.MTI(0).Class())
	assert.Equal(t, mti.ClassReservedByISO900, mti.MTI(900).Class())
}

func TestMTI_Version(t *testing.T) {
	assert.Equal(t, mti.Version8583To1987, mti.MTI(0).Version())
	assert.Equal(t, mti.VersionPrivateUse, mti.MTI(9000).Version())

}

func TestNew(t *testing.T) {
	assert.Equal(t, mti.MTI(1111), mti.New(1, 10, 100, 1000))
}
