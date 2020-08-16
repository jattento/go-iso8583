package mti

import "strconv"

// MTI is a string representation of a iso8583 MTI field.
// MTI has inbuilt methods that allows compare with others MTI considering definitions of the protocol.
// Some MTI methods are going to panic if some of the characters aren't numeric.
type MTI string

// New creates a new MTI considering the definition of each character.
func New(origin origin, function function, class class, version version) MTI {
	return MTI(itoa(int(origin))[:1] + itoa(int(function))[:1] + itoa(int(class))[:1] + itoa(int(version))[:1])
}

// Origin returns the origin element.
// Origin panics if some of the MTI characters isn't numeric.
func (mti MTI) Origin() origin {
	return origin(atoi(string(mti)[3:4]))
}

// Function returns the function element.
// Function panics if some of the MTI characters isn't numeric.
func (mti MTI) Function() function {
	return function(atoi(string(mti)[2:3] + "0"))
}

// Class returns the class element.
// Class panics if some of the MTI characters isn't numeric.
func (mti MTI) Class() class {
	return class(atoi(string(mti)[1:2] + "00"))
}

// Version returns the version element.
// Version panics if some of the MTI characters isn't numeric.
func (mti MTI) Version() version {
	return version(atoi(string(mti)[0:1] + "000"))
}

// String converts the mti to string type.
func (mti MTI) String() string {
	return string(mti)
}

// Equal performs the comparision.
func (mti MTI) Equal(v MTI) bool {
	return mti == v
}

// LowerThan performs the comparision.
// LowerThan panics if some of the MTI characters isn't numeric.
func (mti MTI) LowerThan(v MTI) bool {
	return atoi(string(mti)) < atoi(string(v))
}

// LowerOrEqualThan performs the comparision.
// LowerOrEqualThan panics if some of the MTI characters isn't numeric.
func (mti MTI) LowerOrEqualThan(v MTI) bool {
	return atoi(string(mti)) <= atoi(string(v))
}

// HigherThan performs the comparision.
// HigherThan panics if some of the MTI characters isn't numeric.
func (mti MTI) HigherThan(v MTI) bool {
	return atoi(string(mti)) > atoi(string(v))
}

// HigherOrEqualThan performs the comparision.
// HigherOrEqualThan panics if some of the MTI characters isn't numeric.
func (mti MTI) HigherOrEqualThan(v MTI) bool {
	return atoi(string(mti)) >= atoi(string(v))
}

type origin int

const (
	OriginAcquirer       origin = iota //
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

type function int

const (
	FunctionRequest                     function = iota * 10 // Request from acquirer to issuer to carry out an action; issuer may accept or reject
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

type class int

const (
	ClassReservedByISO000              class = iota * 100 //
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

type version int

const (
	Version8583To1987        version = iota * 1000 //
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
	n, err := strconv.Atoi(s)
	if err != nil {
		panic("mti: " + err.Error())
	}

	return n
}
