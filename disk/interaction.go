package disk

// Note that the signature is different because
// we might add custom features in nes disk to
// intercept the request
// For each request, we check if mappers has valid
// address translation, then we read/write and return
// status

func (n *NesDisk) CRead(addr uint16, data *uint8) bool {
	var mappedAddr uint32
	if n.mapperHandler.CpuMapRead(addr, &mappedAddr) {
		*data = n.PrgRomData[mappedAddr]
		return true
	}
	return false
}

func (n *NesDisk) CWrite(addr uint16, data uint8) bool {
	var mappedAddr uint32
	if n.mapperHandler.CpuMapWrite(addr, &mappedAddr) {
		n.PrgRomData[mappedAddr] = data
		return true
	}
	return false
}

func (n *NesDisk) PRead(addr uint16, data *uint8) bool {
	var mappedAddr uint32
	if n.mapperHandler.PpuMapRead(addr, &mappedAddr) {
		*data = n.ChrRomData[mappedAddr]
		return true
	}
	return false
}

func (n *NesDisk) PWrite(addr uint16, data uint8) bool {
	var mappedAddr uint32
	if n.mapperHandler.PpuMapWrite(addr, &mappedAddr) {
		n.ChrRomData[mappedAddr] = data
		return true
	}
	return false
}
