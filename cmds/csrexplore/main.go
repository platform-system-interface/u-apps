package main

import (
	"fmt"
	"github.com/u-root/u-root/pkg/memio"
)

/**

0x000 URW ustatus User status register.
0x004 URW uie User interrupt-enable register.
0x005 URW utvec User trap handler base address.
User Trap Handling
0x040 URW uscratch Scratch register for user trap handlers.
0x041 URW uepc User exception program counter.
0x042 URW ucause User trap cause.
0x043 URW ubadaddr User bad address.
0x044 URW uip User interrupt pending.
User Floating-Point CSRs
0x001 URW fflags Floating-Point Accrued Exceptions.
0x002 URW frm Floating-Point Dynamic Rounding Mode.
0x003 URW fcsr Floating-Point Control and Status Register (frm + fflags).
User Counter/Timers
0xC00 URO cycle Cycle counter for RDCYCLE instruction.
0xC01 URO time Timer for RDTIME instruction.
0xC02 URO instret Instructions-retired counter for RDINSTRET instruction.
0xC03 URO hpmcounter3 Performance-monitoring counter.
0xC04 URO hpmcounter4 Performance-monitoring counter.
.
.
.
0xC1F URO hpmcounter31 Performance-monitoring counter.
0xC80 URO cycleh Upper 32 bits of cycle, RV32I only.
0xC81 URO timeh Upper 32 bits of time, RV32I only.
0xC82 URO instreth Upper 32 bits of instret, RV32I only.
0xC83 URO hpmcounter3h Upper 32 bits of hpmcounter3, RV32I only.
0xC84 URO hpmcounter4h Upper 32 bits of hpmcounter4, RV32I only.
.
.
.
0xC9F URO hpmcounter31h Upper 32 bits of hpmcounter31, RV32I only.

*/

type csrpriv string

const (
	URW csrpriv = "URW"
	URO         = "URO"
)

type CSR struct {
	Number      int64
	Privilege   csrpriv
	Name        string
	Description string
}

var CSRS = []CSR{
	{Number: 0x000, Privilege: URW, Name: "ustatus", Description: "User status register"},
	// 0x004 URW uie User interrupt-enable register.
	// 0x005 URW utvec User trap handler base address.

	{Number: 0xC00, Privilege: URO, Name: "cycle", Description: "Cycle counter for RDCYCLE instruction"},
}

func main() {
	var val memio.UintN
	err := memio.Read(CSRS[1].Number, val)
	fmt.Printf("cycle %v %v", val, err)
}
