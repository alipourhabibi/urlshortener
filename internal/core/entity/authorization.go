package entity

type TokenDetails struct {
	AccessToken      string
	RefreshToken     string
	AccessTokenUuid  string
	RefreshTokenUuid string
	ATExpires        int64
	RTExpires        int64
}
