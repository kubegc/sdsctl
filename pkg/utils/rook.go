package utils

func GetSecret() (string, error) {
	scmd := &Command{
		Cmd: "grep key /etc/ceph/keyring | awk '{print $3}'",
	}
	return scmd.Execute()
}
