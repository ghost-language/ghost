package color

import "syscall"

// readTermios is the ioctl that reads a terminal's settings. Linux and the BSDs
// spell it differently, and the number is per-platform, so each names its own.
const readTermios = syscall.TCGETS
