package controllers_test

import (
	"net/http"

	"calculator/backend/api/controllers"
	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
	"calculator/backend/api/routes"
	"calculator/backend/api/services"
	"calculator/backend/api/views"
)

func newHandler() http.Handler {
	return routes.NewRouter(controllers.NewCalculatorController(services.NewCalculatorService(repositories.NewMemoryHistory)))
}

func initialCalculatorState() views.CalculatorState {
	return views.StateFromModel(models.InitialState())
}

func validateState(state views.CalculatorState) error {
	return services.ValidateState(state.ToModel())
}

func textPointer(value string) *string {
	return &value
}
