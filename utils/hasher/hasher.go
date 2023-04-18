package hasher_utils


import "golang.org/x/crypto/bcrypt"

func HashPassword(pwd string) string {
	hpwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hpwd)
}

func CheckPassword(pwd string, hpwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hpwd), []byte(pwd))
}
