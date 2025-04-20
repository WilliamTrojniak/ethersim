package ethergame

type GameObject interface {
	OnEvent(Event) bool
}
