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
<head><meta charset="UTF-8"></head>
<body>
<p>`+c.name+` さん、こんにちは。</p>
<p>Chirp へご登録いただきありがとうございます。</p>
<p>以下のリンクをクリックしてアカウントを有効化してください：</p>
<p><a href="`+c.activationURL+`">アカウントを有効化する</a></p>
<p>このリンクは24時間有効です。</p>
</body>
</html>`)
	return err
}

func AccountActivationHTML(name, activationURL string) accountActivationComponent {
	return accountActivationComponent{name: name, activationURL: activationURL}
}
