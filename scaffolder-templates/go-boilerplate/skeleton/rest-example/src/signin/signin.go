package signin

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/methods"

	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

// Signin is a method for saving user email and password in postgre database
func Signin(email, password string) (bool, string, dbtypes.User) {
	str, err := config.PostgreString()
	if err != nil {
		zap.S().Error(err)
		return false, err.Error(), dbtypes.User{}
	}
	// connecting to db
	db, err := gorm.Open("postgres", str)
	if err != nil {
		zap.S().Error(err)
		return false, err.Error(), dbtypes.User{}
	}
	// close db instance whenever whole work completed
	defer db.Close()

	//Checking whether registered or not
	var acs []dbtypes.User
	db.Where("email ILIKE ?", email).Find(&acs)
	zap.S().Info(len(acs))
	// when no account found
	if len(acs) == 0 {
		return false, "Invalid username", dbtypes.User{}
	}

	if methods.CheckHashForPassword(acs[0].Password, password) {
		// when password matched
		return true, "", acs[0]
	}
	// when password not matched
	return false, "Invalid password.", dbtypes.User{}

}
