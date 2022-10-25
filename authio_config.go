package authio

type AuthioConfig struct {
	Paths struct {
		// OnAuthNotRequired is the redirect path when calling authio.RedirectAuthNotRequired
		// Used if user reaches for example, a login page being authed (has session in cookies)
		OnAuthNotRequired string
	}
}

func (ac *AuthioConfig) Defaults() {
	ac.Paths.OnAuthNotRequired = "/"
}
