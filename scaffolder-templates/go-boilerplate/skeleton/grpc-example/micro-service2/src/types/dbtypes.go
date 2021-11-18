package types

type Address struct {
	Country string `bson:"country"`
	State   string `bson:"state"`
	City    string `bson:"city"`
	Postal  int64  `bson:"postal"`
}

type Interests struct {
	Interest string `bson:"interest"`
	Priority int64  `bson:"priority"`
}

type Info struct {
	Email     string      `bson:"email"`
	Address   Address     `bson:"address"`
	Interests []Interests `bson:"interests"`
}
