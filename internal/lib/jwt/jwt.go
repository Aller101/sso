package jwt

import (
	"time"

	"ginhub.com/Aller101/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uuid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix() //.Unix - конвертируем в таймштамп - передалим клиенту как метку времени
	claims["app_id"] = app.ID

	// секрет хранить в модели не оч хорошо - имеет риск быть залогированной,
	// а в логи нельзя секреты помещать - лучше секрет передавать отдельно
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
	// добавить тесты
}
