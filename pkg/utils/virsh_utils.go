package utils

func GetPoolInfo(name string) (string, error) {
	comm := &Command{
		Cmd: "virsh pool-info",
		Params: map[string]string{
			"--pool":    name,
			"--details": "",
		},
	}
	return comm.Execute()
}
