package server

import (
	"database/sql"
	"net/http"
	"strconv"

	dashboard_components "github.com/briancbarrow/gitfit-go/cmd/web/ui/components/dashboard"
	"github.com/briancbarrow/gitfit-go/internal/database/tenancy"
	tenant_database "github.com/briancbarrow/gitfit-go/internal/database/tenancy/db"
)

func (app *application) HandleNewSet(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "authenticatedUserID")

	user, err := app.queries.GetUser(r.Context(), userID.(string))
	if err != nil {
		app.serverError(w, r, err)
	}
	dbID := user.DatabaseID

	tableAndFormOptions := dashboard_components.TableAndFormOptions{}

	tableAndFormOptions.TemplateData = app.newTemplateData(r)
	db, err := tenancy.OpenTenantDB(dbID)
	if err != nil {
		app.serverError(w, r, err)
	}

	tenantQueries := tenant_database.New(db)
	// parse data
	r.ParseForm()
	exerciseId, err := strconv.Atoi(r.Form.Get("exercise"))
	if err != nil {
		app.serverError(w, r, err)
	}

	reps, err := strconv.Atoi(r.Form.Get("reps"))
	if err != nil {
		app.serverError(w, r, err)
	}

	note := r.Form.Get("notes")
	var noteNullString sql.NullString
	if note != "" {
		noteNullString = sql.NullString{String: note, Valid: true}
	} else {
		noteNullString = sql.NullString{Valid: false}
	}

	// create new set in db
	_, err = tenantQueries.CreateWorkoutSet(r.Context(), tenant_database.CreateWorkoutSetParams{
		Exercise: int64(exerciseId),
		Reps:     int64(reps),
		Date:     r.Form.Get("date"),
		Note:     noteNullString,
	})
	if err != nil {
		app.serverError(w, r, err)
	}

	app.dashboardGet(w, r)
}
