package model

const RunningEnvironment = "RUNNING_ENVIRONMENT"

type Environment string

const (
	Local      Environment = "local"
	Production Environment = "production"
)
