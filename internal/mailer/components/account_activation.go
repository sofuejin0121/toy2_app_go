package components

import (
	"context"
	"io"
)

type accountActivationComponent struct {
	name          string
	activationURL string
}

func (c accountActivationComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := io.WriteString(w, `<!DOCTYPE html>
<html>
<body>
<p>Hi `+c.name+`,</p>
<p>Click the link below to activate your account:</p>
<p><a href="`+c.activationURL+`">Activate Account</a></p>
</body>
</html>`)
	return err
}

func AccountActivationHTML(name, activationURL string) accountActivationComponent {
	return accountActivationComponent{name: name, activationURL: activationURL}
}
