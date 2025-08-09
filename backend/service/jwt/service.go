package jwt

import (
	"fmt"
	"log/slog"
	"time"
	"wedding-app/config"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	secretKey     string
	tokenDuration time.Duration
	logger        *slog.Logger
}

func NewJWTService(config config.AuthConfig, logger *slog.Logger) service.JWTService {
	return &jwtService{secretKey: config.SecretKey, tokenDuration: config.Duration, logger: logger}
}

type claims struct {
	UserID string `json:"UserID"`
	// Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generate implements service.JWTService.
func (j *jwtService) Generate(user *model.User) (*model.AccessToken, error) {
	expiresAt := time.Now().Add(j.tokenDuration)

	claims := claims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		j.logger.Error("didnt get signed token: ", utils.ErrAttr(err))
		return nil, err
	}

	accessToken := &model.AccessToken{
		Token:     signedToken,
		ExpiresAt: expiresAt,
	}

	return accessToken, nil
}

// Verify implements service.JWTService.
func (j *jwtService) Verify(tokenString string) (*model.AccessToken, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("token not parsed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in claims: %q: %w", claims.UserID, err)
	}

	accessToken := &model.AccessToken{
		Token:     tokenString,
		ExpiresAt: claims.ExpiresAt.Time,
		UserID:    userID,
		// Role:      claims.Role,
	}

	return accessToken, nil
}
