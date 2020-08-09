package mti

import "strconv"


type MTI int

func New(origin Origin, function Function, class Class, version Version) MTI {
	return MTI(origin + function + class + version)
}

func (mti MTI) Origin() Origin {
	return atoi(itoa(int(mti) + 10000)[4:5])
}

func (mti MTI) Function() Function {
	return atoi(itoa(int(mti) + 10000)[3:5])
}

func (mti MTI) Class() Class {
	return atoi(itoa(int(mti) + 10000)[2:5])
}

func (mti MTI) Version() Version {
	return atoi(itoa(int(mti) + 10000)[1:5])
}

type Origin = int

const (
	OriginAcquirer       Origin = iota //
	OriginAcquirerRepeat               //
	OriginIssuer                       //
	OriginIssuerRepeat                 //
	OriginOther                        //
	OriginOtherRepeat                  //
	OriginReservedByISO6               //
	OriginReservedByISO7               //
	OriginReservedByISO8               //
	OriginReservedByISO9               //
)

type Function = int

const (
	FunctionRequest                     Function = iota * 10 // Request from acquirer to issuer to carry out an action; issuer may accept or reject
	FunctionRequestResponse                                  // Issuer response to a request
	FunctionAdvice                                           // Advice that an action has taken place; receiver can only accept, not reject
	FunctionAdviceResponse                                   // Response to an advice
	FunctionNotification                                     // Notification that an event has taken place; receiver can only accept, not reject
	FunctionNotificationAcknowledgement                      // Response to a notification
	FunctionInstruction                                      // ISO 8583:2003
	FunctionInstructionAcknowledgement                       // Instruction acknowledgement
	FunctionReservedByISO8                                   // Reserved for ISO.
	FunctionReservedByISO9                                   // Reserved for ISO.
)

type Class = int

const (
	ClassReservedByISO000              Class = iota * 100 //
	ClassAuthorizationMessage                             // Determine if funds are available, get an approval but do not post to account for reconciliation. Dual message system (DMS), awaits file exchange for posting to the account.
	ClassFinancialMessages                                // Determine if funds are available, get an approval and post directly to the account. Single message system (SMS), no file exchange after this.
	ClassFileActionsMessage                               // Used for hot-card, TMS and other exchanges
	ClassReversalAndChargebackMessages                    // Reversal (x4x0 or x4x1): Reverses the action of a previous authorization. Chargeback (x4x2 or x4x3): Charges back a previously cleared financial message.
	ClassReconciliationMessage                            // Transmits settlement information message.
	ClassAdministrativeMessage                            // Transmits administrative advice. Often used for failure messages (e.g., message reject or failure to apply).
	ClassFeeCollectionMessages                            //
	ClassNetworkManagementMessage                         // Used for secure key exchange, logon, echo test and other network functions.
	ClassReservedByISO900                                 //
)

type Version = int

const (
	Version8583To1987        Version = iota * 1000 //
	Version8583To1993                              //
	Version8583To2003                              //
	VersionReservedByISO3000                       //
	VersionReservedByISO4000                       //
	VersionReservedByISO5000                       //
	VersionReservedByISO6000                       //
	VersionReservedByISO7000                       //
	VersionNationalUse                             //
	VersionPrivateUse                              //
)

func itoa(n int) string { return strconv.Itoa(n) }
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
