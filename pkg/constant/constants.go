package constant

import "net/http"

// grpc status & msg
const (
	STATUS_OK         = http.StatusOK
	STATUS_BADREQUEST = http.StatusBadRequest
	STATUS_ERR        = http.StatusInternalServerError
)

const (
	MESSAGE_OK = "ok"
)

// log
const (
	DefualtLogPath = "/var/log/sdsctl.log"
)

// k8s GVK
const (
	DefaultGroup   = "doslab.io"
	DefaultVersion = "v1"
	VMD_Kind       = "VirtualMachineDisk"
	VMDS_Kind      = "virtualmachinedisks"
	VMP_Kind       = "VirtualMachinePool"
	VMPS_Kind      = "virtualmachinepools"
)
