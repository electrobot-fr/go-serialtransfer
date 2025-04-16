package main

func cobsEncode(arr []byte, overheadByte int) {
	length := len(arr)
	refByte := overheadByte

	if refByte != -1 {
		for i := length - 1; i >= 0; i-- {
			if arr[i] == StartByte {
				arr[i] = byte(refByte - i)
				refByte = i
			}
		}
	}
}

func findLast(arr []byte) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == StartByte {
			return i
		}
	}
	return -1
}
