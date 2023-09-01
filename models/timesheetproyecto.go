package models

type TimesheetProyectoRequestBody struct {
	EmpID       string `json:"emp_id,omitempty"`
	FechaInicio string `json:"fecha_inicio,omitempty"`
	FechaFin    string `json:"fecha_fin,omitempty"`
	ProyectoID  string `json:"proy_id,omitempty"`
}

type TimesheetProyectoResponseBody struct {
	ID           string `json:"id,omitempty"`
	EmpleadoName string `json:"empleado_name,omitempty"`
	Proyecto     string `json:"proyecto,omitempty"`
	Tarea        string `json:"tarea,omitempty"`
	Descripcion  string `json:"descripcion,omitempty"`
	EmpleadoID   string `json:"empleado_id,omitempty"`
	FechaDia     string `json:"fecha_dia,omitempty"`
	Supervisor   string `json:"supervisor,omitempty"`
	TotalHoras   string `json:"total_horas,omitempty"`
}
