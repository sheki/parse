parse
=====

A go client library for interacting with parse.com.


A few todos:
* A full usage example
* Structured querying
* Roles / ACLs
* Facebook/Twitter Auth linking/unlinking.
* Push sending support


# parse
    import "github.com/tmc/parse"

Package parse provides a programmatic interface to the Parse.com API


## Variables

``` go
var BaseURL = "https://api.parse.com/"
```
BaseURL is the common URL prefix for all API requests.


## type Client
``` go
type Client struct {
    // contains filtered or unexported fields
}
```
Client is the primary struct that this package provides. It represents the
connection to the Parse API


### func NewClient
``` go
func NewClient(parseAppID string, RESTAPIKey string) (*Client, error)
```
NewClient creates a new Client to interact with the Parse API.


### func (\*Client) CallCloudFunction
``` go
func (c *Client) CallCloudFunction(functionName string, arguments interface{}) ([]byte, error)
```
CallCloudFunction invokes the given cloud code function.
The arguments parameter is serialized as JSON and provided as parameters.


### func (\*Client) CallCloudJob
``` go
func (c *Client) CallCloudJob(jobName string, arguments interface{}) ([]byte, error)
```
CallCloudJob schedules the given cloud code job.
The arguments parameter is serialized as JSON and provided as parameters.


### func (\*Client) Create
``` go
func (c *Client) Create(object Object) (objectID string, err error)
```
Create creates a Parse object. On success the new object's ID is returned.
The provided object is not modified.


### func (\*Client) CreateInstallation
``` go
func (c *Client) CreateInstallation(installation Installation) (installationID string, err error)
```
CreateInstallation creates an installation from an Installation object.


### func (\*Client) CreateUser
``` go
func (c *Client) CreateUser(user User) (userID, sessionToken string, err error)
```
CreateUser creates a user from the specified object. On success the new user's
ID and session token are returned. The provided object is not modified.


### func (\*Client) Delete
``` go
func (c *Client) Delete(object Object) error
```
Delete removes the provided object from the Parse data store.


### func (\*Client) DeleteFile
``` go
func (c *Client) DeleteFile(fileName string) error
```
DeleteFile removes a Parse File.


### func (\*Client) DeleteInstallation
``` go
func (c *Client) DeleteInstallation(installation Installation) error
```
DeleteInstallation removes an installation by ID.


### func (\*Client) DeleteUser
``` go
func (c *Client) DeleteUser(user User) error
```
DeleteUser deletes the provided user.


### func (\*Client) Get
``` go
func (c *Client) Get(objectID string, object Object) error
```
Get populates the passed object by looking up based on objectID.


### func (\*Client) GetInstallation
``` go
func (c *Client) GetInstallation(installationID string, installation Installation) error
```
GetInstallation populates the provided installation based on the installationID.


### func (\*Client) GetUser
``` go
func (c *Client) GetUser(userID string, user User) error
```
GetUser looks up a user by ID. The provided user is populated on success.


### func (\*Client) LoginUser
``` go
func (c *Client) LoginUser(username, password string, user User) error
```
LoginUser attempts to log in a user given the provided name an password.
The provided object is populated with the user fields.


### func (\*Client) PasswordResetRequest
``` go
func (c *Client) PasswordResetRequest(email string) error
```
PasswordResetRequest sends a password reset email to the provided email address.


### func (\*Client) Query
``` go
func (c *Client) Query(options *QueryOptions, destination interface{}) error
```
Query performs a lookup of objects based on query options. destination must be
a slice of types satisfying the Object interface.


### func (\*Client) QueryInstallations
``` go
func (c *Client) QueryInstallations(options *QueryOptions, destination []Installation) error
```
QueryInstallations queries Installation objects based on the provided options.


### func (\*Client) SetMasterKey
``` go
func (c *Client) SetMasterKey(masterKey string)
```
SetMasterKey attaches a master key to subsequest API requests.NewClient
in lieu of the REST API Key. Setting to an empty string removes this behavior.


### func (\*Client) SetSessionToken
``` go
func (c *Client) SetSessionToken(sessionToken string)
```
SetSessionToken attaches a session token to subsequent requests, authenticating
them as the user associated with the token.


### func (\*Client) TraceOff
``` go
func (c *Client) TraceOff()
```
TraceOff turns on API response tracing


### func (\*Client) TraceOn
``` go
func (c *Client) TraceOn(logger *log.Logger)
```
TraceOn turns on API response tracing to the given logger.


### func (\*Client) Update
``` go
func (c *Client) Update(object Object) (updateTime time.Time, err error)
```
Update submits the JSON serialization of object and on success returns the
updaated time. The provided object is not modified.


### func (\*Client) UpdateInstallation
``` go
func (c *Client) UpdateInstallation(installation Installation) (updateTime time.Time, err error)
```
UpdateInstallation updates an installation. The provided installation is not modified.
On success the updated time is returned.


### func (\*Client) UpdateUser
``` go
func (c *Client) UpdateUser(user User) (updateTime time.Time, err error)
```
UpdateUser updates the provided user with any provided fields and on success
returns the updated at time.


### func (\*Client) UploadFile
``` go
func (c *Client) UploadFile(name string, contents io.Reader, contentType string) (*ParseFile, error)
```
UploadFile uploads a Parse File from the provided filename, contents and
content type.


## type Error
``` go
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"error"`
}
```
Error represents a Parse API error.


### func (Error) Error
``` go
func (e Error) Error() string
```


## type Installation
``` go
type Installation interface {
    Object
}
```
Installation is the minimal interface a type reprenting an installation must satisfy.



## type Object
``` go
type Object interface {
    ObjectID() string
}
```
Object is the minimal interface a Parse Object must satisfy.


## type ParseFile
``` go
type ParseFile struct {
    URL  string `json:"url"`
    Name string `json:"name"`
}
```
ParseFile is a Parse File that has been uploaded.



## type ParseInstallation
``` go
type ParseInstallation struct {
    ParseObject
    Channels    []string `json:"channels,omitempty"`
    DeviceToken string   `json:"deviceToken,omitempty"`
    DeviceType  string   `json:"deviceType,omitempty"`
}
```
ParseInstallation is a representation of an application installed on a device.


## type ParseObject
``` go
type ParseObject struct {
    ID        string `json:"objectId,omitempty"`
    CreatedAt string `json:"createdAt,omitempty"`
    UpdatedAt string `json:"updatedAt,omitempty"`
}
```
ParseObject is a type that satisifies the Object interface and is provided for
embedding.



### func (ParseObject) ObjectID
``` go
func (o ParseObject) ObjectID() string
```
ObjectID returns the ID of the object.



## type ParseUser
``` go
type ParseUser struct {
    ParseObject
    Username     string `json:"username,omitempty"`
    Password     string `json:"password,omitempty"`
    SessionToken string `json:"sessionToken,omitempty"`
}
```
ParseUser is a type that should be embedded in any custom User object.



## type QueryOptions
``` go
type QueryOptions struct {
    Where string
}
```
QueryOptions represents the parameters to a Parse query.



## type User
``` go
type User interface {
    Object
}
```

User is the minimal interface a type reprenting a user must satisfy.

## type ClassNamer
``` go
type ClassNamer interface {
    ParseClassName() string
}
```
ClassNamer is an interface that allows a type to provide its associated
Parse.com object class name.
