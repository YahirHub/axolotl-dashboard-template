package controllers

import "github.com/goravel/framework/contracts/http"

type DashboardMetric struct {
	Label   string
	Value   string
	Change  string
	Detail  string
	SparkID string
	Ring    bool
	Tone    string
}

type DashboardProject struct {
	Key        string
	Name       string
	Area       string
	Budget     string
	Change     string
	ChangeDown bool
	Status     string
	SparkID    string
	Icon       string
	Tone       string
}

type DashboardTask struct {
	ID       string
	Title    string
	Project  string
	Owner    string
	Due      string
	Priority string
	Progress int
	Tone     string
}

type DashboardActivity struct {
	Title  string
	Detail string
	Time   string
	Icon   string
	Tone   string
}

type DashboardTeamMember struct {
	Name     string
	Role     string
	Initials string
	Status   string
	Load     int
}

type DashboardScheduleItem struct {
	Day   string
	Month string
	Title string
	Meta  string
	Tone  string
}

type DashboardIntegration struct {
	Name   string
	Detail string
	Status string
	Icon   string
	Tone   string
}

type DashboardPageData struct {
	Title        string
	Metrics      []DashboardMetric
	Projects     []DashboardProject
	Tasks        []DashboardTask
	Activities   []DashboardActivity
	Team         []DashboardTeamMember
	Schedule     []DashboardScheduleItem
	Integrations []DashboardIntegration
}

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (r *DashboardController) Index(ctx http.Context) http.Response {
	data := DashboardPageData{
		Title: "Axolotl — Dashboard Multipropósito",
		Metrics: []DashboardMetric{
			{Label: "Ingresos del mes", Value: "$1.28 M", Change: "+8.4%", Detail: "vs. mes anterior", SparkID: "sp1", Tone: "primary"},
			{Label: "Tareas completadas", Value: "327", Change: "+12.1%", Detail: "últimos 30 días", SparkID: "sp2", Tone: "success"},
			{Label: "Usuarios activos", Value: "1,248", Change: "+6.8%", Detail: "últimos 7 días", SparkID: "sp3", Tone: "secondary"},
			{Label: "Disponibilidad", Value: "99.97%", Change: "Operativo", Detail: "últimas 24 horas", Ring: true, Tone: "info"},
		},
		Projects: []DashboardProject{
			{Key: "WEB", Name: "Portal web", Area: "Producto", Budget: "$420,000", Change: "+12.4%", Status: "En ritmo", SparkID: "a1", Icon: "fa-globe", Tone: "primary"},
			{Key: "CRM", Name: "CRM comercial", Area: "Ventas", Budget: "$315,000", Change: "+7.2%", Status: "En ritmo", SparkID: "a2", Icon: "fa-address-book", Tone: "info"},
			{Key: "APP", Name: "App móvil", Area: "Producto", Budget: "$580,000", Change: "-3.1%", ChangeDown: true, Status: "Atención", SparkID: "a3", Icon: "fa-mobile-screen-button", Tone: "danger"},
			{Key: "DATA", Name: "Data Hub", Area: "Analítica", Budget: "$265,000", Change: "+4.9%", Status: "Estable", SparkID: "a4", Icon: "fa-database", Tone: "secondary"},
		},
		Tasks: []DashboardTask{
			{ID: "T-184", Title: "Validar flujo de onboarding", Project: "Portal web", Owner: "Mariana", Due: "Hoy · 16:30", Priority: "Alta", Progress: 72, Tone: "danger"},
			{ID: "T-207", Title: "Revisar segmentación de clientes", Project: "CRM comercial", Owner: "Diego", Due: "Mañana · 10:00", Priority: "Media", Progress: 48, Tone: "warning"},
			{ID: "T-231", Title: "Publicar build interno", Project: "App móvil", Owner: "Sofía", Due: "30 Sep", Priority: "Alta", Progress: 61, Tone: "primary"},
			{ID: "T-249", Title: "Documentar modelo de eventos", Project: "Data Hub", Owner: "Alex", Due: "2 Oct", Priority: "Baja", Progress: 34, Tone: "info"},
		},
		Activities: []DashboardActivity{
			{Title: "Despliegue completado", Detail: "Portal web · versión 2.4.0", Time: "Hace 18 min", Icon: "fa-rocket", Tone: "success"},
			{Title: "Nuevo comentario", Detail: "CRM comercial · revisión de pipeline", Time: "Hace 42 min", Icon: "fa-comment", Tone: "secondary"},
			{Title: "Alerta de capacidad", Detail: "App móvil · cola de procesamiento al 82%", Time: "Hace 1 h", Icon: "fa-triangle-exclamation", Tone: "warning"},
		},
		Team: []DashboardTeamMember{
			{Name: "Mariana López", Role: "Product Lead", Initials: "ML", Status: "Disponible", Load: 64},
			{Name: "Diego Ramos", Role: "Growth", Initials: "DR", Status: "En reunión", Load: 78},
			{Name: "Sofía Torres", Role: "Mobile", Initials: "ST", Status: "Disponible", Load: 53},
			{Name: "Alex Vega", Role: "Data", Initials: "AV", Status: "Concentrado", Load: 71},
		},
		Schedule: []DashboardScheduleItem{
			{Day: "29", Month: "SEP", Title: "Revisión de sprint", Meta: "09:30 · Producto", Tone: "primary"},
			{Day: "30", Month: "SEP", Title: "Demo App móvil", Meta: "13:00 · Stakeholders", Tone: "secondary"},
			{Day: "02", Month: "OCT", Title: "Cierre operativo", Meta: "17:30 · Operaciones", Tone: "warning"},
			{Day: "03", Month: "OCT", Title: "Planeación Q4", Meta: "11:00 · Dirección", Tone: "info"},
		},
		Integrations: []DashboardIntegration{
			{Name: "SQLite", Detail: "Persistencia local", Status: "Conectado", Icon: "fa-database", Tone: "success"},
			{Name: "Correo", Detail: "Notificaciones", Status: "Demo", Icon: "fa-envelope", Tone: "warning"},
			{Name: "Webhooks", Detail: "Automatizaciones", Status: "Listo", Icon: "fa-code-branch", Tone: "primary"},
			{Name: "Storage", Detail: "Archivos y respaldos", Status: "Activo", Icon: "fa-hard-drive", Tone: "info"},
		},
	}

	return ctx.Response().View().Make("dashboard.tmpl", data)
}
