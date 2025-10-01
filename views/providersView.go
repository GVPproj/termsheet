package views

func RenderProviders() string {
	s := "PROVIDERS VIEW\n\n"
	s += "• Provider 1\n"
	s += "• Provider 2\n"
	s += "• Provider 3\n\n"
	s += "Press ESC to return to menu"
	return s
}
