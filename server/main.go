package main

import "chess/server/application"

func getHost() string {
	return "127.0.0.1:8080"
}
func main() {
	r := &application.Rand{}
	serv := &application.Serv{Rand:r}
	app := application.App{Rand:r, Serv: serv}

	host := getHost()

	app.Start(host)

}


