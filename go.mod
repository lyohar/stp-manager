module stpManager

go 1.15

replace stpCommon => ../stp-common

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/lib/pq v1.9.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	stpCommon v0.0.0-00010101000000-000000000000
)
