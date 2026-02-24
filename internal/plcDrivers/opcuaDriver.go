package plcdrivers

// struct that implements the iterface plcComunication/PLCClient
type OPCUAClient struct {
	ip   string
	rack int
	slot int
}

/*
// NewS7Client creates the client but doesn't open the comunication
func NewS7Client(ip string, rack, slot int) *S7Client {

}
*/
