package service

import (
	"fmt"

	"github.com/arstrel/rest-banking/auth/domain"
	"github.com/arstrel/rest-banking/auth/dto"
	"github.com/arstrel/rest-banking/errs"
	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
	Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {
			// continue with refresh token functionality
			var appErr *errs.AppError
			if appErr = s.repo.RefreshTokenExists(request.RefreshToken); appErr != nil {
				return nil, appErr
			}
			// generate a access token from refresh token
			var accessToken string
			if accessToken, appErr = domain.NewAccessTokenFromRefreshToken(request.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errs.NewAuthenticationError("invalid token")
	}
	return nil, errs.NewAuthorizationError("cannot generate a new access token until the current one expires")
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.Login

	if login, appErr = s.repo.FindBy(req.Username, req.Password); appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); appErr != nil {
		return nil, appErr
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	// convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthorizationError(err.Error())
	}
	// Checking the validity of the token, this varifies the expiry
	// time and the signature of the token
	if jwtToken.Valid {
		// type case the token claims to jwt.MapClaims
		claims := jwtToken.Claims.(*domain.AccessTokenClaims)
		// if Role is user then check if the account_id and customer_id
		// coming in the URL belongs to the same token
		if claims.IsUserRole() {
			if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
				return errs.NewAuthorizationError("request not verified with the token claims")
			}
		}
		// verify if the role is authorized to use the route
		isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
		if !isAuthorized {
			return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
		}
		return nil
	}
	return errs.NewAuthorizationError("Invalid token")
}