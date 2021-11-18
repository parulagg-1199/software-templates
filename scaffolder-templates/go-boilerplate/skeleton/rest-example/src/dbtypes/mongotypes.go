package dbtypes

// Address is a strcuture it stores the information related to user address
type Address struct {
	City    string `json:"city" bson:"City"  binding:"required"`
	State   string `json:"state" bson:"State"  binding:"required"`
	Country string `json:"country" bson:"Country"  binding:"required"`
	ZipCode string `json:"zipCode" bson:"ZipCode"`
}

// Interests is a strcuture it stores the information related to user Interests
type Interests struct {
	Type string `json:"type" bson:"Type" binding:"required"`
	Name string `json:"name" bson:"Name" binding:"required"`
}

// UserInfo is a strcuture it stores the information related to user
type UserInfo struct {
	Email     string      `json:"-" bson:"Email"`
	Address   Address     `json:"address" bson:"Address" binding:"required"`
	Interests []Interests `json:"interests" bson:"Interests" binding:"required"`
}
