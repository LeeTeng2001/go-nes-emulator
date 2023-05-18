package bus

// Device represents a device connected to a bus
type Device interface {
	Write(addr uint16, data uint8)
	Read(addr uint16) (data uint8)
}
