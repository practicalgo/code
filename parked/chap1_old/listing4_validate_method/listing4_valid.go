package main

func (c httpConfig) Set(value string) error {
	allowedVerbs := []string{"HEAD", "GET", "POST"}
	allowed := false
	for _, c := range allowedVerbs {
		if value == c {
			allowed = true
		}
	}
	if !allowed {
		return errors.New(fmt.Sprintf("Verb not allowed: %s", value))
	}
	c.Verb = value
	return nil
}

func (c httpConfig) String() string {
	return c.Verb
}
