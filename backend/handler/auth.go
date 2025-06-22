package handler

//var jwtKey = []byte("ladfioadfjia-jdhfoahsfo-dahlghasf") // use env var in production
//
//func GenerateJWT(email string) (string, error) {
//	claims := jwt.MapClaims{
//		"email": email,
//		"exp":   time.Now().Add(24 * time.Hour).Unix(),
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString(jwtKey)
//}
//
//func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		cookie, err := r.Cookie("session_token")
//		if err != nil {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		tokenStr := cookie.Value
//		claims := jwt.MapClaims{}
//
//		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
//			return jwtKey, nil
//		})
//		if err != nil || !token.Valid {
//			http.Error(w, "Invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		// Store user info in context if needed
//		r = r.WithContext(context.WithValue(r.Context(), "email", claims["email"]))
//		next.ServeHTTP(w, r)
//	}
//}
