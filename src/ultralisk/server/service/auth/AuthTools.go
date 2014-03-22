package auth

import (
	"ultralisk/util"
)

func (this *ServiceAuth) showInfo(a ...interface{}) {
	util.ShowInfo(this.titleOnShow, a)
}

func (this *ServiceAuth) showDebug(a ...interface{}) {
	util.ShowDebug(this.titleOnShow, a)
}

func (this *ServiceAuth) showWarnning(a ...interface{}) {
	util.ShowWarnning(this.titleOnShow, a)
}

func (this *ServiceAuth) showError(a ...interface{}) {
	util.ShowError(this.titleOnShow, a)
}
