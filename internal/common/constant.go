package common

const (
	developmentEnvironment = "development"
	localEnvironment       = "local"
	productionEnvironment  = "production"

	emailRegexp        = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	passwordRegexp     = `^[a-zA-Z0-9]{8,}([._]?[a-zA-Z0-9]+)*$`
	phoneRegexp        = `^[0-9]{10,12}$`
	timeDurationRegexp = `^([01]*\d|2[0-3])h?([0-5]\d)m?([0-5]\d)s$`
	uRLRegexp          = `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	usernameRegexp     = `^[a-zA-Z0-9]{8,}([._]?[a-zA-Z0-9]+)*$`
)
