package utils

var (
	AutoCodeVerify = Rules{"DbAlias": {NotEmpty()}}
	IdVerify       = Rules{"ID": []string{NotEmpty()}}
)
