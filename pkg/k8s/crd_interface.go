package k8s

type crd interface {
	Exist(name string) (bool, error)
	Get(name string)
}
