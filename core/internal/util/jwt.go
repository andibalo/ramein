package util

//func GenerateToken(user model.User) (tokenString string, err error) {
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"name":      user.Name,
//		"email":     user.Phone,
//		"role":      user.Role,
//		"timestamp": user.Timestampz,
//	})
//
//	tokenString, err = token.SignedString([]byte(viper.GetString("JWT_SECRET")))
//	if err != nil {
//		log.Println(err)
//
//		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to sign JWT token")
//	}
//
//	return tokenString, nil
//}
