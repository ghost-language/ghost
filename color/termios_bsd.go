//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package color

import "syscall"

const readTermios = syscall.TIOCGETA
