package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrUnknown           = errors.New("An unknown error occurred")
	ErrUnauthorized      = errors.New("Unauthorized")
	ErrRequiresMasterKey = errors.New("Operation requires Master key")

	ErrAccountAlreadyLinked              = Error{Code: 208, Message: "An existing account already linked to another user."}
	ErrCacheMiss                         = Error{Code: 120, Message: "The results were not found in the cache."}
	ErrCommandUnavailable                = Error{Code: 108, Message: "Tried to access a feature only available internally."}
	ErrConnectionFailed                  = Error{Code: 100, Message: "The connection to the Parse servers failed."}
	ErrDuplicateValue                    = Error{Code: 137, Message: "A unique field was given a value that is already taken."}
	ErrExceededQuota                     = Error{Code: 140, Message: "Exceeded an application quota. Upgrade to resolve."}
	ErrFacebookAccountAlreadyLinked      = Error{Code: 208, Message: "An existing Facebook account already linked to another user."}
	ErrFacebookIDMissing                 = Error{Code: 250, Message: "Facebook id missing from request"}
	ErrFacebookInvalidSession            = Error{Code: 251, Message: "Invalid Facebook session"}
	ErrFileDeleteFailure                 = Error{Code: 153, Message: "Fail to delete file."}
	ErrIncorrectType                     = Error{Code: 111, Message: "Field set to incorrect type."}
	ErrInternalServer                    = Error{Code: 1, Message: "Internal server error. No information available."}
	ErrInvalidACL                        = Error{Code: 123, Message: "Invalid ACL. An ACL with an invalid format was saved. This should not happen if you use PFACL."}
	ErrInvalidChannelName                = Error{Code: 112, Message: "Invalid channel name. A channel name is either an empty string (the broadcast channel) or contains only a-zA-Z0-9_ characters and starts with a letter."}
	ErrInvalidClassName                  = Error{Code: 103, Message: "Missing or invalid classname. Classnames are case-sensitive. They must start with a letter, and a-zA-Z0-9_ are the only valid characters."}
	ErrInvalidDeviceToken                = Error{Code: 114, Message: "Invalid device token."}
	ErrInvalidEmailAddress               = Error{Code: 125, Message: "The email address was invalid."}
	ErrInvalidFileName                   = Error{Code: 122, Message: "Invalid file name. A file name contains only a-zA-Z0-9_. characters and is between 1 and 36 characters."}
	ErrInvalidImageData                  = Error{Code: 150, Message: "Fail to convert data to image."}
	ErrInvalidJSON                       = Error{Code: 107, Message: "Malformed json object. A json dictionary is expected."}
	ErrInvalidKeyName                    = Error{Code: 105, Message: "Invalid key name. Keys are case-sensitive. They must start with a letter, and a-zA-Z0-9_ are the only valid characters."}
	ErrInvalidLinkedSession              = Error{Code: 251, Message: "Invalid linked session"}
	ErrInvalidNestedKey                  = Error{Code: 121, Message: "Keys may not include '$' or '.'."}
	ErrInvalidPointer                    = Error{Code: 106, Message: "Malformed pointer. Pointers must be arrays of a classname and an object id."}
	ErrInvalidProductIDentifier          = Error{Code: 146, Message: "The product identifier is invalid"}
	ErrInvalidPurchaseReceipt            = Error{Code: 144, Message: "Product purchase receipt is invalid"}
	ErrInvalidQuery                      = Error{Code: 102, Message: "You tried to find values matching a datatype that doesn't support exact database matching, like an array or a dictionary."}
	ErrInvalidRoleName                   = Error{Code: 139, Message: "Role's name is invalid."}
	ErrInvalidServerResponse             = Error{Code: 148, Message: "The Apple server response is not valid"}
	ErrLinkedIDMissing                   = Error{Code: 250, Message: "Linked id missing from request"}
	ErrMissingObjectID                   = Error{Code: 104, Message: "Missing object id."}
	ErrObjectNotFound                    = Error{Code: 101, Message: "Object doesn't exist, or has an incorrect password."}
	ErrObjectTooLarge                    = Error{Code: 116, Message: "The object is too large."}
	ErrOperationForbidden                = Error{Code: 119, Message: "That operation isn't allowed for clients."}
	ErrPaymentDisabled                   = Error{Code: 145, Message: "Payment is disabled on this device"}
	ErrProductDownloadFileSystemFailure  = Error{Code: 149, Message: "Product fails to download due to file system error"}
	ErrProductNotFoundInAppStore         = Error{Code: 147, Message: "The product is not found in the App Store"}
	ErrPushMisconfigured                 = Error{Code: 115, Message: "Push is misconfigured. See details to find out how."}
	ErrReceiptMissing                    = Error{Code: 143, Message: "Product purchase receipt is missing"}
	ErrTimeout                           = Error{Code: 124, Message: "The request timed out on the server. Typically this indicates the request is too expensive."}
	ErrUnsavedFile                       = Error{Code: 151, Message: "Unsaved file."}
	ErrUserCannotBeAlteredWithoutSession = Error{Code: 206, Message: "The user cannot be altered by a client without the session."}
	ErrUserCanOnlyBeCreatedThroughSignUp = Error{Code: 207, Message: "Users can only be created through sign up"}
	ErrUserEmailMissing                  = Error{Code: 204, Message: "The email is missing, and must be specified"}
	ErrUserEmailTaken                    = Error{Code: 203, Message: "Email has already been taken"}
	ErrUserIDMismatch                    = Error{Code: 209, Message: "User ID mismatch"}
	ErrUsernameMissing                   = Error{Code: 200, Message: "Username is missing or empty"}
	ErrUsernameTaken                     = Error{Code: 202, Message: "Username has already been taken"}
	ErrUserPasswordMissing               = Error{Code: 201, Message: "Password is missing or empty"}
	ErrUserWithEmailNotFound             = Error{Code: 205, Message: "A user with the specified email was not found"}
	ErrScriptError                       = Error{Code: 141, Message: "Cloud Code script had an error."}
	ErrValidationError                   = Error{Code: 142, Message: "Cloud Code validation failed."}
)

// Error represents a Parse API error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e Error) Error() string {
	// TODO(tmc): improve formatting
	return fmt.Sprintf("Parse.com Error %v: %s", e.Code, e.Message)
}

func unmarshalError(r io.Reader) (error, bool) {
	err := &Error{}
	if r == nil {
		return ErrUnauthorized, false
	}
	if marshalErr := json.NewDecoder(r).Decode(&err); marshalErr != nil {
		return ErrUnknown, false
	}
	if byCode := errorByCode(err.Code); byCode != ErrUnknown {
		return byCode, true
	}
	return err, false
}

var errsByCode map[int]Error

func errorByCode(code int) error {
	err, ok := errsByCode[code]
	if !ok {
		return ErrUnknown
	}
	return err
}

func init() {
	errsByCode = make(map[int]Error)
	errs := []Error{ErrAccountAlreadyLinked, ErrCacheMiss, ErrCommandUnavailable, ErrConnectionFailed, ErrDuplicateValue, ErrExceededQuota, ErrFacebookAccountAlreadyLinked, ErrFacebookIDMissing, ErrFacebookInvalidSession, ErrFileDeleteFailure, ErrIncorrectType, ErrInternalServer, ErrInvalidACL, ErrInvalidChannelName, ErrInvalidClassName, ErrInvalidDeviceToken, ErrInvalidEmailAddress, ErrInvalidFileName, ErrInvalidImageData, ErrInvalidJSON, ErrInvalidKeyName, ErrInvalidLinkedSession, ErrInvalidNestedKey, ErrInvalidPointer, ErrInvalidProductIDentifier, ErrInvalidPurchaseReceipt, ErrInvalidQuery, ErrInvalidRoleName, ErrInvalidServerResponse, ErrLinkedIDMissing, ErrMissingObjectID, ErrObjectNotFound, ErrObjectTooLarge, ErrOperationForbidden, ErrPaymentDisabled, ErrProductDownloadFileSystemFailure, ErrProductNotFoundInAppStore, ErrPushMisconfigured, ErrReceiptMissing, ErrTimeout, ErrUnsavedFile, ErrUserCannotBeAlteredWithoutSession, ErrUserCanOnlyBeCreatedThroughSignUp, ErrUserEmailMissing, ErrUserEmailTaken, ErrUserIDMismatch, ErrUsernameMissing, ErrUsernameTaken, ErrUserPasswordMissing, ErrUserWithEmailNotFound, ErrScriptError, ErrValidationError}
	for _, err := range errs {
		errsByCode[err.Code] = err
	}
}
