package errors

type JWTError struct {
	Err string
}

func (e *JWTError) Error() string {
	return e.Err
}

type CryptError struct {
	Err string
}

func (e *CryptError) Error() string {
	return e.Err
}

type ConfigError struct {
	Err string
}

func (e *ConfigError) Error() string {
	return e.Err
}

type DbConnectionError struct {
	Err string
}

func (e *DbConnectionError) Error() string {
	return e.Err
}

type DbQueryError struct {
	Err string
}

func (e *DbQueryError) Error() string {
	return e.Err
}

type EmailError struct {
	Err string
}

func (e *EmailError) Error() string {
	return e.Err
}
