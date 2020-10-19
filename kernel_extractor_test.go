package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnBytes(t *testing.T) {
	data := make([]byte, 12)
	data[3] = '\037'
	data[4] = '\213'
	data[5] = '\010'

	data[9] = '\037'
	data[10] = '\213'
	data[11] = '\010'

	ke := NewKernelExtractor(&data, true)

	subBytes := ke.ReturnBytes(3)
	assert.Equal(t, []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x1f, 0x8b, 0x8}, subBytes)
}

func TestNoPattern(t *testing.T) {
	data := make([]byte, 50)
	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchGZIPPattern()
	assert.False(t, found)
	assert.Equal(t, []int(nil), idx)
}

func TestListHeaders(t *testing.T) {
	data := make([]byte, 12)
	data[3] = '\037'
	data[4] = '\213'
	data[5] = '\010'

	data[9] = '\037'
	data[10] = '\213'
	data[11] = '\010'

	ke := NewKernelExtractor(&data, true)

	err := ke.ListAllHeadersFound()
	assert.Nil(t, err)
}

func TestKernelDetect(t *testing.T) {
	data := []byte("             Linux  \n \n kernel/params.c \n ")
	ke := NewKernelExtractor(&data, false)

	isKernel := ke.isKernelImage(data)
	assert.True(t, isKernel)

	data2 := []byte("             Lin  \n \n kernel/params.c \n ")
	isKernel = ke.isKernelImage(data2)
	assert.False(t, isKernel)

	data = []byte("")
	isKernel = ke.isKernelImage(data)
	assert.False(t, isKernel)

	// test ignore
	ke = NewKernelExtractor(&data, true)
	isKernel = ke.isKernelImage([]byte(""))
	assert.True(t, isKernel)
}

func TestGZIPIndex(t *testing.T) {
	data := make([]byte, 50)
	data[10] = '\037'
	data[11] = '\213'
	data[12] = '\010'

	data[30] = '\037'
	data[31] = '\213'
	data[32] = '\010'

	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchGZIPPattern()
	assert.True(t, found)
	assert.Equal(t, []int{10, 30}, idx)
}

func TestXZIndex(t *testing.T) {
	data := make([]byte, 50)
	data[10] = '\375'
	data[11] = '7'
	data[12] = 'z'
	data[13] = 'X'
	data[14] = 'Z'
	data[15] = '\000'

	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchXZPattern()
	assert.True(t, found)
	assert.Equal(t, []int{10}, idx)
}

func TestBZIPIndex(t *testing.T) {
	data := make([]byte, 50)
	data[10] = 'B'
	data[11] = 'Z'
	data[12] = 'h'
	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchBZIPPattern()
	assert.True(t, found)
	assert.Equal(t, []int{10}, idx)
}

func TestLZMAIndex(t *testing.T) {
	data := make([]byte, 50)
	data[10] = '\135'
	data[11] = '\000'
	data[12] = '\000'
	data[13] = '\000'

	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchLZMAPattern()
	assert.True(t, found)
	assert.Equal(t, []int{10}, idx)
}

func TestLZOPIndex(t *testing.T) {
	data := make([]byte, 50)
	data[10] = '\211'
	data[11] = '\114'
	data[12] = '\132'

	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchLZOPPattern()
	assert.True(t, found)
	assert.Equal(t, []int{10}, idx)
}

func TestLZ4Index(t *testing.T) {
	data := make([]byte, 50)
	data[10] = '\002'
	data[11] = '!'
	data[12] = 'L'
	data[13] = '\030'

	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchLZ4Pattern()
	assert.True(t, found)
	assert.Equal(t, []int{10}, idx)
}

func TestZSTDIndex(t *testing.T) {
	data := make([]byte, 50)
	data[10] = '('
	data[11] = '\265'
	data[12] = '/'
	data[13] = '\375'

	ke := NewKernelExtractor(&data, true)

	found, idx := ke.searchZSTDPattern()
	assert.True(t, found)
	assert.Equal(t, []int{10}, idx)
}
