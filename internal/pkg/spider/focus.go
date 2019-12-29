package spider

type Focus struct {
	Value string
}

var (
	Href      = Focus{"href"}
	HrefText  = Focus{"href-text"}
	HrefImage = Focus{"href-image"}
	Image     = Focus{"image"}
	Script    = Focus{"script"}
	None      = Focus{"none"}
)

func NewFocus(name string) Focus {
	switch name {
	case Href.Value:
		return Href
	case HrefText.Value:
		return HrefText
	case HrefImage.Value:
		return HrefImage
	case Image.Value:
		return Image
	case Script.Value:
		return Script
	default:
		return None
	}
}
