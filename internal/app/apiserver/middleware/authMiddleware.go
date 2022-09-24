package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ressley/hacknu/internal/app/apiserver/helpers"
)

/*func IsAdmin(response http.ResponseWriter, request *http.Request) (bool, error) {
	err := Authentication(response, request)
	if err != nil {
		return false, err
	}
	uid := response.Header().Get("uid")
	_, err = role.GetAdminOne(&uid)
	if err != nil {
		return false, err
	}
	return true, nil
}*/

func Authentication(response http.ResponseWriter, request *http.Request) error {
	response.Header().Set("Content-Type", "application/json")
	authHeader, err := FromAuthHeader(request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"error" : "` + err.Error() + `"}`))
		return err
	}
	claims, _err := helpers.ValidateToken(authHeader)
	if _err != "" {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"error" : "` + _err + `"}`))
		return errors.New(_err)
	}
	response.Header().Set("login", claims.Login)
	response.Header().Set("firstName", claims.First_name)
	response.Header().Set("lastName", claims.Last_name)
	response.Header().Set("uid", claims.Uid)
	return nil
}

func FromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
