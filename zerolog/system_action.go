package zerolog

func SystemStart(config any) {
	ev := logger.Info().Interface("config", config)
	Write(ev, KindSystem, TypeSystemStart)
}

func SystemShutdown(cause string) {
	ev := logger.Info().Str("cause", cause)
	Write(ev, KindSystem, TypeSystemShutdown)
}
