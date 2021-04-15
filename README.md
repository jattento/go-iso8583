[![codecov](https://codecov.io/gh/jattento/go-iso8583/branch/master/graph/badge.svg)](https://codecov.io/gh/jattento/go-iso8583)
[![Maintainability](https://api.codeclimate.com/v1/badges/94a2058a2b0823cf31be/maintainability)](https://codeclimate.com/github/jattento/go-iso8583/maintainability)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Last release](https://img.shields.io/github/v/release/jattento/go-iso8583?style=plastic)](https://github.com/jattento/go-iso8583/releases)


| Version   |      Build      |
|----------|:-------------:|
| Go 1.16 |[![Build Status](https://travis-ci.com/jattento/go-iso8583.svg?branch=master)](https://travis-ci.com/jattento/go-iso8583)|
| Go 1.15 |[![Build Status](https://travis-ci.com/jattento/go-iso8583.svg?branch=master)](https://travis-ci.com/jattento/go-iso8583)|
| Go 1.14 |[![Build Status](https://travis-ci.com/jattento/go-iso8583.svg?branch=master)](https://travis-ci.com/jattento/go-iso8583)|
| Go 1.13 |[![Build Status](https://travis-ci.com/jattento/go-iso8583.svg?branch=master)](https://travis-ci.com/jattento/go-iso8583)|
# go-iso8583

<img align="right" width="200px" src="https://raw.githubusercontent.com/jattento/go-iso8583-logo/main/iso%20gopher.png">

An easy to use, yet flexible marshaler for ISO-8583.

API [godoc](https://godoc.org/github.com/jattento/go-iso8583/pkg/iso8583) documentation.

This library provides:
- Marshal and Unmarshal functions with his respective interfaces
including MTI, VAR, LLVAR, LLLVAR and bitmaps fields ready for use
but with the possibility to easily add new field types.
- Inbuid Support for ASCII and EBCDIC but not limited to them


## Installation

To install go-iso8583 package, you need to install Go and set your Go workspace first.

1. First you need [Go](https://golang.org/) installed, then you can use the below Go command to install go-iso8583.
```sh
$ go get -u github.com/jattento/go-iso8583/iso8583
```

2. Import it in your code:
```go
import "github.com/jattento/go-iso8583/pkg/iso8583"
```

## Quick start

The API of this package is inspired in the go native json package
therefore it's pretty intuitive to use. Take a look at this!

```go
import "github.com/jattento/go-iso8583/pkg/iso8583"

type exampleMessage struct {
	MessageTypeIdentifier                     iso8583.MTI       `iso8583:"mti,length:4,encoding:ebcdic"`
	Bitmap                                    iso8583.BITMAP    `iso8583:"bitmap"`
	SecondaryBitmap                           iso8583.BITMAP    `iso8583:"1,omitempty"`
	PrimaryAccountNumber                      iso8583.LLVAR     `iso8583:"2,length:64,encoding:ebcdic,omitempty"`
	ProcessingCode                            iso8583.VAR       `iso8583:"3,length:6,encoding:ebcdic,omitempty"` 
	AmountTransaction                         iso8583.VAR       `iso8583:"4,length:12,encoding:ebcdic,omitempty"`
	AmountSettlement                          iso8583.VAR       `iso8583:"5,length:12,encoding:ebcdic,omitempty"`
	AmountCardholderBilling                   iso8583.VAR       `iso8583:"6,length:12,encoding:ebcdic,omitempty"`
	AcquiringInstitutionIDCode                iso8583.LLVAR     `iso8583:"32,length:11,encoding:ebcdic,omitempty"`
	ForwardingInstitutionIDCode               iso8583.LLVAR     `iso8583:"33,length:11,encoding:ebcdic,omitempty"`
	PrimaryAccountNumberExtended              iso8583.LLVAR     `iso8583:"34,length:28,encoding:ebcdic,omitempty"`
	Track2Data                                iso8583.LLVAR     `iso8583:"35,length:37,encoding:ebcdic,omitempty"`
	Track3Data                                iso8583.LLLVAR    `iso8583:"36,length:104,encoding:ebcdic,omitempty"`
	AdditionalResponseData                    iso8583.LLVAR     `iso8583:"44,length:25,encoding:ebcdic,omitempty"`
	Track1Data                                iso8583.LLVAR     `iso8583:"45,length:76,encoding:ebcdic,omitempty"`
	ExpandedAdditionalAmounts                 iso8583.LLLVAR    `iso8583:"46,length:999,encoding:ebcdic,omitempty"`
	AdditionalDataNationalUse                 iso8583.LLLVAR    `iso8583:"47,length:999,encoding:ebcdic,omitempty"`
	AdditionalDataPrivateUse                  iso8583.LLLVAR    `iso8583:"48,length:999,encoding:ebcdic,omitempty"`
	CurrencyCodeTransaction                   iso8583.VAR       `iso8583:"49,length:3,encoding:ebcdic,omitempty"`
	CurrencyCodeSettlement                    iso8583.VAR       `iso8583:"50,length:3,encoding:ebcdic,omitempty"`
	CurrencyCodeCardholderBilling             iso8583.VAR       `iso8583:"51,length:3,encoding:ebcdic,omitempty"`
	PersonalIDNumberData                      iso8583.BINARY    `iso8583:"52,length:8,omitempty"`
	SecurityRelatedControlInformation         iso8583.VAR       `iso8583:"53,length:16,encoding:ebcdic,omitempty"`
	AdditionalAmounts                         iso8583.LLLVAR    `iso8583:"54,length:120,encoding:ebcdic,omitempty"`
	IntegratedCircuitCardSystemRelatedData    iso8583.LLLBINARY `iso8583:"55,length:999,encoding:ebcdic,omitempty"`
    MessageSecurityCode                       iso8583.VAR       `iso8583:"96,length:8,encoding:ebcdic,omitempty"`
}
```

```go
func GenerateStaticReqBytes() ([]byte, error) {
	req := exampleMessage{
		MTI: "0100",
		// FirstBitmap is generated by library
		// SecondBitmap is generated by library
		PAN: "54000000000000111", // LL part is added by library!
		ProcessingCode: "1000",
		Amount: "000000000100",
		MessageNumber: "1",
	}
	
	byt, err := iso8583.Marshal(req)
	if err != nil {
		return nil, err
	}

	return byt, nil
}
```

```go
import "github.com/jattento/go-iso8583/pkg/iso8583"

func ReadResp(byt []byte) (exampleMessage,error){
	var resp exampleMessage

	_, err :=iso8583.Unmarshal(byt,&resp)
	if err != nil{
		return exampleMessage{}, err
	}

	return resp,nil
}
```

### [Changelog](changelog.md)