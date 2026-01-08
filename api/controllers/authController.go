package controllers

import (
	"encoding/json"
	"net/http"

	"server/common/appError"
	"server/common/dto"
	"server/common/services"
	"server/common/validator"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	user, err := services.CreateUser(dto.UserCreateDtoToModel(req), req.Password)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	token, err := services.GenerateJWTToken(user)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]any{
		"user":  dto.ToUserResponseDto(*user),
		"token": token,
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req dto.Login

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	user, token, err := services.Login(req.Email, req.Password)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]any{
		"user":  dto.ToUserResponseDto(*user),
		"token": token,
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestPasswordResetDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.RequestPasswordReset(req.Email)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Always return success to prevent email enumeration
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "If an account with that email exists, a password reset code has been sent.",
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.ResetPassword(req.Email, req.ResetCode, req.NewPassword)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been reset successfully.",
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req dto.GoogleAuthDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	user, token, err := services.AuthenticateWithGoogle(req.IDToken)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]any{
		"user":  dto.ToUserResponseDto(*user),
		"token": token,
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func AppleAuth(w http.ResponseWriter, r *http.Request) {
	var req dto.AppleAuthDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	user, token, err := services.AuthenticateWithApple(req.IDToken, req.Email, req.FirstName, req.LastName)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]any{
		"user":  dto.ToUserResponseDto(*user),
		"token": token,
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
