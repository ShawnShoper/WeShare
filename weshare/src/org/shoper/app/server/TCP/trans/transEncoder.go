package trans

type TransEncoder struct {
}

func (t *TransEncoder) decode() error       {}
func (t *TransEncoder) finishDecode() error {}
func (t *TransEncoder) dispose() error      {}
