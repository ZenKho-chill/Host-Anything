// Copyright 2026 Host Anything Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/host-anything/hostanything/pkg/types"
)

// LoginRequest defines the expected payload for the login endpoint.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse defines the successful login payload.
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthHandler handles POST /api/v1/auth/login.
func AuthHandler(cfg *types.SystemConfig, logger *slog.Logger) http.HandlerFunc {
	// Parse timeout, fallback to 24h if invalid
	timeout, err := time.ParseDuration(cfg.Auth.SessionTimeout)
	if err != nil {
		timeout = 24 * time.Hour
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		// Verify credentials
		// Note: In a real system we'd use bcrypt.CompareHashAndPassword.
		// For this milestone, we support plain text match if not hashed, or simulated.
		if req.Username != cfg.Auth.AdminUsername || req.Password != cfg.Auth.AdminPassword {
			// Fail2ban compatible logging: must include client IP and failure keyword
			if cfg.Auth.Fail2BanEnabled {
				logger.Warn("authentication failed",
					"event", "auth_failure",
					"username", req.Username,
					"ip", r.RemoteAddr,
				)
			}
			writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		// Generate JWT
		claims := jwt.RegisteredClaims{
			Subject:   req.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
		if err != nil {
			logger.Error("failed to sign token", "error", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to generate token")
			return
		}

		logger.Info("authentication successful",
			"event", "auth_success",
			"username", req.Username,
			"ip", r.RemoteAddr,
		)

		writeJSON(w, http.StatusOK, LoginResponse{Token: signedToken})
	}
}
