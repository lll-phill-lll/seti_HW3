package application

type App struct {
	Rand *Rand
	Serv *Serv
}

func (a *App) Start(host string) {
	a.Serv.StartServe(host)
}

