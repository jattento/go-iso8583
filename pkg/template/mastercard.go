package template

import (
	"github.com/jattento/go-iso8583/pkg/iso8583"
)

// MasterCardISO87 is a template of the communication settings used for MasterCard connectivity on ISO8583 on 1987 version.
type MasterCardISO87 struct {
	MessageTypeIdentifier                     iso8583.MTI       `iso8583:"mti,length:4,encoding:ebcdic"`
	Bitmap                                    iso8583.BITMAP    `iso8583:"bitmap"`
	SecondaryBitmap                           iso8583.BITMAP    `iso8583:"1,omitempty"`
	PrimaryAccountNumber                      iso8583.LLVAR     `iso8583:"2,length:64,encoding:ebcdic,omitempty"`
	ProcessingCode                            iso8583.VAR       `iso8583:"3,length:6,encoding:ebcdic,omitempty"`
	AmountTransaction                         iso8583.VAR       `iso8583:"4,length:12,encoding:ebcdic,omitempty"`
	AmountSettlement                          iso8583.VAR       `iso8583:"5,length:12,encoding:ebcdic,omitempty"`
	AmountCardholderBilling                   iso8583.VAR       `iso8583:"6,length:12,encoding:ebcdic,omitempty"`
	TransmissionDateAndTime                   iso8583.VAR       `iso8583:"7,length:10,encoding:ebcdic,omitempty"`
	AmountCardholderBillingFee                iso8583.VAR       `iso8583:"8,length:8,encoding:ebcdic,omitempty"`
	ConversionRateSettlement                  iso8583.VAR       `iso8583:"9,length:8,encoding:ebcdic,omitempty"`
	ConversionRateCardholderBilling           iso8583.VAR       `iso8583:"10,length:8,encoding:ebcdic,omitempty"`
	SystemTraceAuditNumber                    iso8583.VAR       `iso8583:"11,length:6,encoding:ebcdic,omitempty"`
	TimeLocalTransaction                      iso8583.VAR       `iso8583:"12,length:6,encoding:ebcdic,omitempty"`
	DateLocalTransaction                      iso8583.VAR       `iso8583:"13,length:4,encoding:ebcdic,omitempty"`
	DateExpiration                            iso8583.VAR       `iso8583:"14,length:4,encoding:ebcdic,omitempty"`
	DateSettlement                            iso8583.VAR       `iso8583:"15,length:4,encoding:ebcdic,omitempty"`
	DateConversion                            iso8583.VAR       `iso8583:"16,length:4,encoding:ebcdic,omitempty"`
	DateCapture                               iso8583.VAR       `iso8583:"17,length:4,encoding:ebcdic,omitempty"`
	MerchantType                              iso8583.VAR       `iso8583:"18,length:4,encoding:ebcdic,omitempty"`
	AcquiringInstitutionCountryCode           iso8583.VAR       `iso8583:"19,length:3,encoding:ebcdic,omitempty"`
	PrimaryAccountNumberCountryCode           iso8583.VAR       `iso8583:"20,length:3,encoding:ebcdic,omitempty"`
	ForwardingInstitutionCountryCode          iso8583.VAR       `iso8583:"21,length:3,encoding:ebcdic,omitempty"`
	PointOfServiceEntryMode                   iso8583.VAR       `iso8583:"22,length:3,encoding:ebcdic,omitempty"`
	CardSequenceNumber                        iso8583.VAR       `iso8583:"23,length:3,encoding:ebcdic,omitempty"`
	NetworkInternationalID                    iso8583.VAR       `iso8583:"24,length:3,encoding:ebcdic,omitempty"`
	PointOfServiceConditionCode               iso8583.VAR       `iso8583:"25,length:2,encoding:ebcdic,omitempty"`
	PointOfServicePersonalIDNumberCaptureCode iso8583.VAR       `iso8583:"26,length:2,encoding:ebcdic,omitempty"`
	AuthorizationIDResponseLength             iso8583.VAR       `iso8583:"27,length:1,encoding:ebcdic,omitempty"`
	AmountTransactionFee                      iso8583.VAR       `iso8583:"28,length:9,encoding:ebcdic,omitempty"`
	AmountSettlementFee                       iso8583.VAR       `iso8583:"29,length:9,encoding:ebcdic,omitempty"`
	AmountTransactionProcessingFee            iso8583.VAR       `iso8583:"30,length:9,encoding:ebcdic,omitempty"`
	AmountSettlementProcessingFee             iso8583.VAR       `iso8583:"31,length:9,encoding:ebcdic,omitempty"`
	AcquiringInstitutionIDCode                iso8583.LLVAR     `iso8583:"32,length:11,encoding:ebcdic,omitempty"`
	ForwardingInstitutionIDCode               iso8583.LLVAR     `iso8583:"33,length:11,encoding:ebcdic,omitempty"`
	PrimaryAccountNumberExtended              iso8583.LLVAR     `iso8583:"34,length:28,encoding:ebcdic,omitempty"`
	Track2Data                                iso8583.LLVAR     `iso8583:"35,length:37,encoding:ebcdic,omitempty"`
	Track3Data                                iso8583.LLLVAR    `iso8583:"36,length:104,encoding:ebcdic,omitempty"`
	RetrievalReferenceNumber                  iso8583.VAR       `iso8583:"37,length:12,encoding:ebcdic,omitempty"`
	AuthorizationIDResponse                   iso8583.VAR       `iso8583:"38,length:6,encoding:ebcdic,omitempty"`
	ResponseCode                              iso8583.VAR       `iso8583:"39,length:2,encoding:ebcdic,omitempty"`
	ServiceRestrictionCode                    iso8583.VAR       `iso8583:"40,length:3,encoding:ebcdic,omitempty"`
	CardAcceptorTerminalID                    iso8583.VAR       `iso8583:"41,length:8,encoding:ebcdic,omitempty"`
	CardAcceptorIDCode                        iso8583.VAR       `iso8583:"42,length:15,encoding:ebcdic,omitempty"`
	CardAcceptorNameLocation                  iso8583.VAR       `iso8583:"43,length:40,encoding:ebcdic,omitempty"`
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
	PaymentAccountData                        iso8583.LLLVAR    `iso8583:"56,length:999,encoding:ebcdic,omitempty"`
	ReservedForNationalUse57                  iso8583.LLLVAR    `iso8583:"57,length:999,encoding:ebcdic,omitempty"`
	ReservedForNationalUse58                  iso8583.LLLVAR    `iso8583:"58,length:999,encoding:ebcdic,omitempty"`
	ReservedForNationalUse59                  iso8583.LLLVAR    `iso8583:"59,length:999,encoding:ebcdic,omitempty"`
	AdviceReasonCode                          iso8583.LLLVAR    `iso8583:"60,length:999,encoding:ebcdic,omitempty"`
	PointOfServiceData                        iso8583.LLLVAR    `iso8583:"61,length:999,encoding:ebcdic,omitempty"`
	IntermediateNetworkFacilityData           iso8583.LLLVAR    `iso8583:"62,length:999,encoding:ebcdic,omitempty"`
	NetworkData                               iso8583.LLLVAR    `iso8583:"63,length:999,encoding:ebcdic,omitempty"`
	NetworkManagementInformationCode          iso8583.VAR       `iso8583:"70,length:3,encoding:ebcdic,omitempty"`
	OriginalDataElements                      iso8583.VAR       `iso8583:"90,length:42,encoding:ebcdic,omitempty"`
	ServiceIndicator                          iso8583.VAR       `iso8583:"94,length:7,encoding:ebcdic,omitempty"`
	ReplacementAmounts                        iso8583.VAR       `iso8583:"95,length:42,encoding:ebcdic,omitempty"`
	MessageSecurityCode                       iso8583.VAR       `iso8583:"96,length:8,encoding:ebcdic,omitempty"`
	AccountID1                                iso8583.LLVAR     `iso8583:"102,length:28,encoding:ebcdic,omitempty"`
	AccountID2                                iso8583.LLVAR     `iso8583:"103,length:28,encoding:ebcdic,omitempty"`
	DigitalPaymentData                        iso8583.LLLVAR    `iso8583:"104,length:999,encoding:ebcdic,omitempty"`
	MoneySendReferenceData                    iso8583.LLLVAR    `iso8583:"108,length:999,encoding:ebcdic,omitempty"`
	AdditionalData                            iso8583.LLLVAR    `iso8583:"112,length:999,encoding:ebcdic,omitempty"`
	RecordData                                iso8583.LLLVAR    `iso8583:"120,length:999,encoding:ebcdic,omitempty"`
	AuthorizingAgentIDCode                    iso8583.LLLVAR    `iso8583:"121,length:999,encoding:ebcdic,omitempty"`
	ReceiptFreeText                           iso8583.LLLVAR    `iso8583:"123,length:999,encoding:ebcdic,omitempty"`
	MemberDefinedData                         iso8583.LLLVAR    `iso8583:"124,length:999,encoding:ebcdic,omitempty"`
	PrivateData126                            iso8583.LLLVAR    `iso8583:"126,length:999,encoding:ebcdic,omitempty"`
	PrivateData127                            iso8583.LLLVAR    `iso8583:"127,length:999,encoding:ebcdic,omitempty"`
}
