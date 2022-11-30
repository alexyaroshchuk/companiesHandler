package middlewares

func AccessibleRoles() map[string][]string {
	return map[string][]string{
		"/companies.CompanyService/Create": {"admin"},
		"/companies.CompanyService/Patch":  {"admin"},
		"/companies.CompanyService/Put":    {"admin"},
		"/companies.CompanyService/Delete": {"admin"},
	}
}
