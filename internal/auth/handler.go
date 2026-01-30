package auth

import (
	"encoding/json"
	"errors"
	middlewares "finance-app/middleware"
	"finance-app/utils"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service ServiceHandler
}

func NewHandler(r chi.Router, service ServiceHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", func(w http.ResponseWriter, r *http.Request) {
			var reqBody SignUpParams
			fmt.Println(reqBody.Email)
			fmt.Println(reqBody.Password)

			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			user, err := service.SignUp(ctx, reqBody)
			if err != nil {
				if errors.Is(err, utils.ErrorAccountConflict) {
					utils.Send(w, "not-ok", http.StatusConflict, utils.ErrorAccountConflict.Error(), nil)
					return
				}

				// SEND ANYWAY
				utils.Send(w, "not-ok", http.StatusInternalServerError, "internal server error", nil)
				return
			}

			utils.Send(w, "ok", http.StatusCreated, "create user successfull", user)
		})

		r.Post("/sign-in", func(w http.ResponseWriter, r *http.Request) {
			var reqBody SignInParams
			meta := SignInUserMeta{
				UserAgent: r.Header.Get("User-Agent"),
				IpAddress: r.RemoteAddr,
			}

			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			session, err := service.SignIn(ctx, reqBody, meta)
			if err != nil {
				if errors.Is(err, utils.ErrorEmailNotFound) {
					utils.Send(w, "not-ok", http.StatusUnauthorized, utils.ErrorEmailNotFound.Error(), nil)
					return
				}

				if errors.Is(err, utils.ErrorAccountWrongPassword) {
					utils.Send(w, "not-ok", http.StatusUnauthorized, utils.ErrorAccountWrongPassword.Error(), nil)
					return
				}

				if errors.Is(err, utils.ErrorAccountNotFound) {
					utils.Send(w, "not-ok", http.StatusNotFound, utils.ErrorAccountNotFound.Error(), nil)
					return
				}

				// SEND ANYWAY
				utils.Send(w, "not-ok", http.StatusInternalServerError, "internal server error", nil)
				return
			}

			utils.SetCookie(w, session.CreateSessionRow.Token)
			utils.SetCsrfCookie(w, session.CsrfToken)

			signInResponse := struct {
				ID        string `json:"id"`
				UserID    string `json:"user_id"`
				Email     string `json:"email"`
				CsrfToken string `json:"csrf_token"`
			}{
				ID:        session.CreateSessionRow.ID,
				UserID:    session.CreateSessionRow.UserID,
				Email:     reqBody.Email,
				CsrfToken: session.CsrfToken,
			}
			utils.Send(w, "ok", http.StatusOK, "sign in successfull", signInResponse)
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)
			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				token := ctx.Value("token").(string)

				user, err := service.Me(ctx, token)
				if err != nil {
					if errors.Is(err, utils.ErrorNoUserFound) {
						utils.Send(w, "not-ok", http.StatusUnauthorized, utils.ErrorNoUserFound.Error(), nil)
						return
					}

					if errors.Is(err, utils.ErrorSessionNotFound) {
						utils.Send(w, "not-ok", http.StatusUnauthorized, utils.ErrorSessionNotFound.Error(), nil)
						return
					}

					utils.Send(w, "not-ok", http.StatusInternalServerError, "internal server error", nil)
				}
				utils.Send(w, "ok", http.StatusOK, "success get me", user)
			})
		})

		r.Get("/csrf-token", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			csrfToken, err := service.Csrf(ctx)
			if err != nil {
				utils.Send(w, "not-ok", http.StatusInternalServerError, "internal server error", nil)
				return
			}

			utils.Send(w, "ok", http.StatusOK, "success get csrf token", csrfToken)
		})
	})
}
