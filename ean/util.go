package ean

func c2i(c byte) int {
	return int(c) - '0'
}

func i2c(i int) byte {
	return byte('0' + i)
}
