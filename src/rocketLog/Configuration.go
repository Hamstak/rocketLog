package rocketLog

type configuration struct{
	webservice string
}

func readConfiguration() configuration{
	return configuration{"something"}
}