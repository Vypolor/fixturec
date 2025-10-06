package utils

import (
	"os"
	"syscall"
)

const FilePerm0644 = os.FileMode(syscall.S_IRUSR | syscall.S_IWUSR | syscall.S_IRGRP | syscall.S_IROTH)
