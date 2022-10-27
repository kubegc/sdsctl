package constant

import "net/http"

// grpc status & msg
const SocketPath = "/var/lib/libvirt/ks.sock"

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
	DefaultGroup     = "doslab.io"
	DefaultVersion   = "v1"
	DefaultNamespace = "default"

	VMD_Kind  = "VirtualMachineDisk"
	VMDS_Kind = "virtualmachinedisks"
	VMP_Kind  = "VirtualMachinePool"
	VMPS_Kind = "virtualmachinepools"

	CRD_Pool_Key   = "pool"
	CRD_Volume_Key = "volume"

	CRD_Pool_Active   = "active"
	CRD_Pool_Inactive = "inactive"

	CRD_Ready_Msg    = "The resource is ready."
	CRD_Ready_Reason = "Ready"
)
