package moonms

func adaptNoReturn(fn func()) ServerEvent {
	return func() any {
		fn()
		return nil
	}
}
