package handlers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"goservices/databases"
	"goservices/models"
	"log"
)

func GetProyectoReport(c *fiber.Ctx) error {
	var requestBody models.TimesheetProyectoRequestBody

	// Parse the request body into the RequestBody struct
	if err := c.BodyParser(&requestBody); err != nil {
		return err
	}

	// Access the ID from the parsed request body
	EmpID := requestBody.EmpID
	FechaInicio := requestBody.FechaInicio
	FechaFin := requestBody.FechaFin
	ProyectoID := requestBody.ProyectoID

	// Define the base SQL query
	sqlQuery := "SELECT " +
		"timesheet_horas.id, " +
		"empleados.name as nombreEmpleado, " +
		"timesheet_proyectos.proyecto as proyecto, " +
		"timesheet_tareas.tarea as tarea, " +
		"timesheet_horas.descripcion as descripcion, " +
		"timesheet.empleado_id, " +
		"timesheet.fecha_dia, " +
		"(SELECT supervisor.name FROM empleados AS supervisor WHERE supervisor.id = empleados.supervisor_id) AS supervisor_id, " +
		"(SELECT SUM(COALESCE(horas_lunes::numeric, 0) + COALESCE(horas_martes::numeric, 0) + COALESCE(horas_miercoles::numeric,0) + COALESCE(horas_jueves::numeric,0) + COALESCE(horas_viernes::numeric,0) + COALESCE(horas_sabado::numeric,0) + COALESCE(horas_domingo::numeric,0)) FROM timesheet_horas AS subquery WHERE subquery.id = timesheet_horas.id) AS totalHoras " +
		"FROM timesheet_horas " +
		"JOIN timesheet ON timesheet.id = timesheet_horas.timesheet_id " +
		"JOIN empleados ON empleados.id = timesheet.empleado_id " +
		"JOIN timesheet_tareas ON timesheet_tareas.id = timesheet_horas.tarea_id " +
		"JOIN timesheet_proyectos ON timesheet_proyectos.id = timesheet_horas.proyecto_id"

	// Define a flag to track whether any conditions have been added
	conditionsAdded := false

	if FechaInicio != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
			conditionsAdded = true
		}
		sqlQuery += " timesheet.fecha_dia = '" + FechaInicio + "'"
	}

	if FechaFin != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
			conditionsAdded = true
		}
		sqlQuery += " timesheet.fecha_dia = '" + FechaFin + "'"
	}

	if ProyectoID != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
		}
		sqlQuery += " proyecto.id = '" + ProyectoID + "'"
	}

	if EmpID != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
		}
		sqlQuery += " timesheet.empleado_id = '" + EmpID + "'"
		conditionsAdded = true
	}

	// Execute the query with the updated SQL query string
	QueryResult, err := databases.Database.Raw(sqlQuery).Rows()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("SqlQuery: " + sqlQuery + "\nerror" + err.Error())
	}
	// Ensure QueryResult is properly closed after use
	defer QueryResult.Close()

	// Create a slice to hold the query results
	var timesheetHoras []models.TimesheetProyectoResponseBody

	for QueryResult.Next() {
		var th models.TimesheetProyectoResponseBody

		// Scan into pointers to handle NULL values
		var id, EmpleadoName, Proyecto, Tarea, Descripcion, EmpleadoID, FechaDia, Supervisor, TotalHoras sql.NullString
		if err = QueryResult.Scan(&id, &EmpleadoName, &Proyecto, &Tarea, &Descripcion, &EmpleadoID, &FechaDia, &Supervisor, &TotalHoras); err != nil {
			log.Fatal(err)
		}

		// Check if the values are not NULL and assign to the struct
		if id.Valid {
			th.ID = id.String
		}
		if EmpleadoName.Valid {
			th.EmpleadoName = EmpleadoName.String
		}
		if Proyecto.Valid {
			th.Proyecto = Proyecto.String
		}
		if Tarea.Valid {
			th.Tarea = Tarea.String
		}
		if Descripcion.Valid {
			th.Descripcion = Descripcion.String
		}
		if EmpleadoID.Valid {
			th.EmpleadoID = EmpleadoID.String
		}
		if FechaDia.Valid {
			th.FechaDia = FechaDia.String
		}
		if Supervisor.Valid {
			th.Supervisor = Supervisor.String
		}
		if TotalHoras.Valid {
			th.TotalHoras = TotalHoras.String
		}

		// Append the result to the slice
		timesheetHoras = append(timesheetHoras, th)
	}

	return c.Status(fiber.StatusOK).JSON(timesheetHoras)

}
