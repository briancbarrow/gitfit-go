package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	dashboard_components "github.com/briancbarrow/gitfit-go/cmd/web/ui/components/dashboard"
	"github.com/briancbarrow/gitfit-go/cmd/web/ui/pages"
	"github.com/briancbarrow/gitfit-go/internal/database/tenancy"
	tenant_database "github.com/briancbarrow/gitfit-go/internal/database/tenancy/db"
	"github.com/julienschmidt/httprouter"
)

func (app *application) dashboardGet(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "authenticatedUserID")
	user, err := app.queries.GetUser(r.Context(), userID.(string))
	if err != nil {
		app.serverError(w, r, err)
	}
	dbID := user.DatabaseID

	tableAndFormOptions, err := app.GetTableAndFormOptions(r, dbID, time.Now().Format("2006-01-02"))
	if err != nil {
		app.serverError(w, r, err)
	}
	pages.Dashboard(tableAndFormOptions).Render(r.Context(), w)
}

func (app *application) GetTableAndFormOptions(r *http.Request, dbID string, date string) (dashboard_components.TableAndFormOptions, error) {
	tableAndFormOptions := dashboard_components.TableAndFormOptions{}

	tableAndFormOptions.TemplateData = app.newTemplateData(r)
	tableAndFormOptions.Date = date
	db, err := tenancy.OpenTenantDB(dbID)
	if err != nil {
		return dashboard_components.TableAndFormOptions{}, err
	}

	tenantQueries := tenant_database.New(db)

	exerciseListChan := make(chan []tenant_database.Exercise, 1)
	workoutSetListChan := make(chan []tenant_database.ListWorkoutSetsRow, 1)
	errChan := make(chan error, 2)

	go func() {
		exerciseList, err := tenantQueries.ListExercises(r.Context())
		if err != nil {
			errChan <- err
			return
		}
		exerciseListChan <- exerciseList
	}()

	go func() {
		workoutSetList, err := tenantQueries.ListWorkoutSets(r.Context(), date)
		if err != nil {
			errChan <- err
			return
		}
		workoutSetListChan <- workoutSetList
	}()

	tableAndFormOptions.ExerciseList = <-exerciseListChan
	tableAndFormOptions.WorkoutSetList = <-workoutSetListChan

	select {
	case err := <-errChan:
		return dashboard_components.TableAndFormOptions{}, err
	default:
		return tableAndFormOptions, nil
	}
}

func (app *application) GetWorkoutSetTable(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "authenticatedUserID")
	user, err := app.queries.GetUser(r.Context(), userID.(string))
	if err != nil {
		app.serverError(w, r, err)
	}

	dbID := user.DatabaseID
	date := r.URL.Query().Get("date")

	tableAndFormOptions, err := app.GetTableAndFormOptions(r, dbID, date)
	if err != nil {
		app.serverError(w, r, err)
	}

	dashboard_components.TableAndForm(tableAndFormOptions).Render(r.Context(), w)
}

func (app *application) HandleNewSet(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "authenticatedUserID")

	user, err := app.queries.GetUser(r.Context(), userID.(string))
	if err != nil {
		app.serverError(w, r, err)
	}
	dbID := user.DatabaseID

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

	date := r.Form.Get("date")
	// create new set in db
	_, err = tenantQueries.CreateWorkoutSet(r.Context(), tenant_database.CreateWorkoutSetParams{
		Exercise: int64(exerciseId),
		Reps:     int64(reps),
		Date:     date,
		Note:     noteNullString,
	})
	if err != nil {
		app.serverError(w, r, err)
	}

	tableAndFormOptions, err := app.GetTableAndFormOptions(r, dbID, date)
	if err != nil {
		app.serverError(w, r, err)
	}

	dashboard_components.TableAndForm(tableAndFormOptions).Render(r.Context(), w)
}

func (app *application) HandleDeleteSet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userID := app.sessionManager.Get(r.Context(), "authenticatedUserID")

	user, err := app.queries.GetUser(r.Context(), userID.(string))
	if err != nil {
		app.serverError(w, r, err)
	}
	dbID := user.DatabaseID

	db, err := tenancy.OpenTenantDB(dbID)
	if err != nil {
		app.serverError(w, r, err)
	}

	tenantQueries := tenant_database.New(db)
	params := httprouter.ParamsFromContext(r.Context())

	setId := params.ByName("id")
	setIdInt, err := strconv.ParseInt(setId, 10, 64)
	if err != nil {
		app.serverError(w, r, err)
	}

	// delete set
	err = tenantQueries.DeleteWorkoutSet(r.Context(), setIdInt)
	if err != nil {
		app.serverError(w, r, err)
	}

	tableAndFormOptions, err := app.GetTableAndFormOptions(r, dbID, r.Form.Get("date"))
	if err != nil {
		app.serverError(w, r, err)
	}

	dashboard_components.TableAndForm(tableAndFormOptions).Render(r.Context(), w)
}
