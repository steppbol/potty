package services

import (
	"context"
	"fmt"
	"time"

	"github.com/steppbol/activity-manager/internal/cache"
	"github.com/steppbol/activity-manager/internal/middleware"
)

type AuthenticationService struct {
	cache   *cache.RedisCache
	context *context.Context
}

func NewAuthenticationService(c *cache.RedisCache) *AuthenticationService {
	ctx := context.Background()
	return &AuthenticationService{
		cache:   c,
		context: &ctx,
	}
}

func (as AuthenticationService) CreateAuthentication(userId uint, td *middleware.TokenDetail) error {
	at := time.Unix(td.AccessTokenExpireDate, 0)
	rt := time.Unix(td.RefreshTokenExpireDate, 0)
	now := time.Now()

	err := as.cache.Client.Set(*as.context, td.AccessID.String(), fmt.Sprintf("%d", userId), at.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = as.cache.Client.Set(*as.context, td.RefreshID.String(), fmt.Sprintf("%d", userId), rt.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (as AuthenticationService) DeleteAuthentication(id string) (int64, error) {
	deleted, err := as.cache.Client.Del(*as.context, id).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}

func (as AuthenticationService) DeleteTokens(ad *middleware.AccessDetail) (bool, error) {
	dA, err := as.DeleteAuthentication(ad.AccessID)
	if err != nil {
		return false, err
	}

	dR, err := as.DeleteAuthentication(fmt.Sprintf("%s_%d", ad.AccessID, ad.UserID))
	if err != nil {
		return false, err
	}

	if dA != 1 || dR != 1 {
		return false, err
	}
	return true, nil
}
