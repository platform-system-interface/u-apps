package main

import (
	"github.com/u-root/u-root/pkg/msr"
)

const (
	EFER                           msr.MSR = 0xc0000080
	STAR                           msr.MSR = 0xc0000081
	LSTAR                          msr.MSR = 0xc0000082
	CSTAR                          msr.MSR = 0xc0000083
	SYSCALL_MASK                   msr.MSR = 0xc0000084
	FS_BASE                        msr.MSR = 0xc0000100
	GS_BASE                        msr.MSR = 0xc0000101
	KERNEL_GS_BASE                 msr.MSR = 0xc0000102
	TSC_AUX                        msr.MSR = 0xc0000103
	HV_X64_RESET                   msr.MSR = 0x40000003
	HV_X64_TSC_FREQUENCY           msr.MSR = 0x40000022
	HV_X64_APIC_FREQUENCY          msr.MSR = 0x40000023
	HV_X64_REENLIGHTENMENT_CONTROL msr.MSR = 0x40000106
	HV_X64_TSC_EMULATION_CONTROL   msr.MSR = 0x40000107
	HV_X64_TSC_EMULATION_STATUS    msr.MSR = 0x40000108
)

var MSRS = []msr.MSRVal{
	/* https://github.com/bytedance/kvm-utils/blob/master/microbenchmark/msr-bench/msr-index.h */
	{Name: "EFER", Addr: EFER, Set: 0},
	{Name: "STAR", Addr: STAR, Set: 0},
	{Name: "LSTAR", Addr: LSTAR, Set: 0},
	{Name: "CSTAR", Addr: CSTAR, Set: 0},
	{Name: "SYSCALL_MASK", Addr: SYSCALL_MASK, Set: 0},
	{Name: "FS_BASE", Addr: FS_BASE, Set: 0},
	{Name: "GS_BASE", Addr: GS_BASE, Set: 0},
	{Name: "KERNEL_GS_BASE", Addr: KERNEL_GS_BASE, Set: 0},
	{Name: "TSC_AUX", Addr: TSC_AUX, Set: 0},
	/* https://readthedocs.org/projects/qemu/downloads/pdf/latest/ p280 */
	{Name: "HV_X64_RESET", Addr: HV_X64_RESET, Set: 0},
	{Name: "HV_X64_TSC_FREQUENCY", Addr: HV_X64_TSC_FREQUENCY, Set: 0},
	{Name: "HV_X64_APIC_FREQUENCY", Addr: HV_X64_APIC_FREQUENCY, Set: 0},
	{Name: "HV_X64_REENLIGHTENMENT_CONTROL", Addr: HV_X64_REENLIGHTENMENT_CONTROL, Set: 0},
	{Name: "HV_X64_TSC_EMULATION_CONTROL  ", Addr: HV_X64_TSC_EMULATION_CONTROL, Set: 0},
	{Name: "HV_X64_TSC_EMULATION_STATUS  ", Addr: HV_X64_TSC_EMULATION_STATUS, Set: 0},
}
