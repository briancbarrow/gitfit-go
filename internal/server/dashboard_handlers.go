package server

import (
	"context"
	"database/sql"
	"fmt"
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
	chartDataChan := make(chan dashboard_components.ChartData, 1)
	errChan := make(chan error, 3)

	go func() {
		chartData, err := app.GatherChartData(tenantQueries, r.Context())
		if err != nil {
			errChan <- err
			return
		}
		chartDataChan <- chartData
	}()

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
		fmt.Println("Workout set list: ", workoutSetList)
		workoutSetListChan <- workoutSetList
	}()

	tableAndFormOptions.ExerciseList = <-exerciseListChan
	tableAndFormOptions.WorkoutSetList = <-workoutSetListChan
	tableAndFormOptions.ChartData = <-chartDataChan

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
	time.Sleep(5 * time.Second)
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

type DeleteWorkoutSetPayload struct {
	Date string `json:"date"`
}

func (app *application) HandleDeleteSet(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Params: ", params)
	setId := params.ByName("id")
	setIdInt, err := strconv.ParseInt(setId, 10, 64)
	if err != nil {
		app.serverError(w, r, err)
	}
	fmt.Println("Deleting set with id: ", setId)
	// delete set
	err = tenantQueries.DeleteWorkoutSet(r.Context(), setIdInt)
	if err != nil {
		app.serverError(w, r, err)
	}
	fmt.Println("Deleted set with id: ", setId)
	var payload DeleteWorkoutSetPayload

	// err = json.NewDecoder(r.Body).Decode(&payload)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }

	q := r.URL.Query()
	payload.Date = q.Get("date")
	fmt.Println("Got payload: ", payload)
	tableAndFormOptions, err := app.GetTableAndFormOptions(r, dbID, payload.Date)
	if err != nil {
		app.serverError(w, r, err)
	}
	fmt.Println("Table and form options: ", tableAndFormOptions.WorkoutSetList)

	dashboard_components.TableAndForm(tableAndFormOptions).Render(r.Context(), w)
}

func (app *application) GatherChartData(tenantQueries *tenant_database.Queries, reqContext context.Context) (dashboard_components.ChartData, error) {

	workoutSetCounts, err := tenantQueries.GetWorkoutSetCounts(reqContext)
	if err != nil {
		return dashboard_components.ChartData{}, err
	}

	firstDayOfYear := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	weekDay := firstDayOfYear.Weekday()
	days := make([][]time.Time, 7)

	// Iterate over each day of the year
	currentDay := firstDayOfYear
	for currentDay.Year() == 2024 {
		// Calculate the day of the week and add the day to the corresponding list
		days[int(currentDay.Weekday())] = append(days[int(currentDay.Weekday())], currentDay)

		// Move to the next day
		currentDay = currentDay.AddDate(0, 0, 1)
	}

	chartInfo := dashboard_components.ChartData{
		FirstWeekday:     weekDay,
		DayData:          days,
		WorkoutSetCounts: workoutSetCounts,
		FirstDayFound:    false,
	}
	return chartInfo, nil
}
