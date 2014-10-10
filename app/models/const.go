package models

type Role int

const (
	AnonymousRole Role = iota
	UserRole
	AdminRole
)

const (
	CSessionRole = "CSessionRole"
)

const (
	ArticlesInHomePanel     = 15
	ArticlesInUserHomePanel = 20
	ArticlesInSinglePage    = 30
	QiNiuSpace              = "http://sanwenjia.qiniudn.com/"
	QiNiuBucket             = "sanwenjia"
	QiNiuAccessKey          = "ZJ-KESjY89JYeZGT3dCvgQIngru-Qnkzt9PScvH1"
	QiNiuSecretKey          = "gRmeHtZ20UhTsUh0AfnlMHcUstwK7fR3vW-BSQQe"
	DefaultBoyAvatarUrl     = "http://sanwenjia.qiniudn.com/boy.PNG"
	DefaultGirlAvatarUrl    = "http://sanwenjia.qiniudn.com/girl.PNG"
	AnonymousAvatarUrl      = "http://sanwenjia.qiniudn.com/unknow.jpeg"
)
