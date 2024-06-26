package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ermapula/golang-project/pkg/model"
	"github.com/ermapula/golang-project/pkg/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &model.User{
		Name: input.Name,
		Email: input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidatorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.Id, "games:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.Id, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return 
	}

	data := struct{
		Token *string `json:"token"`
		User *model.User `json:"user"`
	}{
		Token: &token.Plaintext,
		User: user,
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": data}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func(app *application) activateUserHandler(w http.ResponseWriter,  r *http.Request) {
	var input struct{
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(model.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidatorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(model.ScopeActivation, user.Id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addPermission(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int64 `json:"user_id"`
		Permission string `json:"permission"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidatePermission(v, input.Permission); !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Permissions.AddForUser(input.UserID, input.Permission)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"permission": input.Permission}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}