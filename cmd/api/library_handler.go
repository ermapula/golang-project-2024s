package main

import (
	"errors"
	"net/http"

	"github.com/ermapula/golang-project/pkg/model"
	"github.com/ermapula/golang-project/pkg/validator"
)

func (app *application) showLibraryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	games, err := app.models.Games.GetAllOfUser(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"games": games}, nil)
}

func (app *application) addLibraryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	gameId, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	game, err := app.models.Games.Get(gameId)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	wallet, err := app.models.Users.GetWallet(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	v := validator.New()
	v.Check(wallet.Balance >= game.Price, "wallet", "insufficient funds")
	if !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	err = app.models.Games.AddToLibrary(user.Id, gameId)
	if err != nil {
		switch {
		case err.Error() == "game already in library":
			app.failedValidatorResponse(w, r, map[string]string{"library": "game already in library"})
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	_, err = app.models.Users.UpdateWallet(user.Id, -game.Price)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
}

func (app *application) removeLibraryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Games.DeleteFromLibrary(user.Id, id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getWalletHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	wallet, err := app.models.Users.GetWallet(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"wallet": wallet}, nil)
}

func (app *application) updateWalletHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		Amount    float64 `json:"amount"`
		Operation string  `json:"operation"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	wallet, err := app.models.Users.GetWallet(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	v := validator.New()
	if input.Operation != "-" && input.Operation != "+" {
		v.AddError("operation", "operation must be '+' or '-'")
	} else if input.Operation == "-" {
		v.Check(wallet.Balance >= input.Amount, "amount", "insufficient funds")
	}
	if !v.Valid() {
		app.failedValidatorResponse(w, r, v.Errors)
		return
	}

	if input.Operation == "-" {
		input.Amount = -input.Amount
	}
	wallet, err = app.models.Users.UpdateWallet(user.Id, input.Amount)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"wallet": wallet}, nil)
}
