package mti_test

import (
	"testing"

	"github.com/jattento/go-iso8583/pkg/mti"
	"github.com/stretchr/testify/assert"
)

func TestMTI_Origin(t *testing.T) {
	assert.Equal(t, mti.OriginAcquirer, mti.MTI("0000").Origin())
	assert.Equal(t, mti.OriginReservedByISO9, mti.MTI("0009").Origin())
}

func TestMTI_Function(t *testing.T) {
	assert.Equal(t, mti.FunctionRequest, mti.MTI("0000").Function())
	assert.Equal(t, mti.FunctionRequest, mti.MTI("0001").Function())
	assert.Equal(t, mti.FunctionReservedByISO9, mti.MTI("0091").Function())
}

func TestMTI_Class(t *testing.T) {
	assert.Equal(t, mti.ClassReservedByISO000, mti.MTI("0000").Class())
	assert.Equal(t, mti.ClassReservedByISO000, mti.MTI("0011").Class())
	assert.Equal(t, mti.ClassReservedByISO000, mti.MTI("0010").Class())
	assert.Equal(t, mti.ClassReservedByISO900, mti.MTI("0912").Class())
}

func TestMTI_Version(t *testing.T) {
	assert.Equal(t, mti.Version8583To1987, mti.MTI("0000").Version())
	assert.Equal(t, mti.Version8583To1987, mti.MTI("0111").Version())
	assert.Equal(t, mti.VersionPrivateUse, mti.MTI("9123").Version())

}

func TestMTI_String(t *testing.T) {
	assert.Equal(t, "0204", mti.MTI("0204").String())
}

func TestNew(t *testing.T) {
	assert.Equal(t, mti.MTI("1111"), mti.New(1, 10, 100, 1000))
}

func TestMTI_Equal(t *testing.T) {
	assert.False(t,mti.MTI("1000").Equal("0999"))
	assert.True(t,mti.MTI("1000").Equal("1000"))
	assert.False(t,mti.MTI("1000").Equal("1001"))
}

func TestMTI_HigherOrEqualThan(t *testing.T) {
	assert.True(t,mti.MTI("1000").HigherOrEqualThan("0999"))
	assert.True(t,mti.MTI("1000").HigherOrEqualThan("1000"))
	assert.False(t,mti.MTI("1000").HigherOrEqualThan("1001"))
}

func TestMTI_LowerOrEqualThan(t *testing.T) {
	assert.False(t,mti.MTI("1000").LowerOrEqualThan("0999"))
	assert.True(t,mti.MTI("1000").LowerOrEqualThan("1000"))
	assert.True(t,mti.MTI("1000").LowerOrEqualThan("1001"))
}

func TestMTI_HigherThan(t *testing.T) {
	assert.False(t,mti.MTI("1000").LowerThan("0999"))
	assert.False(t,mti.MTI("1000").LowerThan("1000"))
	assert.True(t,mti.MTI("1000").LowerThan("1001"))
}

func TestMTI_LowerThan(t *testing.T) {
	assert.True(t,mti.MTI("1000").HigherThan("0999"))
	assert.False(t,mti.MTI("1000").HigherThan("1000"))
	assert.False(t,mti.MTI("1000").HigherThan("1001"))
}
