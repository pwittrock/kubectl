package framework

import (
	"k8s.io/kubectl/pkg/framework/internal/inject"
)

func Factory() *inject.Factory {
	return inject.FactorySingleton
}
