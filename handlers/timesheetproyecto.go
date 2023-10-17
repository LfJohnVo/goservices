package handlers

import (
	"database/sql"
	"fmt"
	"goservices/databases"
	"goservices/models"
	"goservices/pkg"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/xuri/excelize/v2"
	_ "github.com/xuri/excelize/v2"
)

func GetProyectoReport(c *fiber.Ctx) error {
	var requestBody models.TimesheetProyectoRequestBody

	// Parse the request body into the RequestBody struct
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error" + err.Error())
	}

	// Access the ID from the parsed request body
	EmpID := requestBody.EmpID
	FechaInicio := requestBody.FechaInicio
	FechaFin := requestBody.FechaFin
	ProyectoID := requestBody.ProyectoID
	Solicitante := requestBody.Solicitante

	// Define the base SQL query
	sqlQuery := `
		SELECT
			timesheet_horas.id,
			empleados.name AS nombreEmpleado,
			timesheet_proyectos.proyecto AS proyecto,
			timesheet_tareas.tarea AS tarea,
			timesheet_horas.descripcion AS descripcion,
			timesheet.empleado_id,
			TO_CHAR(timesheet.fecha_dia, 'DD-MM-YYYY') AS fecha_dia,
			(
				SELECT supervisor.name FROM empleados AS supervisor WHERE supervisor.id = empleados.supervisor_id
			) AS supervisor_id,
			(
				SELECT SUM(
					COALESCE(horas_lunes::numeric, 0) +
					COALESCE(horas_martes::numeric, 0) +
					COALESCE(horas_miercoles::numeric, 0) +
					COALESCE(horas_jueves::numeric, 0) +
					COALESCE(horas_viernes::numeric, 0) +
					COALESCE(horas_sabado::numeric, 0) +
					COALESCE(horas_domingo::numeric, 0)
				)
				FROM timesheet_horas AS subquery
				WHERE subquery.id = timesheet_horas.id
			) AS totalHoras
		FROM
			timesheet_horas
		JOIN timesheet ON timesheet.id = timesheet_horas.timesheet_id
		JOIN empleados ON empleados.id = timesheet.empleado_id
		JOIN timesheet_tareas ON timesheet_tareas.id = timesheet_horas.tarea_id
		JOIN timesheet_proyectos ON timesheet_proyectos.id = timesheet_horas.proyecto_id`

	// Define a flag to track whether any conditions have been added
	conditionsAdded := false

	// Construct the date range condition if both FechaInicio and FechaFin are provided
	if FechaInicio != "" && FechaFin != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
			conditionsAdded = true
		}
		sqlQuery += " timesheet.fecha_dia BETWEEN '" + FechaInicio + "' AND '" + FechaFin + "'"
	} else if FechaInicio != "" {
		// Construct the condition for FechaInicio if it's provided
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
			conditionsAdded = true
		}
		sqlQuery += " timesheet.fecha_dia = '" + FechaInicio + "'"
	} else if FechaFin != "" {
		// Construct the condition for FechaFin if it's provided
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
			conditionsAdded = true
		}
		sqlQuery += " timesheet.fecha_dia = '" + FechaFin + "'"
	}

	// Add conditions for ProyectoID and EmpID if provided
	if ProyectoID != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
		}
		sqlQuery += " timesheet_proyectos.id = '" + ProyectoID + "'"
	}

	if EmpID != "" {
		if conditionsAdded {
			sqlQuery += " AND"
		} else {
			sqlQuery += " WHERE"
		}
		sqlQuery += " timesheet.empleado_id = '" + EmpID + "'"
	}

	conditionsAdded = true

	// Execute the query with the updated SQL query string
	QueryResult, err := databases.Database.Raw(sqlQuery).Rows()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("SqlQuery: " + sqlQuery + "\nerror" + err.Error())
	}
	// Ensure QueryResult is properly closed after use
	defer QueryResult.Close()

	// Create a slice to hold the query results
	var timesheetHoras []models.TimesheetProyectoResponseBody

	// Concatenate the date and time to the sheetName
	sheetName := "ProyectoReport_" + carbon.Now().ToDateString()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Set column headers
	headers := []string{"Empleado", "Proyecto", "Tarea", "Descripción", "Empleado ID", "Fecha Día", "Supervisor", "Total Horas"}
	for i, header := range headers {
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)
	}

	var i int = 2
	for QueryResult.Next() {

		var th models.TimesheetProyectoResponseBody

		// Scan into pointers to handle NULL values
		var id, EmpleadoName, Proyecto, Tarea, Descripcion, EmpleadoID, FechaDia, Supervisor, TotalHoras sql.NullString
		if err = QueryResult.Scan(&id, &EmpleadoName, &Proyecto, &Tarea, &Descripcion, &EmpleadoID, &FechaDia, &Supervisor, &TotalHoras); err != nil {
			log.Fatal(err)
			return c.Status(fiber.StatusInternalServerError).JSON("error" + err.Error())
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
		//en caso de querer retornar el slice como json
		timesheetHoras = append(timesheetHoras, th)

		// Add the values to the sheet
		//f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+0)), i), i-1)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+0)), i), th.EmpleadoName)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+1)), i), th.Proyecto)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+2)), i), th.Tarea)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+3)), i), th.Descripcion)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+4)), i), th.EmpleadoID)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+5)), i), th.FechaDia)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+6)), i), th.Supervisor)
		f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+7)), i), th.TotalHoras)

		i++
	}

	if err := f.SaveAs("storage/" + sheetName + ".xlsx"); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create the XLSX file",
		})
	}
	defer os.Remove("storage/" + sheetName + ".xlsx") // Delete the temporary file after sending

	// Set the response headers to indicate that it's an XLSX file
	c.Set("Content-Disposition", "attachment; filename="+sheetName+".xlsx")
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Specify the file path
	filePath := "pkg/emailtemplate.html" // Update with your file's path

	// Read the file
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		return err
	}

	// Convert the byte slice to a string
	contentString := string(fileContents)
	pkg.SendEmail([]string{Solicitante}, "Reporte registro colaboradores tareas", contentString, true, "storage/"+sheetName+".xlsx")

	// Send the XLSX file as the response
	//return c.Status(fiber.StatusOK).SendFile("storage/" + sheetName + ".xlsx")
	return c.Status(fiber.StatusOK).JSON("Reporte generado exitosamente")
}
