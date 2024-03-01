package common

var (
	Environment = struct {
		Development, Local, Production string
	}{
		Development: developmentEnvironment,
		Local:       localEnvironment,
		Production:  productionEnvironment,
	}
	Regexp = struct {
		Email, Password, Phone, TimeDuration, URL, Username string
	}{
		Email:        emailRegexp,
		Password:     passwordRegexp,
		Phone:        phoneRegexp,
		TimeDuration: timeDurationRegexp,
		URL:          uRLRegexp,
		Username:     usernameRegexp,
	}
)
