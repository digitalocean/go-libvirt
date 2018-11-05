// Copyright 2018 The go-libvirt Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// WARNING: This file has automatically been generated
// by https://git.io/c-for-go. DO NOT EDIT.

package libvirt

const (
	// ExportVar as defined in libvirt/libvirt-common.h:58
	ExportVar = 0
	// TypedParamFieldLength as defined in libvirt/libvirt-common.h:171
	TypedParamFieldLength = 80
	// SecurityLabelBuflen as defined in libvirt/libvirt-host.h:85
	SecurityLabelBuflen = 4097
	// SecurityModelBuflen as defined in libvirt/libvirt-host.h:113
	SecurityModelBuflen = 257
	// SecurityDoiBuflen as defined in libvirt/libvirt-host.h:120
	SecurityDoiBuflen = 257
	// NodeCPUStatsFieldLength as defined in libvirt/libvirt-host.h:177
	NodeCPUStatsFieldLength = 80
	// NodeCPUStatsKernel as defined in libvirt/libvirt-host.h:194
	NodeCPUStatsKernel = "kernel"
	// NodeCPUStatsUser as defined in libvirt/libvirt-host.h:202
	NodeCPUStatsUser = "user"
	// NodeCPUStatsIdle as defined in libvirt/libvirt-host.h:210
	NodeCPUStatsIdle = "idle"
	// NodeCPUStatsIowait as defined in libvirt/libvirt-host.h:218
	NodeCPUStatsIowait = "iowait"
	// NodeCPUStatsIntr as defined in libvirt/libvirt-host.h:226
	NodeCPUStatsIntr = "intr"
	// NodeCPUStatsUtilization as defined in libvirt/libvirt-host.h:235
	NodeCPUStatsUtilization = "utilization"
	// NodeMemoryStatsFieldLength as defined in libvirt/libvirt-host.h:255
	NodeMemoryStatsFieldLength = 80
	// NodeMemoryStatsTotal as defined in libvirt/libvirt-host.h:272
	NodeMemoryStatsTotal = "total"
	// NodeMemoryStatsFree as defined in libvirt/libvirt-host.h:281
	NodeMemoryStatsFree = "free"
	// NodeMemoryStatsBuffers as defined in libvirt/libvirt-host.h:289
	NodeMemoryStatsBuffers = "buffers"
	// NodeMemoryStatsCached as defined in libvirt/libvirt-host.h:297
	NodeMemoryStatsCached = "cached"
	// NodeMemorySharedPagesToScan as defined in libvirt/libvirt-host.h:318
	NodeMemorySharedPagesToScan = "shm_pages_to_scan"
	// NodeMemorySharedSleepMillisecs as defined in libvirt/libvirt-host.h:326
	NodeMemorySharedSleepMillisecs = "shm_sleep_millisecs"
	// NodeMemorySharedPagesShared as defined in libvirt/libvirt-host.h:334
	NodeMemorySharedPagesShared = "shm_pages_shared"
	// NodeMemorySharedPagesSharing as defined in libvirt/libvirt-host.h:342
	NodeMemorySharedPagesSharing = "shm_pages_sharing"
	// NodeMemorySharedPagesUnshared as defined in libvirt/libvirt-host.h:350
	NodeMemorySharedPagesUnshared = "shm_pages_unshared"
	// NodeMemorySharedPagesVolatile as defined in libvirt/libvirt-host.h:358
	NodeMemorySharedPagesVolatile = "shm_pages_volatile"
	// NodeMemorySharedFullScans as defined in libvirt/libvirt-host.h:366
	NodeMemorySharedFullScans = "shm_full_scans"
	// NodeMemorySharedMergeAcrossNodes as defined in libvirt/libvirt-host.h:378
	NodeMemorySharedMergeAcrossNodes = "shm_merge_across_nodes"
	// NodeSevPdh as defined in libvirt/libvirt-host.h:446
	NodeSevPdh = "pdh"
	// NodeSevCertChain as defined in libvirt/libvirt-host.h:455
	NodeSevCertChain = "cert-chain"
	// NodeSevCbitpos as defined in libvirt/libvirt-host.h:462
	NodeSevCbitpos = "cbitpos"
	// NodeSevReducedPhysBits as defined in libvirt/libvirt-host.h:470
	NodeSevReducedPhysBits = "reduced-phys-bits"
	// UUIDBuflen as defined in libvirt/libvirt-host.h:555
	UUIDBuflen = 16
	// UUIDStringBuflen as defined in libvirt/libvirt-host.h:564
	UUIDStringBuflen = 37
	// DomainSchedulerCPUShares as defined in libvirt/libvirt-domain.h:315
	DomainSchedulerCPUShares = "cpu_shares"
	// DomainSchedulerGlobalPeriod as defined in libvirt/libvirt-domain.h:323
	DomainSchedulerGlobalPeriod = "global_period"
	// DomainSchedulerGlobalQuota as defined in libvirt/libvirt-domain.h:331
	DomainSchedulerGlobalQuota = "global_quota"
	// DomainSchedulerVCPUPeriod as defined in libvirt/libvirt-domain.h:339
	DomainSchedulerVCPUPeriod = "vcpu_period"
	// DomainSchedulerVCPUQuota as defined in libvirt/libvirt-domain.h:347
	DomainSchedulerVCPUQuota = "vcpu_quota"
	// DomainSchedulerEmulatorPeriod as defined in libvirt/libvirt-domain.h:356
	DomainSchedulerEmulatorPeriod = "emulator_period"
	// DomainSchedulerEmulatorQuota as defined in libvirt/libvirt-domain.h:365
	DomainSchedulerEmulatorQuota = "emulator_quota"
	// DomainSchedulerIothreadPeriod as defined in libvirt/libvirt-domain.h:373
	DomainSchedulerIothreadPeriod = "iothread_period"
	// DomainSchedulerIothreadQuota as defined in libvirt/libvirt-domain.h:381
	DomainSchedulerIothreadQuota = "iothread_quota"
	// DomainSchedulerWeight as defined in libvirt/libvirt-domain.h:389
	DomainSchedulerWeight = "weight"
	// DomainSchedulerCap as defined in libvirt/libvirt-domain.h:397
	DomainSchedulerCap = "cap"
	// DomainSchedulerReservation as defined in libvirt/libvirt-domain.h:405
	DomainSchedulerReservation = "reservation"
	// DomainSchedulerLimit as defined in libvirt/libvirt-domain.h:413
	DomainSchedulerLimit = "limit"
	// DomainSchedulerShares as defined in libvirt/libvirt-domain.h:421
	DomainSchedulerShares = "shares"
	// DomainBlockStatsFieldLength as defined in libvirt/libvirt-domain.h:479
	DomainBlockStatsFieldLength = 80
	// DomainBlockStatsReadBytes as defined in libvirt/libvirt-domain.h:487
	DomainBlockStatsReadBytes = "rd_bytes"
	// DomainBlockStatsReadReq as defined in libvirt/libvirt-domain.h:495
	DomainBlockStatsReadReq = "rd_operations"
	// DomainBlockStatsReadTotalTimes as defined in libvirt/libvirt-domain.h:503
	DomainBlockStatsReadTotalTimes = "rd_total_times"
	// DomainBlockStatsWriteBytes as defined in libvirt/libvirt-domain.h:511
	DomainBlockStatsWriteBytes = "wr_bytes"
	// DomainBlockStatsWriteReq as defined in libvirt/libvirt-domain.h:519
	DomainBlockStatsWriteReq = "wr_operations"
	// DomainBlockStatsWriteTotalTimes as defined in libvirt/libvirt-domain.h:527
	DomainBlockStatsWriteTotalTimes = "wr_total_times"
	// DomainBlockStatsFlushReq as defined in libvirt/libvirt-domain.h:535
	DomainBlockStatsFlushReq = "flush_operations"
	// DomainBlockStatsFlushTotalTimes as defined in libvirt/libvirt-domain.h:543
	DomainBlockStatsFlushTotalTimes = "flush_total_times"
	// DomainBlockStatsErrs as defined in libvirt/libvirt-domain.h:550
	DomainBlockStatsErrs = "errs"
	// MigrateParamURI as defined in libvirt/libvirt-domain.h:849
	MigrateParamURI = "migrate_uri"
	// MigrateParamDestName as defined in libvirt/libvirt-domain.h:859
	MigrateParamDestName = "destination_name"
	// MigrateParamDestXML as defined in libvirt/libvirt-domain.h:878
	MigrateParamDestXML = "destination_xml"
	// MigrateParamPersistXML as defined in libvirt/libvirt-domain.h:893
	MigrateParamPersistXML = "persistent_xml"
	// MigrateParamBandwidth as defined in libvirt/libvirt-domain.h:903
	MigrateParamBandwidth = "bandwidth"
	// MigrateParamGraphicsURI as defined in libvirt/libvirt-domain.h:924
	MigrateParamGraphicsURI = "graphics_uri"
	// MigrateParamListenAddress as defined in libvirt/libvirt-domain.h:935
	MigrateParamListenAddress = "listen_address"
	// MigrateParamMigrateDisks as defined in libvirt/libvirt-domain.h:944
	MigrateParamMigrateDisks = "migrate_disks"
	// MigrateParamDisksPort as defined in libvirt/libvirt-domain.h:954
	MigrateParamDisksPort = "disks_port"
	// MigrateParamCompression as defined in libvirt/libvirt-domain.h:964
	MigrateParamCompression = "compression"
	// MigrateParamCompressionMtLevel as defined in libvirt/libvirt-domain.h:973
	MigrateParamCompressionMtLevel = "compression.mt.level"
	// MigrateParamCompressionMtThreads as defined in libvirt/libvirt-domain.h:981
	MigrateParamCompressionMtThreads = "compression.mt.threads"
	// MigrateParamCompressionMtDthreads as defined in libvirt/libvirt-domain.h:989
	MigrateParamCompressionMtDthreads = "compression.mt.dthreads"
	// MigrateParamCompressionXbzrleCache as defined in libvirt/libvirt-domain.h:997
	MigrateParamCompressionXbzrleCache = "compression.xbzrle.cache"
	// MigrateParamAutoConvergeInitial as defined in libvirt/libvirt-domain.h:1006
	MigrateParamAutoConvergeInitial = "auto_converge.initial"
	// MigrateParamAutoConvergeIncrement as defined in libvirt/libvirt-domain.h:1016
	MigrateParamAutoConvergeIncrement = "auto_converge.increment"
	// DomainCPUStatsCputime as defined in libvirt/libvirt-domain.h:1269
	DomainCPUStatsCputime = "cpu_time"
	// DomainCPUStatsUsertime as defined in libvirt/libvirt-domain.h:1275
	DomainCPUStatsUsertime = "user_time"
	// DomainCPUStatsSystemtime as defined in libvirt/libvirt-domain.h:1281
	DomainCPUStatsSystemtime = "system_time"
	// DomainCPUStatsVcputime as defined in libvirt/libvirt-domain.h:1288
	DomainCPUStatsVcputime = "vcpu_time"
	// DomainBlkioWeight as defined in libvirt/libvirt-domain.h:1317
	DomainBlkioWeight = "weight"
	// DomainBlkioDeviceWeight as defined in libvirt/libvirt-domain.h:1327
	DomainBlkioDeviceWeight = "device_weight"
	// DomainBlkioDeviceReadIops as defined in libvirt/libvirt-domain.h:1338
	DomainBlkioDeviceReadIops = "device_read_iops_sec"
	// DomainBlkioDeviceWriteIops as defined in libvirt/libvirt-domain.h:1349
	DomainBlkioDeviceWriteIops = "device_write_iops_sec"
	// DomainBlkioDeviceReadBps as defined in libvirt/libvirt-domain.h:1360
	DomainBlkioDeviceReadBps = "device_read_bytes_sec"
	// DomainBlkioDeviceWriteBps as defined in libvirt/libvirt-domain.h:1371
	DomainBlkioDeviceWriteBps = "device_write_bytes_sec"
	// DomainMemoryParamUnlimited as defined in libvirt/libvirt-domain.h:1390
	DomainMemoryParamUnlimited = 9007199254740991
	// DomainMemoryHardLimit as defined in libvirt/libvirt-domain.h:1399
	DomainMemoryHardLimit = "hard_limit"
	// DomainMemorySoftLimit as defined in libvirt/libvirt-domain.h:1408
	DomainMemorySoftLimit = "soft_limit"
	// DomainMemoryMinGuarantee as defined in libvirt/libvirt-domain.h:1417
	DomainMemoryMinGuarantee = "min_guarantee"
	// DomainMemorySwapHardLimit as defined in libvirt/libvirt-domain.h:1427
	DomainMemorySwapHardLimit = "swap_hard_limit"
	// DomainNumaNodeset as defined in libvirt/libvirt-domain.h:1472
	DomainNumaNodeset = "numa_nodeset"
	// DomainNumaMode as defined in libvirt/libvirt-domain.h:1480
	DomainNumaMode = "numa_mode"
	// DomainBandwidthInAverage as defined in libvirt/libvirt-domain.h:1592
	DomainBandwidthInAverage = "inbound.average"
	// DomainBandwidthInPeak as defined in libvirt/libvirt-domain.h:1599
	DomainBandwidthInPeak = "inbound.peak"
	// DomainBandwidthInBurst as defined in libvirt/libvirt-domain.h:1606
	DomainBandwidthInBurst = "inbound.burst"
	// DomainBandwidthInFloor as defined in libvirt/libvirt-domain.h:1613
	DomainBandwidthInFloor = "inbound.floor"
	// DomainBandwidthOutAverage as defined in libvirt/libvirt-domain.h:1620
	DomainBandwidthOutAverage = "outbound.average"
	// DomainBandwidthOutPeak as defined in libvirt/libvirt-domain.h:1627
	DomainBandwidthOutPeak = "outbound.peak"
	// DomainBandwidthOutBurst as defined in libvirt/libvirt-domain.h:1634
	DomainBandwidthOutBurst = "outbound.burst"
	// PerfParamCmt as defined in libvirt/libvirt-domain.h:2095
	PerfParamCmt = "cmt"
	// PerfParamMbmt as defined in libvirt/libvirt-domain.h:2106
	PerfParamMbmt = "mbmt"
	// PerfParamMbml as defined in libvirt/libvirt-domain.h:2116
	PerfParamMbml = "mbml"
	// PerfParamCacheMisses as defined in libvirt/libvirt-domain.h:2126
	PerfParamCacheMisses = "cache_misses"
	// PerfParamCacheReferences as defined in libvirt/libvirt-domain.h:2136
	PerfParamCacheReferences = "cache_references"
	// PerfParamInstructions as defined in libvirt/libvirt-domain.h:2146
	PerfParamInstructions = "instructions"
	// PerfParamCPUCycles as defined in libvirt/libvirt-domain.h:2156
	PerfParamCPUCycles = "cpu_cycles"
	// PerfParamBranchInstructions as defined in libvirt/libvirt-domain.h:2166
	PerfParamBranchInstructions = "branch_instructions"
	// PerfParamBranchMisses as defined in libvirt/libvirt-domain.h:2176
	PerfParamBranchMisses = "branch_misses"
	// PerfParamBusCycles as defined in libvirt/libvirt-domain.h:2186
	PerfParamBusCycles = "bus_cycles"
	// PerfParamStalledCyclesFrontend as defined in libvirt/libvirt-domain.h:2197
	PerfParamStalledCyclesFrontend = "stalled_cycles_frontend"
	// PerfParamStalledCyclesBackend as defined in libvirt/libvirt-domain.h:2208
	PerfParamStalledCyclesBackend = "stalled_cycles_backend"
	// PerfParamRefCPUCycles as defined in libvirt/libvirt-domain.h:2219
	PerfParamRefCPUCycles = "ref_cpu_cycles"
	// PerfParamCPUClock as defined in libvirt/libvirt-domain.h:2230
	PerfParamCPUClock = "cpu_clock"
	// PerfParamTaskClock as defined in libvirt/libvirt-domain.h:2241
	PerfParamTaskClock = "task_clock"
	// PerfParamPageFaults as defined in libvirt/libvirt-domain.h:2251
	PerfParamPageFaults = "page_faults"
	// PerfParamContextSwitches as defined in libvirt/libvirt-domain.h:2261
	PerfParamContextSwitches = "context_switches"
	// PerfParamCPUMigrations as defined in libvirt/libvirt-domain.h:2271
	PerfParamCPUMigrations = "cpu_migrations"
	// PerfParamPageFaultsMin as defined in libvirt/libvirt-domain.h:2281
	PerfParamPageFaultsMin = "page_faults_min"
	// PerfParamPageFaultsMaj as defined in libvirt/libvirt-domain.h:2291
	PerfParamPageFaultsMaj = "page_faults_maj"
	// PerfParamAlignmentFaults as defined in libvirt/libvirt-domain.h:2301
	PerfParamAlignmentFaults = "alignment_faults"
	// PerfParamEmulationFaults as defined in libvirt/libvirt-domain.h:2311
	PerfParamEmulationFaults = "emulation_faults"
	// DomainBlockCopyBandwidth as defined in libvirt/libvirt-domain.h:2475
	DomainBlockCopyBandwidth = "bandwidth"
	// DomainBlockCopyGranularity as defined in libvirt/libvirt-domain.h:2486
	DomainBlockCopyGranularity = "granularity"
	// DomainBlockCopyBufSize as defined in libvirt/libvirt-domain.h:2495
	DomainBlockCopyBufSize = "buf-size"
	// DomainBlockIotuneTotalBytesSec as defined in libvirt/libvirt-domain.h:2536
	DomainBlockIotuneTotalBytesSec = "total_bytes_sec"
	// DomainBlockIotuneReadBytesSec as defined in libvirt/libvirt-domain.h:2544
	DomainBlockIotuneReadBytesSec = "read_bytes_sec"
	// DomainBlockIotuneWriteBytesSec as defined in libvirt/libvirt-domain.h:2552
	DomainBlockIotuneWriteBytesSec = "write_bytes_sec"
	// DomainBlockIotuneTotalIopsSec as defined in libvirt/libvirt-domain.h:2560
	DomainBlockIotuneTotalIopsSec = "total_iops_sec"
	// DomainBlockIotuneReadIopsSec as defined in libvirt/libvirt-domain.h:2568
	DomainBlockIotuneReadIopsSec = "read_iops_sec"
	// DomainBlockIotuneWriteIopsSec as defined in libvirt/libvirt-domain.h:2575
	DomainBlockIotuneWriteIopsSec = "write_iops_sec"
	// DomainBlockIotuneTotalBytesSecMax as defined in libvirt/libvirt-domain.h:2583
	DomainBlockIotuneTotalBytesSecMax = "total_bytes_sec_max"
	// DomainBlockIotuneReadBytesSecMax as defined in libvirt/libvirt-domain.h:2591
	DomainBlockIotuneReadBytesSecMax = "read_bytes_sec_max"
	// DomainBlockIotuneWriteBytesSecMax as defined in libvirt/libvirt-domain.h:2599
	DomainBlockIotuneWriteBytesSecMax = "write_bytes_sec_max"
	// DomainBlockIotuneTotalIopsSecMax as defined in libvirt/libvirt-domain.h:2607
	DomainBlockIotuneTotalIopsSecMax = "total_iops_sec_max"
	// DomainBlockIotuneReadIopsSecMax as defined in libvirt/libvirt-domain.h:2615
	DomainBlockIotuneReadIopsSecMax = "read_iops_sec_max"
	// DomainBlockIotuneWriteIopsSecMax as defined in libvirt/libvirt-domain.h:2622
	DomainBlockIotuneWriteIopsSecMax = "write_iops_sec_max"
	// DomainBlockIotuneTotalBytesSecMaxLength as defined in libvirt/libvirt-domain.h:2630
	DomainBlockIotuneTotalBytesSecMaxLength = "total_bytes_sec_max_length"
	// DomainBlockIotuneReadBytesSecMaxLength as defined in libvirt/libvirt-domain.h:2638
	DomainBlockIotuneReadBytesSecMaxLength = "read_bytes_sec_max_length"
	// DomainBlockIotuneWriteBytesSecMaxLength as defined in libvirt/libvirt-domain.h:2646
	DomainBlockIotuneWriteBytesSecMaxLength = "write_bytes_sec_max_length"
	// DomainBlockIotuneTotalIopsSecMaxLength as defined in libvirt/libvirt-domain.h:2654
	DomainBlockIotuneTotalIopsSecMaxLength = "total_iops_sec_max_length"
	// DomainBlockIotuneReadIopsSecMaxLength as defined in libvirt/libvirt-domain.h:2662
	DomainBlockIotuneReadIopsSecMaxLength = "read_iops_sec_max_length"
	// DomainBlockIotuneWriteIopsSecMaxLength as defined in libvirt/libvirt-domain.h:2670
	DomainBlockIotuneWriteIopsSecMaxLength = "write_iops_sec_max_length"
	// DomainBlockIotuneSizeIopsSec as defined in libvirt/libvirt-domain.h:2677
	DomainBlockIotuneSizeIopsSec = "size_iops_sec"
	// DomainBlockIotuneGroupName as defined in libvirt/libvirt-domain.h:2684
	DomainBlockIotuneGroupName = "group_name"
	// KeycodeSetRfb as defined in libvirt/libvirt-domain.h:2765
	KeycodeSetRfb = 0
	// DomainSendKeyMaxKeys as defined in libvirt/libvirt-domain.h:2772
	DomainSendKeyMaxKeys = 16
	// DomainJobOperationStr as defined in libvirt/libvirt-domain.h:3184
	DomainJobOperationStr = "operation"
	// DomainJobTimeElapsed as defined in libvirt/libvirt-domain.h:3194
	DomainJobTimeElapsed = "time_elapsed"
	// DomainJobTimeElapsedNet as defined in libvirt/libvirt-domain.h:3204
	DomainJobTimeElapsedNet = "time_elapsed_net"
	// DomainJobTimeRemaining as defined in libvirt/libvirt-domain.h:3214
	DomainJobTimeRemaining = "time_remaining"
	// DomainJobDowntime as defined in libvirt/libvirt-domain.h:3224
	DomainJobDowntime = "downtime"
	// DomainJobDowntimeNet as defined in libvirt/libvirt-domain.h:3233
	DomainJobDowntimeNet = "downtime_net"
	// DomainJobSetupTime as defined in libvirt/libvirt-domain.h:3242
	DomainJobSetupTime = "setup_time"
	// DomainJobDataTotal as defined in libvirt/libvirt-domain.h:3257
	DomainJobDataTotal = "data_total"
	// DomainJobDataProcessed as defined in libvirt/libvirt-domain.h:3267
	DomainJobDataProcessed = "data_processed"
	// DomainJobDataRemaining as defined in libvirt/libvirt-domain.h:3277
	DomainJobDataRemaining = "data_remaining"
	// DomainJobMemoryTotal as defined in libvirt/libvirt-domain.h:3287
	DomainJobMemoryTotal = "memory_total"
	// DomainJobMemoryProcessed as defined in libvirt/libvirt-domain.h:3297
	DomainJobMemoryProcessed = "memory_processed"
	// DomainJobMemoryRemaining as defined in libvirt/libvirt-domain.h:3307
	DomainJobMemoryRemaining = "memory_remaining"
	// DomainJobMemoryConstant as defined in libvirt/libvirt-domain.h:3319
	DomainJobMemoryConstant = "memory_constant"
	// DomainJobMemoryNormal as defined in libvirt/libvirt-domain.h:3329
	DomainJobMemoryNormal = "memory_normal"
	// DomainJobMemoryNormalBytes as defined in libvirt/libvirt-domain.h:3339
	DomainJobMemoryNormalBytes = "memory_normal_bytes"
	// DomainJobMemoryBps as defined in libvirt/libvirt-domain.h:3347
	DomainJobMemoryBps = "memory_bps"
	// DomainJobMemoryDirtyRate as defined in libvirt/libvirt-domain.h:3355
	DomainJobMemoryDirtyRate = "memory_dirty_rate"
	// DomainJobMemoryPageSize as defined in libvirt/libvirt-domain.h:3366
	DomainJobMemoryPageSize = "memory_page_size"
	// DomainJobMemoryIteration as defined in libvirt/libvirt-domain.h:3377
	DomainJobMemoryIteration = "memory_iteration"
	// DomainJobDiskTotal as defined in libvirt/libvirt-domain.h:3387
	DomainJobDiskTotal = "disk_total"
	// DomainJobDiskProcessed as defined in libvirt/libvirt-domain.h:3397
	DomainJobDiskProcessed = "disk_processed"
	// DomainJobDiskRemaining as defined in libvirt/libvirt-domain.h:3407
	DomainJobDiskRemaining = "disk_remaining"
	// DomainJobDiskBps as defined in libvirt/libvirt-domain.h:3415
	DomainJobDiskBps = "disk_bps"
	// DomainJobCompressionCache as defined in libvirt/libvirt-domain.h:3424
	DomainJobCompressionCache = "compression_cache"
	// DomainJobCompressionBytes as defined in libvirt/libvirt-domain.h:3432
	DomainJobCompressionBytes = "compression_bytes"
	// DomainJobCompressionPages as defined in libvirt/libvirt-domain.h:3440
	DomainJobCompressionPages = "compression_pages"
	// DomainJobCompressionCacheMisses as defined in libvirt/libvirt-domain.h:3449
	DomainJobCompressionCacheMisses = "compression_cache_misses"
	// DomainJobCompressionOverflow as defined in libvirt/libvirt-domain.h:3459
	DomainJobCompressionOverflow = "compression_overflow"
	// DomainJobAutoConvergeThrottle as defined in libvirt/libvirt-domain.h:3468
	DomainJobAutoConvergeThrottle = "auto_converge_throttle"
	// DomainTunableCPUVcpupin as defined in libvirt/libvirt-domain.h:4022
	DomainTunableCPUVcpupin = "cputune.vcpupin%u"
	// DomainTunableCPUEmulatorpin as defined in libvirt/libvirt-domain.h:4030
	DomainTunableCPUEmulatorpin = "cputune.emulatorpin"
	// DomainTunableCPUIothreadspin as defined in libvirt/libvirt-domain.h:4039
	DomainTunableCPUIothreadspin = "cputune.iothreadpin%u"
	// DomainTunableCPUCpuShares as defined in libvirt/libvirt-domain.h:4047
	DomainTunableCPUCpuShares = "cputune.cpu_shares"
	// DomainTunableCPUGlobalPeriod as defined in libvirt/libvirt-domain.h:4055
	DomainTunableCPUGlobalPeriod = "cputune.global_period"
	// DomainTunableCPUGlobalQuota as defined in libvirt/libvirt-domain.h:4063
	DomainTunableCPUGlobalQuota = "cputune.global_quota"
	// DomainTunableCPUVCPUPeriod as defined in libvirt/libvirt-domain.h:4071
	DomainTunableCPUVCPUPeriod = "cputune.vcpu_period"
	// DomainTunableCPUVCPUQuota as defined in libvirt/libvirt-domain.h:4079
	DomainTunableCPUVCPUQuota = "cputune.vcpu_quota"
	// DomainTunableCPUEmulatorPeriod as defined in libvirt/libvirt-domain.h:4088
	DomainTunableCPUEmulatorPeriod = "cputune.emulator_period"
	// DomainTunableCPUEmulatorQuota as defined in libvirt/libvirt-domain.h:4097
	DomainTunableCPUEmulatorQuota = "cputune.emulator_quota"
	// DomainTunableCPUIothreadPeriod as defined in libvirt/libvirt-domain.h:4105
	DomainTunableCPUIothreadPeriod = "cputune.iothread_period"
	// DomainTunableCPUIothreadQuota as defined in libvirt/libvirt-domain.h:4113
	DomainTunableCPUIothreadQuota = "cputune.iothread_quota"
	// DomainTunableBlkdevDisk as defined in libvirt/libvirt-domain.h:4121
	DomainTunableBlkdevDisk = "blkdeviotune.disk"
	// DomainTunableBlkdevTotalBytesSec as defined in libvirt/libvirt-domain.h:4129
	DomainTunableBlkdevTotalBytesSec = "blkdeviotune.total_bytes_sec"
	// DomainTunableBlkdevReadBytesSec as defined in libvirt/libvirt-domain.h:4137
	DomainTunableBlkdevReadBytesSec = "blkdeviotune.read_bytes_sec"
	// DomainTunableBlkdevWriteBytesSec as defined in libvirt/libvirt-domain.h:4145
	DomainTunableBlkdevWriteBytesSec = "blkdeviotune.write_bytes_sec"
	// DomainTunableBlkdevTotalIopsSec as defined in libvirt/libvirt-domain.h:4153
	DomainTunableBlkdevTotalIopsSec = "blkdeviotune.total_iops_sec"
	// DomainTunableBlkdevReadIopsSec as defined in libvirt/libvirt-domain.h:4161
	DomainTunableBlkdevReadIopsSec = "blkdeviotune.read_iops_sec"
	// DomainTunableBlkdevWriteIopsSec as defined in libvirt/libvirt-domain.h:4169
	DomainTunableBlkdevWriteIopsSec = "blkdeviotune.write_iops_sec"
	// DomainTunableBlkdevTotalBytesSecMax as defined in libvirt/libvirt-domain.h:4177
	DomainTunableBlkdevTotalBytesSecMax = "blkdeviotune.total_bytes_sec_max"
	// DomainTunableBlkdevReadBytesSecMax as defined in libvirt/libvirt-domain.h:4185
	DomainTunableBlkdevReadBytesSecMax = "blkdeviotune.read_bytes_sec_max"
	// DomainTunableBlkdevWriteBytesSecMax as defined in libvirt/libvirt-domain.h:4193
	DomainTunableBlkdevWriteBytesSecMax = "blkdeviotune.write_bytes_sec_max"
	// DomainTunableBlkdevTotalIopsSecMax as defined in libvirt/libvirt-domain.h:4201
	DomainTunableBlkdevTotalIopsSecMax = "blkdeviotune.total_iops_sec_max"
	// DomainTunableBlkdevReadIopsSecMax as defined in libvirt/libvirt-domain.h:4209
	DomainTunableBlkdevReadIopsSecMax = "blkdeviotune.read_iops_sec_max"
	// DomainTunableBlkdevWriteIopsSecMax as defined in libvirt/libvirt-domain.h:4217
	DomainTunableBlkdevWriteIopsSecMax = "blkdeviotune.write_iops_sec_max"
	// DomainTunableBlkdevSizeIopsSec as defined in libvirt/libvirt-domain.h:4225
	DomainTunableBlkdevSizeIopsSec = "blkdeviotune.size_iops_sec"
	// DomainTunableBlkdevGroupName as defined in libvirt/libvirt-domain.h:4233
	DomainTunableBlkdevGroupName = "blkdeviotune.group_name"
	// DomainTunableBlkdevTotalBytesSecMaxLength as defined in libvirt/libvirt-domain.h:4242
	DomainTunableBlkdevTotalBytesSecMaxLength = "blkdeviotune.total_bytes_sec_max_length"
	// DomainTunableBlkdevReadBytesSecMaxLength as defined in libvirt/libvirt-domain.h:4251
	DomainTunableBlkdevReadBytesSecMaxLength = "blkdeviotune.read_bytes_sec_max_length"
	// DomainTunableBlkdevWriteBytesSecMaxLength as defined in libvirt/libvirt-domain.h:4260
	DomainTunableBlkdevWriteBytesSecMaxLength = "blkdeviotune.write_bytes_sec_max_length"
	// DomainTunableBlkdevTotalIopsSecMaxLength as defined in libvirt/libvirt-domain.h:4269
	DomainTunableBlkdevTotalIopsSecMaxLength = "blkdeviotune.total_iops_sec_max_length"
	// DomainTunableBlkdevReadIopsSecMaxLength as defined in libvirt/libvirt-domain.h:4278
	DomainTunableBlkdevReadIopsSecMaxLength = "blkdeviotune.read_iops_sec_max_length"
	// DomainTunableBlkdevWriteIopsSecMaxLength as defined in libvirt/libvirt-domain.h:4287
	DomainTunableBlkdevWriteIopsSecMaxLength = "blkdeviotune.write_iops_sec_max_length"
	// DomainSchedFieldLength as defined in libvirt/libvirt-domain.h:4575
	DomainSchedFieldLength = 80
	// DomainBlkioFieldLength as defined in libvirt/libvirt-domain.h:4619
	DomainBlkioFieldLength = 80
	// DomainMemoryFieldLength as defined in libvirt/libvirt-domain.h:4663
	DomainMemoryFieldLength = 80
	// DomainLaunchSecuritySevMeasurement as defined in libvirt/libvirt-domain.h:4789
	DomainLaunchSecuritySevMeasurement = "sev-measurement"
)

// ConnectCloseReason as declared in libvirt/libvirt-common.h:120
type ConnectCloseReason int32

// ConnectCloseReason enumeration from libvirt/libvirt-common.h:120
const (
	ConnectCloseReasonError     ConnectCloseReason = iota
	ConnectCloseReasonEOF       ConnectCloseReason = 1
	ConnectCloseReasonKeepalive ConnectCloseReason = 2
	ConnectCloseReasonClient    ConnectCloseReason = 3
)

// TypedParameterType as declared in libvirt/libvirt-common.h:139
type TypedParameterType int32

// TypedParameterType enumeration from libvirt/libvirt-common.h:139
const (
	TypedParamInt     TypedParameterType = 1
	TypedParamUint    TypedParameterType = 2
	TypedParamLlong   TypedParameterType = 3
	TypedParamUllong  TypedParameterType = 4
	TypedParamDouble  TypedParameterType = 5
	TypedParamBoolean TypedParameterType = 6
	TypedParamString  TypedParameterType = 7
)

// TypedParameterFlags as declared in libvirt/libvirt-common.h:164
type TypedParameterFlags int32

// TypedParameterFlags enumeration from libvirt/libvirt-common.h:164
const (
	TypedParamStringOkay TypedParameterFlags = 4
)

// NodeSuspendTarget as declared in libvirt/libvirt-host.h:62
type NodeSuspendTarget int32

// NodeSuspendTarget enumeration from libvirt/libvirt-host.h:62
const (
	NodeSuspendTargetMem    NodeSuspendTarget = iota
	NodeSuspendTargetDisk   NodeSuspendTarget = 1
	NodeSuspendTargetHybrid NodeSuspendTarget = 2
)

// NodeGetCPUStatsAllCPUs as declared in libvirt/libvirt-host.h:186
type NodeGetCPUStatsAllCPUs int32

// NodeGetCPUStatsAllCPUs enumeration from libvirt/libvirt-host.h:186
const (
	NodeCPUStatsAllCpus NodeGetCPUStatsAllCPUs = -1
)

// NodeGetMemoryStatsAllCells as declared in libvirt/libvirt-host.h:264
type NodeGetMemoryStatsAllCells int32

// NodeGetMemoryStatsAllCells enumeration from libvirt/libvirt-host.h:264
const (
	NodeMemoryStatsAllCells NodeGetMemoryStatsAllCells = -1
)

// ConnectFlags as declared in libvirt/libvirt-host.h:485
type ConnectFlags int32

// ConnectFlags enumeration from libvirt/libvirt-host.h:485
const (
	ConnectRo        ConnectFlags = 1
	ConnectNoAliases ConnectFlags = 2
)

// ConnectCredentialType as declared in libvirt/libvirt-host.h:502
type ConnectCredentialType int32

// ConnectCredentialType enumeration from libvirt/libvirt-host.h:502
const (
	CredUsername     ConnectCredentialType = 1
	CredAuthname     ConnectCredentialType = 2
	CredLanguage     ConnectCredentialType = 3
	CredCnonce       ConnectCredentialType = 4
	CredPassphrase   ConnectCredentialType = 5
	CredEchoprompt   ConnectCredentialType = 6
	CredNoechoprompt ConnectCredentialType = 7
	CredRealm        ConnectCredentialType = 8
	CredExternal     ConnectCredentialType = 9
)

// CPUCompareResult as declared in libvirt/libvirt-host.h:675
type CPUCompareResult int32

// CPUCompareResult enumeration from libvirt/libvirt-host.h:675
const (
	CPUCompareError        CPUCompareResult = -1
	CPUCompareIncompatible CPUCompareResult = 0
	CPUCompareIdentical    CPUCompareResult = 1
	CPUCompareSuperset     CPUCompareResult = 2
)

// ConnectCompareCPUFlags as declared in libvirt/libvirt-host.h:680
type ConnectCompareCPUFlags int32

// ConnectCompareCPUFlags enumeration from libvirt/libvirt-host.h:680
const (
	ConnectCompareCPUFailIncompatible ConnectCompareCPUFlags = 1
)

// ConnectBaselineCPUFlags as declared in libvirt/libvirt-host.h:706
type ConnectBaselineCPUFlags int32

// ConnectBaselineCPUFlags enumeration from libvirt/libvirt-host.h:706
const (
	ConnectBaselineCPUExpandFeatures ConnectBaselineCPUFlags = 1
	ConnectBaselineCPUMigratable     ConnectBaselineCPUFlags = 2
)

// NodeAllocPagesFlags as declared in libvirt/libvirt-host.h:736
type NodeAllocPagesFlags int32

// NodeAllocPagesFlags enumeration from libvirt/libvirt-host.h:736
const (
	NodeAllocPagesAdd NodeAllocPagesFlags = iota
	NodeAllocPagesSet NodeAllocPagesFlags = 1
)

// DomainState as declared in libvirt/libvirt-domain.h:71
type DomainState int32

// DomainState enumeration from libvirt/libvirt-domain.h:71
const (
	DomainNostate     DomainState = iota
	DomainRunning     DomainState = 1
	DomainBlocked     DomainState = 2
	DomainPaused      DomainState = 3
	DomainShutdown    DomainState = 4
	DomainShutoff     DomainState = 5
	DomainCrashed     DomainState = 6
	DomainPmsuspended DomainState = 7
)

// DomainNostateReason as declared in libvirt/libvirt-domain.h:79
type DomainNostateReason int32

// DomainNostateReason enumeration from libvirt/libvirt-domain.h:79
const (
	DomainNostateUnknown DomainNostateReason = iota
)

// DomainRunningReason as declared in libvirt/libvirt-domain.h:98
type DomainRunningReason int32

// DomainRunningReason enumeration from libvirt/libvirt-domain.h:98
const (
	DomainRunningUnknown           DomainRunningReason = iota
	DomainRunningBooted            DomainRunningReason = 1
	DomainRunningMigrated          DomainRunningReason = 2
	DomainRunningRestored          DomainRunningReason = 3
	DomainRunningFromSnapshot      DomainRunningReason = 4
	DomainRunningUnpaused          DomainRunningReason = 5
	DomainRunningMigrationCanceled DomainRunningReason = 6
	DomainRunningSaveCanceled      DomainRunningReason = 7
	DomainRunningWakeup            DomainRunningReason = 8
	DomainRunningCrashed           DomainRunningReason = 9
	DomainRunningPostcopy          DomainRunningReason = 10
)

// DomainBlockedReason as declared in libvirt/libvirt-domain.h:106
type DomainBlockedReason int32

// DomainBlockedReason enumeration from libvirt/libvirt-domain.h:106
const (
	DomainBlockedUnknown DomainBlockedReason = iota
)

// DomainPausedReason as declared in libvirt/libvirt-domain.h:127
type DomainPausedReason int32

// DomainPausedReason enumeration from libvirt/libvirt-domain.h:127
const (
	DomainPausedUnknown        DomainPausedReason = iota
	DomainPausedUser           DomainPausedReason = 1
	DomainPausedMigration      DomainPausedReason = 2
	DomainPausedSave           DomainPausedReason = 3
	DomainPausedDump           DomainPausedReason = 4
	DomainPausedIoerror        DomainPausedReason = 5
	DomainPausedWatchdog       DomainPausedReason = 6
	DomainPausedFromSnapshot   DomainPausedReason = 7
	DomainPausedShuttingDown   DomainPausedReason = 8
	DomainPausedSnapshot       DomainPausedReason = 9
	DomainPausedCrashed        DomainPausedReason = 10
	DomainPausedStartingUp     DomainPausedReason = 11
	DomainPausedPostcopy       DomainPausedReason = 12
	DomainPausedPostcopyFailed DomainPausedReason = 13
)

// DomainShutdownReason as declared in libvirt/libvirt-domain.h:136
type DomainShutdownReason int32

// DomainShutdownReason enumeration from libvirt/libvirt-domain.h:136
const (
	DomainShutdownUnknown DomainShutdownReason = iota
	DomainShutdownUser    DomainShutdownReason = 1
)

// DomainShutoffReason as declared in libvirt/libvirt-domain.h:151
type DomainShutoffReason int32

// DomainShutoffReason enumeration from libvirt/libvirt-domain.h:151
const (
	DomainShutoffUnknown      DomainShutoffReason = iota
	DomainShutoffShutdown     DomainShutoffReason = 1
	DomainShutoffDestroyed    DomainShutoffReason = 2
	DomainShutoffCrashed      DomainShutoffReason = 3
	DomainShutoffMigrated     DomainShutoffReason = 4
	DomainShutoffSaved        DomainShutoffReason = 5
	DomainShutoffFailed       DomainShutoffReason = 6
	DomainShutoffFromSnapshot DomainShutoffReason = 7
)

// DomainCrashedReason as declared in libvirt/libvirt-domain.h:160
type DomainCrashedReason int32

// DomainCrashedReason enumeration from libvirt/libvirt-domain.h:160
const (
	DomainCrashedUnknown  DomainCrashedReason = iota
	DomainCrashedPanicked DomainCrashedReason = 1
)

// DomainPMSuspendedReason as declared in libvirt/libvirt-domain.h:168
type DomainPMSuspendedReason int32

// DomainPMSuspendedReason enumeration from libvirt/libvirt-domain.h:168
const (
	DomainPmsuspendedUnknown DomainPMSuspendedReason = iota
)

// DomainPMSuspendedDiskReason as declared in libvirt/libvirt-domain.h:176
type DomainPMSuspendedDiskReason int32

// DomainPMSuspendedDiskReason enumeration from libvirt/libvirt-domain.h:176
const (
	DomainPmsuspendedDiskUnknown DomainPMSuspendedDiskReason = iota
)

// DomainControlState as declared in libvirt/libvirt-domain.h:196
type DomainControlState int32

// DomainControlState enumeration from libvirt/libvirt-domain.h:196
const (
	DomainControlOk       DomainControlState = iota
	DomainControlJob      DomainControlState = 1
	DomainControlOccupied DomainControlState = 2
	DomainControlError    DomainControlState = 3
)

// DomainControlErrorReason as declared in libvirt/libvirt-domain.h:216
type DomainControlErrorReason int32

// DomainControlErrorReason enumeration from libvirt/libvirt-domain.h:216
const (
	DomainControlErrorReasonNone     DomainControlErrorReason = iota
	DomainControlErrorReasonUnknown  DomainControlErrorReason = 1
	DomainControlErrorReasonMonitor  DomainControlErrorReason = 2
	DomainControlErrorReasonInternal DomainControlErrorReason = 3
)

// DomainModificationImpact as declared in libvirt/libvirt-domain.h:264
type DomainModificationImpact int32

// DomainModificationImpact enumeration from libvirt/libvirt-domain.h:264
const (
	DomainAffectCurrent DomainModificationImpact = iota
	DomainAffectLive    DomainModificationImpact = 1
	DomainAffectConfig  DomainModificationImpact = 2
)

// DomainCreateFlags as declared in libvirt/libvirt-domain.h:304
type DomainCreateFlags int32

// DomainCreateFlags enumeration from libvirt/libvirt-domain.h:304
const (
	DomainNone             DomainCreateFlags = iota
	DomainStartPaused      DomainCreateFlags = 1
	DomainStartAutodestroy DomainCreateFlags = 2
	DomainStartBypassCache DomainCreateFlags = 4
	DomainStartForceBoot   DomainCreateFlags = 8
	DomainStartValidate    DomainCreateFlags = 16
)

// DomainMemoryStatTags as declared in libvirt/libvirt-domain.h:647
type DomainMemoryStatTags int32

// DomainMemoryStatTags enumeration from libvirt/libvirt-domain.h:647
const (
	DomainMemoryStatSwapIn        DomainMemoryStatTags = iota
	DomainMemoryStatSwapOut       DomainMemoryStatTags = 1
	DomainMemoryStatMajorFault    DomainMemoryStatTags = 2
	DomainMemoryStatMinorFault    DomainMemoryStatTags = 3
	DomainMemoryStatUnused        DomainMemoryStatTags = 4
	DomainMemoryStatAvailable     DomainMemoryStatTags = 5
	DomainMemoryStatActualBalloon DomainMemoryStatTags = 6
	DomainMemoryStatRss           DomainMemoryStatTags = 7
	DomainMemoryStatUsable        DomainMemoryStatTags = 8
	DomainMemoryStatLastUpdate    DomainMemoryStatTags = 9
	DomainMemoryStatDiskCaches    DomainMemoryStatTags = 10
	DomainMemoryStatNr            DomainMemoryStatTags = 11
)

// DomainCoreDumpFlags as declared in libvirt/libvirt-domain.h:666
type DomainCoreDumpFlags int32

// DomainCoreDumpFlags enumeration from libvirt/libvirt-domain.h:666
const (
	DumpCrash       DomainCoreDumpFlags = 1
	DumpLive        DomainCoreDumpFlags = 2
	DumpBypassCache DomainCoreDumpFlags = 4
	DumpReset       DomainCoreDumpFlags = 8
	DumpMemoryOnly  DomainCoreDumpFlags = 16
)

// DomainCoreDumpFormat as declared in libvirt/libvirt-domain.h:689
type DomainCoreDumpFormat int32

// DomainCoreDumpFormat enumeration from libvirt/libvirt-domain.h:689
const (
	DomainCoreDumpFormatRaw         DomainCoreDumpFormat = iota
	DomainCoreDumpFormatKdumpZlib   DomainCoreDumpFormat = 1
	DomainCoreDumpFormatKdumpLzo    DomainCoreDumpFormat = 2
	DomainCoreDumpFormatKdumpSnappy DomainCoreDumpFormat = 3
)

// DomainMigrateFlags as declared in libvirt/libvirt-domain.h:833
type DomainMigrateFlags int32

// DomainMigrateFlags enumeration from libvirt/libvirt-domain.h:833
const (
	MigrateLive             DomainMigrateFlags = 1
	MigratePeer2peer        DomainMigrateFlags = 2
	MigrateTunnelled        DomainMigrateFlags = 4
	MigratePersistDest      DomainMigrateFlags = 8
	MigrateUndefineSource   DomainMigrateFlags = 16
	MigratePaused           DomainMigrateFlags = 32
	MigrateNonSharedDisk    DomainMigrateFlags = 64
	MigrateNonSharedInc     DomainMigrateFlags = 128
	MigrateChangeProtection DomainMigrateFlags = 256
	MigrateUnsafe           DomainMigrateFlags = 512
	MigrateOffline          DomainMigrateFlags = 1024
	MigrateCompressed       DomainMigrateFlags = 2048
	MigrateAbortOnError     DomainMigrateFlags = 4096
	MigrateAutoConverge     DomainMigrateFlags = 8192
	MigrateRdmaPinAll       DomainMigrateFlags = 16384
	MigratePostcopy         DomainMigrateFlags = 32768
	MigrateTLS              DomainMigrateFlags = 65536
)

// DomainShutdownFlagValues as declared in libvirt/libvirt-domain.h:1128
type DomainShutdownFlagValues int32

// DomainShutdownFlagValues enumeration from libvirt/libvirt-domain.h:1128
const (
	DomainShutdownDefault      DomainShutdownFlagValues = iota
	DomainShutdownAcpiPowerBtn DomainShutdownFlagValues = 1
	DomainShutdownGuestAgent   DomainShutdownFlagValues = 2
	DomainShutdownInitctl      DomainShutdownFlagValues = 4
	DomainShutdownSignal       DomainShutdownFlagValues = 8
	DomainShutdownParavirt     DomainShutdownFlagValues = 16
)

// DomainRebootFlagValues as declared in libvirt/libvirt-domain.h:1141
type DomainRebootFlagValues int32

// DomainRebootFlagValues enumeration from libvirt/libvirt-domain.h:1141
const (
	DomainRebootDefault      DomainRebootFlagValues = iota
	DomainRebootAcpiPowerBtn DomainRebootFlagValues = 1
	DomainRebootGuestAgent   DomainRebootFlagValues = 2
	DomainRebootInitctl      DomainRebootFlagValues = 4
	DomainRebootSignal       DomainRebootFlagValues = 8
	DomainRebootParavirt     DomainRebootFlagValues = 16
)

// DomainDestroyFlagsValues as declared in libvirt/libvirt-domain.h:1159
type DomainDestroyFlagsValues int32

// DomainDestroyFlagsValues enumeration from libvirt/libvirt-domain.h:1159
const (
	DomainDestroyDefault  DomainDestroyFlagsValues = iota
	DomainDestroyGraceful DomainDestroyFlagsValues = 1
)

// DomainSaveRestoreFlags as declared in libvirt/libvirt-domain.h:1191
type DomainSaveRestoreFlags int32

// DomainSaveRestoreFlags enumeration from libvirt/libvirt-domain.h:1191
const (
	DomainSaveBypassCache DomainSaveRestoreFlags = 1
	DomainSaveRunning     DomainSaveRestoreFlags = 2
	DomainSavePaused      DomainSaveRestoreFlags = 4
)

// DomainMemoryModFlags as declared in libvirt/libvirt-domain.h:1446
type DomainMemoryModFlags int32

// DomainMemoryModFlags enumeration from libvirt/libvirt-domain.h:1446
const (
	DomainMemCurrent DomainMemoryModFlags = iota
	DomainMemLive    DomainMemoryModFlags = 1
	DomainMemConfig  DomainMemoryModFlags = 2
	DomainMemMaximum DomainMemoryModFlags = 4
)

// DomainNumatuneMemMode as declared in libvirt/libvirt-domain.h:1464
type DomainNumatuneMemMode int32

// DomainNumatuneMemMode enumeration from libvirt/libvirt-domain.h:1464
const (
	DomainNumatuneMemStrict     DomainNumatuneMemMode = iota
	DomainNumatuneMemPreferred  DomainNumatuneMemMode = 1
	DomainNumatuneMemInterleave DomainNumatuneMemMode = 2
)

// DomainMetadataType as declared in libvirt/libvirt-domain.h:1526
type DomainMetadataType int32

// DomainMetadataType enumeration from libvirt/libvirt-domain.h:1526
const (
	DomainMetadataDescription DomainMetadataType = iota
	DomainMetadataTitle       DomainMetadataType = 1
	DomainMetadataElement     DomainMetadataType = 2
)

// DomainXMLFlags as declared in libvirt/libvirt-domain.h:1556
type DomainXMLFlags int32

// DomainXMLFlags enumeration from libvirt/libvirt-domain.h:1556
const (
	DomainXMLSecure     DomainXMLFlags = 1
	DomainXMLInactive   DomainXMLFlags = 2
	DomainXMLUpdateCPU  DomainXMLFlags = 4
	DomainXMLMigratable DomainXMLFlags = 8
)

// DomainBlockResizeFlags as declared in libvirt/libvirt-domain.h:1661
type DomainBlockResizeFlags int32

// DomainBlockResizeFlags enumeration from libvirt/libvirt-domain.h:1661
const (
	DomainBlockResizeBytes DomainBlockResizeFlags = 1
)

// DomainMemoryFlags as declared in libvirt/libvirt-domain.h:1724
type DomainMemoryFlags int32

// DomainMemoryFlags enumeration from libvirt/libvirt-domain.h:1724
const (
	MemoryVirtual  DomainMemoryFlags = 1
	MemoryPhysical DomainMemoryFlags = 2
)

// DomainDefineFlags as declared in libvirt/libvirt-domain.h:1734
type DomainDefineFlags int32

// DomainDefineFlags enumeration from libvirt/libvirt-domain.h:1734
const (
	DomainDefineValidate DomainDefineFlags = 1
)

// DomainUndefineFlagsValues as declared in libvirt/libvirt-domain.h:1758
type DomainUndefineFlagsValues int32

// DomainUndefineFlagsValues enumeration from libvirt/libvirt-domain.h:1758
const (
	DomainUndefineManagedSave       DomainUndefineFlagsValues = 1
	DomainUndefineSnapshotsMetadata DomainUndefineFlagsValues = 2
	DomainUndefineNvram             DomainUndefineFlagsValues = 4
	DomainUndefineKeepNvram         DomainUndefineFlagsValues = 8
)

// ConnectListAllDomainsFlags as declared in libvirt/libvirt-domain.h:1794
type ConnectListAllDomainsFlags int32

// ConnectListAllDomainsFlags enumeration from libvirt/libvirt-domain.h:1794
const (
	ConnectListDomainsActive        ConnectListAllDomainsFlags = 1
	ConnectListDomainsInactive      ConnectListAllDomainsFlags = 2
	ConnectListDomainsPersistent    ConnectListAllDomainsFlags = 4
	ConnectListDomainsTransient     ConnectListAllDomainsFlags = 8
	ConnectListDomainsRunning       ConnectListAllDomainsFlags = 16
	ConnectListDomainsPaused        ConnectListAllDomainsFlags = 32
	ConnectListDomainsShutoff       ConnectListAllDomainsFlags = 64
	ConnectListDomainsOther         ConnectListAllDomainsFlags = 128
	ConnectListDomainsManagedsave   ConnectListAllDomainsFlags = 256
	ConnectListDomainsNoManagedsave ConnectListAllDomainsFlags = 512
	ConnectListDomainsAutostart     ConnectListAllDomainsFlags = 1024
	ConnectListDomainsNoAutostart   ConnectListAllDomainsFlags = 2048
	ConnectListDomainsHasSnapshot   ConnectListAllDomainsFlags = 4096
	ConnectListDomainsNoSnapshot    ConnectListAllDomainsFlags = 8192
)

// VCPUState as declared in libvirt/libvirt-domain.h:1825
type VCPUState int32

// VCPUState enumeration from libvirt/libvirt-domain.h:1825
const (
	VCPUOffline VCPUState = iota
	VCPURunning VCPUState = 1
	VCPUBlocked VCPUState = 2
)

// DomainVCPUFlags as declared in libvirt/libvirt-domain.h:1847
type DomainVCPUFlags int32

// DomainVCPUFlags enumeration from libvirt/libvirt-domain.h:1847
const (
	DomainVCPUCurrent      DomainVCPUFlags = iota
	DomainVCPULive         DomainVCPUFlags = 1
	DomainVCPUConfig       DomainVCPUFlags = 2
	DomainVCPUMaximum      DomainVCPUFlags = 4
	DomainVCPUGuest        DomainVCPUFlags = 8
	DomainVCPUHotpluggable DomainVCPUFlags = 16
)

// DomainDeviceModifyFlags as declared in libvirt/libvirt-domain.h:2020
type DomainDeviceModifyFlags int32

// DomainDeviceModifyFlags enumeration from libvirt/libvirt-domain.h:2020
const (
	DomainDeviceModifyCurrent DomainDeviceModifyFlags = iota
	DomainDeviceModifyLive    DomainDeviceModifyFlags = 1
	DomainDeviceModifyConfig  DomainDeviceModifyFlags = 2
	DomainDeviceModifyForce   DomainDeviceModifyFlags = 4
)

// DomainStatsTypes as declared in libvirt/libvirt-domain.h:2051
type DomainStatsTypes int32

// DomainStatsTypes enumeration from libvirt/libvirt-domain.h:2051
const (
	DomainStatsState     DomainStatsTypes = 1
	DomainStatsCPUTotal  DomainStatsTypes = 2
	DomainStatsBalloon   DomainStatsTypes = 4
	DomainStatsVCPU      DomainStatsTypes = 8
	DomainStatsInterface DomainStatsTypes = 16
	DomainStatsBlock     DomainStatsTypes = 32
	DomainStatsPerf      DomainStatsTypes = 64
)

// ConnectGetAllDomainStatsFlags as declared in libvirt/libvirt-domain.h:2069
type ConnectGetAllDomainStatsFlags int32

// ConnectGetAllDomainStatsFlags enumeration from libvirt/libvirt-domain.h:2069
const (
	ConnectGetAllDomainsStatsActive       ConnectGetAllDomainStatsFlags = 1
	ConnectGetAllDomainsStatsInactive     ConnectGetAllDomainStatsFlags = 2
	ConnectGetAllDomainsStatsPersistent   ConnectGetAllDomainStatsFlags = 4
	ConnectGetAllDomainsStatsTransient    ConnectGetAllDomainStatsFlags = 8
	ConnectGetAllDomainsStatsRunning      ConnectGetAllDomainStatsFlags = 16
	ConnectGetAllDomainsStatsPaused       ConnectGetAllDomainStatsFlags = 32
	ConnectGetAllDomainsStatsShutoff      ConnectGetAllDomainStatsFlags = 64
	ConnectGetAllDomainsStatsOther        ConnectGetAllDomainStatsFlags = 128
	ConnectGetAllDomainsStatsNowait       ConnectGetAllDomainStatsFlags = 536870912
	ConnectGetAllDomainsStatsBacking      ConnectGetAllDomainStatsFlags = 1073741824
	ConnectGetAllDomainsStatsEnforceStats ConnectGetAllDomainStatsFlags = -2147483648
)

// DomainBlockJobType as declared in libvirt/libvirt-domain.h:2353
type DomainBlockJobType int32

// DomainBlockJobType enumeration from libvirt/libvirt-domain.h:2353
const (
	DomainBlockJobTypeUnknown      DomainBlockJobType = iota
	DomainBlockJobTypePull         DomainBlockJobType = 1
	DomainBlockJobTypeCopy         DomainBlockJobType = 2
	DomainBlockJobTypeCommit       DomainBlockJobType = 3
	DomainBlockJobTypeActiveCommit DomainBlockJobType = 4
)

// DomainBlockJobAbortFlags as declared in libvirt/libvirt-domain.h:2365
type DomainBlockJobAbortFlags int32

// DomainBlockJobAbortFlags enumeration from libvirt/libvirt-domain.h:2365
const (
	DomainBlockJobAbortAsync DomainBlockJobAbortFlags = 1
	DomainBlockJobAbortPivot DomainBlockJobAbortFlags = 2
)

// DomainBlockJobInfoFlags as declared in libvirt/libvirt-domain.h:2374
type DomainBlockJobInfoFlags int32

// DomainBlockJobInfoFlags enumeration from libvirt/libvirt-domain.h:2374
const (
	DomainBlockJobInfoBandwidthBytes DomainBlockJobInfoFlags = 1
)

// DomainBlockJobSetSpeedFlags as declared in libvirt/libvirt-domain.h:2403
type DomainBlockJobSetSpeedFlags int32

// DomainBlockJobSetSpeedFlags enumeration from libvirt/libvirt-domain.h:2403
const (
	DomainBlockJobSpeedBandwidthBytes DomainBlockJobSetSpeedFlags = 1
)

// DomainBlockPullFlags as declared in libvirt/libvirt-domain.h:2413
type DomainBlockPullFlags int32

// DomainBlockPullFlags enumeration from libvirt/libvirt-domain.h:2413
const (
	DomainBlockPullBandwidthBytes DomainBlockPullFlags = 64
)

// DomainBlockRebaseFlags as declared in libvirt/libvirt-domain.h:2437
type DomainBlockRebaseFlags int32

// DomainBlockRebaseFlags enumeration from libvirt/libvirt-domain.h:2437
const (
	DomainBlockRebaseShallow        DomainBlockRebaseFlags = 1
	DomainBlockRebaseReuseExt       DomainBlockRebaseFlags = 2
	DomainBlockRebaseCopyRaw        DomainBlockRebaseFlags = 4
	DomainBlockRebaseCopy           DomainBlockRebaseFlags = 8
	DomainBlockRebaseRelative       DomainBlockRebaseFlags = 16
	DomainBlockRebaseCopyDev        DomainBlockRebaseFlags = 32
	DomainBlockRebaseBandwidthBytes DomainBlockRebaseFlags = 64
)

// DomainBlockCopyFlags as declared in libvirt/libvirt-domain.h:2456
type DomainBlockCopyFlags int32

// DomainBlockCopyFlags enumeration from libvirt/libvirt-domain.h:2456
const (
	DomainBlockCopyShallow      DomainBlockCopyFlags = 1
	DomainBlockCopyReuseExt     DomainBlockCopyFlags = 2
	DomainBlockCopyTransientJob DomainBlockCopyFlags = 4
)

// DomainBlockCommitFlags as declared in libvirt/libvirt-domain.h:2521
type DomainBlockCommitFlags int32

// DomainBlockCommitFlags enumeration from libvirt/libvirt-domain.h:2521
const (
	DomainBlockCommitShallow        DomainBlockCommitFlags = 1
	DomainBlockCommitDelete         DomainBlockCommitFlags = 2
	DomainBlockCommitActive         DomainBlockCommitFlags = 4
	DomainBlockCommitRelative       DomainBlockCommitFlags = 8
	DomainBlockCommitBandwidthBytes DomainBlockCommitFlags = 16
)

// DomainDiskErrorCode as declared in libvirt/libvirt-domain.h:2712
type DomainDiskErrorCode int32

// DomainDiskErrorCode enumeration from libvirt/libvirt-domain.h:2712
const (
	DomainDiskErrorNone    DomainDiskErrorCode = iota
	DomainDiskErrorUnspec  DomainDiskErrorCode = 1
	DomainDiskErrorNoSpace DomainDiskErrorCode = 2
)

// KeycodeSet as declared in libvirt/libvirt-domain.h:2758
type KeycodeSet int32

// KeycodeSet enumeration from libvirt/libvirt-domain.h:2758
const (
	KeycodeSetLinux  KeycodeSet = iota
	KeycodeSetXt     KeycodeSet = 1
	KeycodeSetAtset1 KeycodeSet = 2
	KeycodeSetAtset2 KeycodeSet = 3
	KeycodeSetAtset3 KeycodeSet = 4
	KeycodeSetOsx    KeycodeSet = 5
	KeycodeSetXtKbd  KeycodeSet = 6
	KeycodeSetUsb    KeycodeSet = 7
	KeycodeSetWin32  KeycodeSet = 8
	KeycodeSetQnum   KeycodeSet = 9
)

// DomainProcessSignal as declared in libvirt/libvirt-domain.h:2867
type DomainProcessSignal int32

// DomainProcessSignal enumeration from libvirt/libvirt-domain.h:2867
const (
	DomainProcessSignalNop    DomainProcessSignal = iota
	DomainProcessSignalHup    DomainProcessSignal = 1
	DomainProcessSignalInt    DomainProcessSignal = 2
	DomainProcessSignalQuit   DomainProcessSignal = 3
	DomainProcessSignalIll    DomainProcessSignal = 4
	DomainProcessSignalTrap   DomainProcessSignal = 5
	DomainProcessSignalAbrt   DomainProcessSignal = 6
	DomainProcessSignalBus    DomainProcessSignal = 7
	DomainProcessSignalFpe    DomainProcessSignal = 8
	DomainProcessSignalKill   DomainProcessSignal = 9
	DomainProcessSignalUsr1   DomainProcessSignal = 10
	DomainProcessSignalSegv   DomainProcessSignal = 11
	DomainProcessSignalUsr2   DomainProcessSignal = 12
	DomainProcessSignalPipe   DomainProcessSignal = 13
	DomainProcessSignalAlrm   DomainProcessSignal = 14
	DomainProcessSignalTerm   DomainProcessSignal = 15
	DomainProcessSignalStkflt DomainProcessSignal = 16
	DomainProcessSignalChld   DomainProcessSignal = 17
	DomainProcessSignalCont   DomainProcessSignal = 18
	DomainProcessSignalStop   DomainProcessSignal = 19
	DomainProcessSignalTstp   DomainProcessSignal = 20
	DomainProcessSignalTtin   DomainProcessSignal = 21
	DomainProcessSignalTtou   DomainProcessSignal = 22
	DomainProcessSignalUrg    DomainProcessSignal = 23
	DomainProcessSignalXcpu   DomainProcessSignal = 24
	DomainProcessSignalXfsz   DomainProcessSignal = 25
	DomainProcessSignalVtalrm DomainProcessSignal = 26
	DomainProcessSignalProf   DomainProcessSignal = 27
	DomainProcessSignalWinch  DomainProcessSignal = 28
	DomainProcessSignalPoll   DomainProcessSignal = 29
	DomainProcessSignalPwr    DomainProcessSignal = 30
	DomainProcessSignalSys    DomainProcessSignal = 31
	DomainProcessSignalRt0    DomainProcessSignal = 32
	DomainProcessSignalRt1    DomainProcessSignal = 33
	DomainProcessSignalRt2    DomainProcessSignal = 34
	DomainProcessSignalRt3    DomainProcessSignal = 35
	DomainProcessSignalRt4    DomainProcessSignal = 36
	DomainProcessSignalRt5    DomainProcessSignal = 37
	DomainProcessSignalRt6    DomainProcessSignal = 38
	DomainProcessSignalRt7    DomainProcessSignal = 39
	DomainProcessSignalRt8    DomainProcessSignal = 40
	DomainProcessSignalRt9    DomainProcessSignal = 41
	DomainProcessSignalRt10   DomainProcessSignal = 42
	DomainProcessSignalRt11   DomainProcessSignal = 43
	DomainProcessSignalRt12   DomainProcessSignal = 44
	DomainProcessSignalRt13   DomainProcessSignal = 45
	DomainProcessSignalRt14   DomainProcessSignal = 46
	DomainProcessSignalRt15   DomainProcessSignal = 47
	DomainProcessSignalRt16   DomainProcessSignal = 48
	DomainProcessSignalRt17   DomainProcessSignal = 49
	DomainProcessSignalRt18   DomainProcessSignal = 50
	DomainProcessSignalRt19   DomainProcessSignal = 51
	DomainProcessSignalRt20   DomainProcessSignal = 52
	DomainProcessSignalRt21   DomainProcessSignal = 53
	DomainProcessSignalRt22   DomainProcessSignal = 54
	DomainProcessSignalRt23   DomainProcessSignal = 55
	DomainProcessSignalRt24   DomainProcessSignal = 56
	DomainProcessSignalRt25   DomainProcessSignal = 57
	DomainProcessSignalRt26   DomainProcessSignal = 58
	DomainProcessSignalRt27   DomainProcessSignal = 59
	DomainProcessSignalRt28   DomainProcessSignal = 60
	DomainProcessSignalRt29   DomainProcessSignal = 61
	DomainProcessSignalRt30   DomainProcessSignal = 62
	DomainProcessSignalRt31   DomainProcessSignal = 63
	DomainProcessSignalRt32   DomainProcessSignal = 64
)

// DomainEventType as declared in libvirt/libvirt-domain.h:2905
type DomainEventType int32

// DomainEventType enumeration from libvirt/libvirt-domain.h:2905
const (
	DomainEventDefined     DomainEventType = iota
	DomainEventUndefined   DomainEventType = 1
	DomainEventStarted     DomainEventType = 2
	DomainEventSuspended   DomainEventType = 3
	DomainEventResumed     DomainEventType = 4
	DomainEventStopped     DomainEventType = 5
	DomainEventShutdown    DomainEventType = 6
	DomainEventPmsuspended DomainEventType = 7
	DomainEventCrashed     DomainEventType = 8
)

// DomainEventDefinedDetailType as declared in libvirt/libvirt-domain.h:2921
type DomainEventDefinedDetailType int32

// DomainEventDefinedDetailType enumeration from libvirt/libvirt-domain.h:2921
const (
	DomainEventDefinedAdded        DomainEventDefinedDetailType = iota
	DomainEventDefinedUpdated      DomainEventDefinedDetailType = 1
	DomainEventDefinedRenamed      DomainEventDefinedDetailType = 2
	DomainEventDefinedFromSnapshot DomainEventDefinedDetailType = 3
)

// DomainEventUndefinedDetailType as declared in libvirt/libvirt-domain.h:2935
type DomainEventUndefinedDetailType int32

// DomainEventUndefinedDetailType enumeration from libvirt/libvirt-domain.h:2935
const (
	DomainEventUndefinedRemoved DomainEventUndefinedDetailType = iota
	DomainEventUndefinedRenamed DomainEventUndefinedDetailType = 1
)

// DomainEventStartedDetailType as declared in libvirt/libvirt-domain.h:2952
type DomainEventStartedDetailType int32

// DomainEventStartedDetailType enumeration from libvirt/libvirt-domain.h:2952
const (
	DomainEventStartedBooted       DomainEventStartedDetailType = iota
	DomainEventStartedMigrated     DomainEventStartedDetailType = 1
	DomainEventStartedRestored     DomainEventStartedDetailType = 2
	DomainEventStartedFromSnapshot DomainEventStartedDetailType = 3
	DomainEventStartedWakeup       DomainEventStartedDetailType = 4
)

// DomainEventSuspendedDetailType as declared in libvirt/libvirt-domain.h:2973
type DomainEventSuspendedDetailType int32

// DomainEventSuspendedDetailType enumeration from libvirt/libvirt-domain.h:2973
const (
	DomainEventSuspendedPaused         DomainEventSuspendedDetailType = iota
	DomainEventSuspendedMigrated       DomainEventSuspendedDetailType = 1
	DomainEventSuspendedIoerror        DomainEventSuspendedDetailType = 2
	DomainEventSuspendedWatchdog       DomainEventSuspendedDetailType = 3
	DomainEventSuspendedRestored       DomainEventSuspendedDetailType = 4
	DomainEventSuspendedFromSnapshot   DomainEventSuspendedDetailType = 5
	DomainEventSuspendedAPIError       DomainEventSuspendedDetailType = 6
	DomainEventSuspendedPostcopy       DomainEventSuspendedDetailType = 7
	DomainEventSuspendedPostcopyFailed DomainEventSuspendedDetailType = 8
)

// DomainEventResumedDetailType as declared in libvirt/libvirt-domain.h:2990
type DomainEventResumedDetailType int32

// DomainEventResumedDetailType enumeration from libvirt/libvirt-domain.h:2990
const (
	DomainEventResumedUnpaused     DomainEventResumedDetailType = iota
	DomainEventResumedMigrated     DomainEventResumedDetailType = 1
	DomainEventResumedFromSnapshot DomainEventResumedDetailType = 2
	DomainEventResumedPostcopy     DomainEventResumedDetailType = 3
)

// DomainEventStoppedDetailType as declared in libvirt/libvirt-domain.h:3009
type DomainEventStoppedDetailType int32

// DomainEventStoppedDetailType enumeration from libvirt/libvirt-domain.h:3009
const (
	DomainEventStoppedShutdown     DomainEventStoppedDetailType = iota
	DomainEventStoppedDestroyed    DomainEventStoppedDetailType = 1
	DomainEventStoppedCrashed      DomainEventStoppedDetailType = 2
	DomainEventStoppedMigrated     DomainEventStoppedDetailType = 3
	DomainEventStoppedSaved        DomainEventStoppedDetailType = 4
	DomainEventStoppedFailed       DomainEventStoppedDetailType = 5
	DomainEventStoppedFromSnapshot DomainEventStoppedDetailType = 6
)

// DomainEventShutdownDetailType as declared in libvirt/libvirt-domain.h:3032
type DomainEventShutdownDetailType int32

// DomainEventShutdownDetailType enumeration from libvirt/libvirt-domain.h:3032
const (
	DomainEventShutdownFinished DomainEventShutdownDetailType = iota
	DomainEventShutdownGuest    DomainEventShutdownDetailType = 1
	DomainEventShutdownHost     DomainEventShutdownDetailType = 2
)

// DomainEventPMSuspendedDetailType as declared in libvirt/libvirt-domain.h:3046
type DomainEventPMSuspendedDetailType int32

// DomainEventPMSuspendedDetailType enumeration from libvirt/libvirt-domain.h:3046
const (
	DomainEventPmsuspendedMemory DomainEventPMSuspendedDetailType = iota
	DomainEventPmsuspendedDisk   DomainEventPMSuspendedDetailType = 1
)

// DomainEventCrashedDetailType as declared in libvirt/libvirt-domain.h:3059
type DomainEventCrashedDetailType int32

// DomainEventCrashedDetailType enumeration from libvirt/libvirt-domain.h:3059
const (
	DomainEventCrashedPanicked DomainEventCrashedDetailType = iota
)

// DomainJobType as declared in libvirt/libvirt-domain.h:3103
type DomainJobType int32

// DomainJobType enumeration from libvirt/libvirt-domain.h:3103
const (
	DomainJobNone      DomainJobType = iota
	DomainJobBounded   DomainJobType = 1
	DomainJobUnbounded DomainJobType = 2
	DomainJobCompleted DomainJobType = 3
	DomainJobFailed    DomainJobType = 4
	DomainJobCancelled DomainJobType = 5
)

// DomainGetJobStatsFlags as declared in libvirt/libvirt-domain.h:3150
type DomainGetJobStatsFlags int32

// DomainGetJobStatsFlags enumeration from libvirt/libvirt-domain.h:3150
const (
	DomainJobStatsCompleted DomainGetJobStatsFlags = 1
)

// DomainJobOperation as declared in libvirt/libvirt-domain.h:3175
type DomainJobOperation int32

// DomainJobOperation enumeration from libvirt/libvirt-domain.h:3175
const (
	DomainJobOperationStrUnknown        DomainJobOperation = iota
	DomainJobOperationStrStart          DomainJobOperation = 1
	DomainJobOperationStrSave           DomainJobOperation = 2
	DomainJobOperationStrRestore        DomainJobOperation = 3
	DomainJobOperationStrMigrationIn    DomainJobOperation = 4
	DomainJobOperationStrMigrationOut   DomainJobOperation = 5
	DomainJobOperationStrSnapshot       DomainJobOperation = 6
	DomainJobOperationStrSnapshotRevert DomainJobOperation = 7
	DomainJobOperationStrDump           DomainJobOperation = 8
)

// DomainEventWatchdogAction as declared in libvirt/libvirt-domain.h:3519
type DomainEventWatchdogAction int32

// DomainEventWatchdogAction enumeration from libvirt/libvirt-domain.h:3519
const (
	DomainEventWatchdogNone      DomainEventWatchdogAction = iota
	DomainEventWatchdogPause     DomainEventWatchdogAction = 1
	DomainEventWatchdogReset     DomainEventWatchdogAction = 2
	DomainEventWatchdogPoweroff  DomainEventWatchdogAction = 3
	DomainEventWatchdogShutdown  DomainEventWatchdogAction = 4
	DomainEventWatchdogDebug     DomainEventWatchdogAction = 5
	DomainEventWatchdogInjectnmi DomainEventWatchdogAction = 6
)

// DomainEventIOErrorAction as declared in libvirt/libvirt-domain.h:3550
type DomainEventIOErrorAction int32

// DomainEventIOErrorAction enumeration from libvirt/libvirt-domain.h:3550
const (
	DomainEventIoErrorNone   DomainEventIOErrorAction = iota
	DomainEventIoErrorPause  DomainEventIOErrorAction = 1
	DomainEventIoErrorReport DomainEventIOErrorAction = 2
)

// DomainEventGraphicsPhase as declared in libvirt/libvirt-domain.h:3613
type DomainEventGraphicsPhase int32

// DomainEventGraphicsPhase enumeration from libvirt/libvirt-domain.h:3613
const (
	DomainEventGraphicsConnect    DomainEventGraphicsPhase = iota
	DomainEventGraphicsInitialize DomainEventGraphicsPhase = 1
	DomainEventGraphicsDisconnect DomainEventGraphicsPhase = 2
)

// DomainEventGraphicsAddressType as declared in libvirt/libvirt-domain.h:3628
type DomainEventGraphicsAddressType int32

// DomainEventGraphicsAddressType enumeration from libvirt/libvirt-domain.h:3628
const (
	DomainEventGraphicsAddressIpv4 DomainEventGraphicsAddressType = iota
	DomainEventGraphicsAddressIpv6 DomainEventGraphicsAddressType = 1
	DomainEventGraphicsAddressUnix DomainEventGraphicsAddressType = 2
)

// ConnectDomainEventBlockJobStatus as declared in libvirt/libvirt-domain.h:3716
type ConnectDomainEventBlockJobStatus int32

// ConnectDomainEventBlockJobStatus enumeration from libvirt/libvirt-domain.h:3716
const (
	DomainBlockJobCompleted ConnectDomainEventBlockJobStatus = iota
	DomainBlockJobFailed    ConnectDomainEventBlockJobStatus = 1
	DomainBlockJobCanceled  ConnectDomainEventBlockJobStatus = 2
	DomainBlockJobReady     ConnectDomainEventBlockJobStatus = 3
)

// ConnectDomainEventDiskChangeReason as declared in libvirt/libvirt-domain.h:3766
type ConnectDomainEventDiskChangeReason int32

// ConnectDomainEventDiskChangeReason enumeration from libvirt/libvirt-domain.h:3766
const (
	DomainEventDiskChangeMissingOnStart ConnectDomainEventDiskChangeReason = iota
	DomainEventDiskDropMissingOnStart   ConnectDomainEventDiskChangeReason = 1
)

// DomainEventTrayChangeReason as declared in libvirt/libvirt-domain.h:3807
type DomainEventTrayChangeReason int32

// DomainEventTrayChangeReason enumeration from libvirt/libvirt-domain.h:3807
const (
	DomainEventTrayChangeOpen  DomainEventTrayChangeReason = iota
	DomainEventTrayChangeClose DomainEventTrayChangeReason = 1
)

// ConnectDomainEventAgentLifecycleState as declared in libvirt/libvirt-domain.h:4322
type ConnectDomainEventAgentLifecycleState int32

// ConnectDomainEventAgentLifecycleState enumeration from libvirt/libvirt-domain.h:4322
const (
	ConnectDomainEventAgentLifecycleStateConnected    ConnectDomainEventAgentLifecycleState = 1
	ConnectDomainEventAgentLifecycleStateDisconnected ConnectDomainEventAgentLifecycleState = 2
)

// ConnectDomainEventAgentLifecycleReason as declared in libvirt/libvirt-domain.h:4332
type ConnectDomainEventAgentLifecycleReason int32

// ConnectDomainEventAgentLifecycleReason enumeration from libvirt/libvirt-domain.h:4332
const (
	ConnectDomainEventAgentLifecycleReasonUnknown       ConnectDomainEventAgentLifecycleReason = iota
	ConnectDomainEventAgentLifecycleReasonDomainStarted ConnectDomainEventAgentLifecycleReason = 1
	ConnectDomainEventAgentLifecycleReasonChannel       ConnectDomainEventAgentLifecycleReason = 2
)

// DomainEventID as declared in libvirt/libvirt-domain.h:4436
type DomainEventID int32

// DomainEventID enumeration from libvirt/libvirt-domain.h:4436
const (
	DomainEventIDLifecycle           DomainEventID = iota
	DomainEventIDReboot              DomainEventID = 1
	DomainEventIDRtcChange           DomainEventID = 2
	DomainEventIDWatchdog            DomainEventID = 3
	DomainEventIDIoError             DomainEventID = 4
	DomainEventIDGraphics            DomainEventID = 5
	DomainEventIDIoErrorReason       DomainEventID = 6
	DomainEventIDControlError        DomainEventID = 7
	DomainEventIDBlockJob            DomainEventID = 8
	DomainEventIDDiskChange          DomainEventID = 9
	DomainEventIDTrayChange          DomainEventID = 10
	DomainEventIDPmwakeup            DomainEventID = 11
	DomainEventIDPmsuspend           DomainEventID = 12
	DomainEventIDBalloonChange       DomainEventID = 13
	DomainEventIDPmsuspendDisk       DomainEventID = 14
	DomainEventIDDeviceRemoved       DomainEventID = 15
	DomainEventIDBlockJob2           DomainEventID = 16
	DomainEventIDTunable             DomainEventID = 17
	DomainEventIDAgentLifecycle      DomainEventID = 18
	DomainEventIDDeviceAdded         DomainEventID = 19
	DomainEventIDMigrationIteration  DomainEventID = 20
	DomainEventIDJobCompleted        DomainEventID = 21
	DomainEventIDDeviceRemovalFailed DomainEventID = 22
	DomainEventIDMetadataChange      DomainEventID = 23
	DomainEventIDBlockThreshold      DomainEventID = 24
)

// DomainConsoleFlags as declared in libvirt/libvirt-domain.h:4463
type DomainConsoleFlags int32

// DomainConsoleFlags enumeration from libvirt/libvirt-domain.h:4463
const (
	DomainConsoleForce DomainConsoleFlags = 1
	DomainConsoleSafe  DomainConsoleFlags = 2
)

// DomainChannelFlags as declared in libvirt/libvirt-domain.h:4479
type DomainChannelFlags int32

// DomainChannelFlags enumeration from libvirt/libvirt-domain.h:4479
const (
	DomainChannelForce DomainChannelFlags = 1
)

// DomainOpenGraphicsFlags as declared in libvirt/libvirt-domain.h:4488
type DomainOpenGraphicsFlags int32

// DomainOpenGraphicsFlags enumeration from libvirt/libvirt-domain.h:4488
const (
	DomainOpenGraphicsSkipauth DomainOpenGraphicsFlags = 1
)

// DomainSetTimeFlags as declared in libvirt/libvirt-domain.h:4545
type DomainSetTimeFlags int32

// DomainSetTimeFlags enumeration from libvirt/libvirt-domain.h:4545
const (
	DomainTimeSync DomainSetTimeFlags = 1
)

// SchedParameterType as declared in libvirt/libvirt-domain.h:4566
type SchedParameterType int32

// SchedParameterType enumeration from libvirt/libvirt-domain.h:4566
const (
	DomainSchedFieldInt     SchedParameterType = 1
	DomainSchedFieldUint    SchedParameterType = 2
	DomainSchedFieldLlong   SchedParameterType = 3
	DomainSchedFieldUllong  SchedParameterType = 4
	DomainSchedFieldDouble  SchedParameterType = 5
	DomainSchedFieldBoolean SchedParameterType = 6
)

// BlkioParameterType as declared in libvirt/libvirt-domain.h:4610
type BlkioParameterType int32

// BlkioParameterType enumeration from libvirt/libvirt-domain.h:4610
const (
	DomainBlkioParamInt     BlkioParameterType = 1
	DomainBlkioParamUint    BlkioParameterType = 2
	DomainBlkioParamLlong   BlkioParameterType = 3
	DomainBlkioParamUllong  BlkioParameterType = 4
	DomainBlkioParamDouble  BlkioParameterType = 5
	DomainBlkioParamBoolean BlkioParameterType = 6
)

// MemoryParameterType as declared in libvirt/libvirt-domain.h:4654
type MemoryParameterType int32

// MemoryParameterType enumeration from libvirt/libvirt-domain.h:4654
const (
	DomainMemoryParamInt     MemoryParameterType = 1
	DomainMemoryParamUint    MemoryParameterType = 2
	DomainMemoryParamLlong   MemoryParameterType = 3
	DomainMemoryParamUllong  MemoryParameterType = 4
	DomainMemoryParamDouble  MemoryParameterType = 5
	DomainMemoryParamBoolean MemoryParameterType = 6
)

// DomainInterfaceAddressesSource as declared in libvirt/libvirt-domain.h:4692
type DomainInterfaceAddressesSource int32

// DomainInterfaceAddressesSource enumeration from libvirt/libvirt-domain.h:4692
const (
	DomainInterfaceAddressesSrcLease DomainInterfaceAddressesSource = iota
	DomainInterfaceAddressesSrcAgent DomainInterfaceAddressesSource = 1
	DomainInterfaceAddressesSrcArp   DomainInterfaceAddressesSource = 2
)

// DomainSetUserPasswordFlags as declared in libvirt/libvirt-domain.h:4720
type DomainSetUserPasswordFlags int32

// DomainSetUserPasswordFlags enumeration from libvirt/libvirt-domain.h:4720
const (
	DomainPasswordEncrypted DomainSetUserPasswordFlags = 1
)

// DomainLifecycle as declared in libvirt/libvirt-domain.h:4759
type DomainLifecycle int32

// DomainLifecycle enumeration from libvirt/libvirt-domain.h:4759
const (
	DomainLifecyclePoweroff DomainLifecycle = iota
	DomainLifecycleReboot   DomainLifecycle = 1
	DomainLifecycleCrash    DomainLifecycle = 2
)

// DomainLifecycleAction as declared in libvirt/libvirt-domain.h:4772
type DomainLifecycleAction int32

// DomainLifecycleAction enumeration from libvirt/libvirt-domain.h:4772
const (
	DomainLifecycleActionDestroy         DomainLifecycleAction = iota
	DomainLifecycleActionRestart         DomainLifecycleAction = 1
	DomainLifecycleActionRestartRename   DomainLifecycleAction = 2
	DomainLifecycleActionPreserve        DomainLifecycleAction = 3
	DomainLifecycleActionCoredumpDestroy DomainLifecycleAction = 4
	DomainLifecycleActionCoredumpRestart DomainLifecycleAction = 5
)

// DomainSnapshotCreateFlags as declared in libvirt/libvirt-domain-snapshot.h:75
type DomainSnapshotCreateFlags int32

// DomainSnapshotCreateFlags enumeration from libvirt/libvirt-domain-snapshot.h:75
const (
	DomainSnapshotCreateRedefine   DomainSnapshotCreateFlags = 1
	DomainSnapshotCreateCurrent    DomainSnapshotCreateFlags = 2
	DomainSnapshotCreateNoMetadata DomainSnapshotCreateFlags = 4
	DomainSnapshotCreateHalt       DomainSnapshotCreateFlags = 8
	DomainSnapshotCreateDiskOnly   DomainSnapshotCreateFlags = 16
	DomainSnapshotCreateReuseExt   DomainSnapshotCreateFlags = 32
	DomainSnapshotCreateQuiesce    DomainSnapshotCreateFlags = 64
	DomainSnapshotCreateAtomic     DomainSnapshotCreateFlags = 128
	DomainSnapshotCreateLive       DomainSnapshotCreateFlags = 256
)

// DomainSnapshotListFlags as declared in libvirt/libvirt-domain-snapshot.h:135
type DomainSnapshotListFlags int32

// DomainSnapshotListFlags enumeration from libvirt/libvirt-domain-snapshot.h:135
const (
	DomainSnapshotListRoots       DomainSnapshotListFlags = 1
	DomainSnapshotListDescendants DomainSnapshotListFlags = 1
	DomainSnapshotListLeaves      DomainSnapshotListFlags = 4
	DomainSnapshotListNoLeaves    DomainSnapshotListFlags = 8
	DomainSnapshotListMetadata    DomainSnapshotListFlags = 2
	DomainSnapshotListNoMetadata  DomainSnapshotListFlags = 16
	DomainSnapshotListInactive    DomainSnapshotListFlags = 32
	DomainSnapshotListActive      DomainSnapshotListFlags = 64
	DomainSnapshotListDiskOnly    DomainSnapshotListFlags = 128
	DomainSnapshotListInternal    DomainSnapshotListFlags = 256
	DomainSnapshotListExternal    DomainSnapshotListFlags = 512
)

// DomainSnapshotRevertFlags as declared in libvirt/libvirt-domain-snapshot.h:192
type DomainSnapshotRevertFlags int32

// DomainSnapshotRevertFlags enumeration from libvirt/libvirt-domain-snapshot.h:192
const (
	DomainSnapshotRevertRunning DomainSnapshotRevertFlags = 1
	DomainSnapshotRevertPaused  DomainSnapshotRevertFlags = 2
	DomainSnapshotRevertForce   DomainSnapshotRevertFlags = 4
)

// DomainSnapshotDeleteFlags as declared in libvirt/libvirt-domain-snapshot.h:206
type DomainSnapshotDeleteFlags int32

// DomainSnapshotDeleteFlags enumeration from libvirt/libvirt-domain-snapshot.h:206
const (
	DomainSnapshotDeleteChildren     DomainSnapshotDeleteFlags = 1
	DomainSnapshotDeleteMetadataOnly DomainSnapshotDeleteFlags = 2
	DomainSnapshotDeleteChildrenOnly DomainSnapshotDeleteFlags = 4
)

// EventHandleType as declared in libvirt/libvirt-event.h:44
type EventHandleType int32

// EventHandleType enumeration from libvirt/libvirt-event.h:44
const (
	EventHandleReadable EventHandleType = 1
	EventHandleWritable EventHandleType = 2
	EventHandleError    EventHandleType = 4
	EventHandleHangup   EventHandleType = 8
)

// ConnectListAllInterfacesFlags as declared in libvirt/libvirt-interface.h:65
type ConnectListAllInterfacesFlags int32

// ConnectListAllInterfacesFlags enumeration from libvirt/libvirt-interface.h:65
const (
	ConnectListInterfacesInactive ConnectListAllInterfacesFlags = 1
	ConnectListInterfacesActive   ConnectListAllInterfacesFlags = 2
)

// InterfaceXMLFlags as declared in libvirt/libvirt-interface.h:81
type InterfaceXMLFlags int32

// InterfaceXMLFlags enumeration from libvirt/libvirt-interface.h:81
const (
	InterfaceXMLInactive InterfaceXMLFlags = 1
)

// NetworkXMLFlags as declared in libvirt/libvirt-network.h:33
type NetworkXMLFlags int32

// NetworkXMLFlags enumeration from libvirt/libvirt-network.h:33
const (
	NetworkXMLInactive NetworkXMLFlags = 1
)

// ConnectListAllNetworksFlags as declared in libvirt/libvirt-network.h:85
type ConnectListAllNetworksFlags int32

// ConnectListAllNetworksFlags enumeration from libvirt/libvirt-network.h:85
const (
	ConnectListNetworksInactive    ConnectListAllNetworksFlags = 1
	ConnectListNetworksActive      ConnectListAllNetworksFlags = 2
	ConnectListNetworksPersistent  ConnectListAllNetworksFlags = 4
	ConnectListNetworksTransient   ConnectListAllNetworksFlags = 8
	ConnectListNetworksAutostart   ConnectListAllNetworksFlags = 16
	ConnectListNetworksNoAutostart ConnectListAllNetworksFlags = 32
)

// NetworkUpdateCommand as declared in libvirt/libvirt-network.h:134
type NetworkUpdateCommand int32

// NetworkUpdateCommand enumeration from libvirt/libvirt-network.h:134
const (
	NetworkUpdateCommandNone     NetworkUpdateCommand = iota
	NetworkUpdateCommandModify   NetworkUpdateCommand = 1
	NetworkUpdateCommandDelete   NetworkUpdateCommand = 2
	NetworkUpdateCommandAddLast  NetworkUpdateCommand = 3
	NetworkUpdateCommandAddFirst NetworkUpdateCommand = 4
)

// NetworkUpdateSection as declared in libvirt/libvirt-network.h:160
type NetworkUpdateSection int32

// NetworkUpdateSection enumeration from libvirt/libvirt-network.h:160
const (
	NetworkSectionNone             NetworkUpdateSection = iota
	NetworkSectionBridge           NetworkUpdateSection = 1
	NetworkSectionDomain           NetworkUpdateSection = 2
	NetworkSectionIP               NetworkUpdateSection = 3
	NetworkSectionIPDhcpHost       NetworkUpdateSection = 4
	NetworkSectionIPDhcpRange      NetworkUpdateSection = 5
	NetworkSectionForward          NetworkUpdateSection = 6
	NetworkSectionForwardInterface NetworkUpdateSection = 7
	NetworkSectionForwardPf        NetworkUpdateSection = 8
	NetworkSectionPortgroup        NetworkUpdateSection = 9
	NetworkSectionDNSHost          NetworkUpdateSection = 10
	NetworkSectionDNSTxt           NetworkUpdateSection = 11
	NetworkSectionDNSSrv           NetworkUpdateSection = 12
)

// NetworkUpdateFlags as declared in libvirt/libvirt-network.h:172
type NetworkUpdateFlags int32

// NetworkUpdateFlags enumeration from libvirt/libvirt-network.h:172
const (
	NetworkUpdateAffectCurrent NetworkUpdateFlags = iota
	NetworkUpdateAffectLive    NetworkUpdateFlags = 1
	NetworkUpdateAffectConfig  NetworkUpdateFlags = 2
)

// NetworkEventLifecycleType as declared in libvirt/libvirt-network.h:230
type NetworkEventLifecycleType int32

// NetworkEventLifecycleType enumeration from libvirt/libvirt-network.h:230
const (
	NetworkEventDefined   NetworkEventLifecycleType = iota
	NetworkEventUndefined NetworkEventLifecycleType = 1
	NetworkEventStarted   NetworkEventLifecycleType = 2
	NetworkEventStopped   NetworkEventLifecycleType = 3
)

// NetworkEventID as declared in libvirt/libvirt-network.h:278
type NetworkEventID int32

// NetworkEventID enumeration from libvirt/libvirt-network.h:278
const (
	NetworkEventIDLifecycle NetworkEventID = iota
)

// IPAddrType as declared in libvirt/libvirt-network.h:287
type IPAddrType int32

// IPAddrType enumeration from libvirt/libvirt-network.h:287
const (
	IPAddrTypeIpv4 IPAddrType = iota
	IPAddrTypeIpv6 IPAddrType = 1
)

// ConnectListAllNodeDeviceFlags as declared in libvirt/libvirt-nodedev.h:85
type ConnectListAllNodeDeviceFlags int32

// ConnectListAllNodeDeviceFlags enumeration from libvirt/libvirt-nodedev.h:85
const (
	ConnectListNodeDevicesCapSystem       ConnectListAllNodeDeviceFlags = 1
	ConnectListNodeDevicesCapPciDev       ConnectListAllNodeDeviceFlags = 2
	ConnectListNodeDevicesCapUsbDev       ConnectListAllNodeDeviceFlags = 4
	ConnectListNodeDevicesCapUsbInterface ConnectListAllNodeDeviceFlags = 8
	ConnectListNodeDevicesCapNet          ConnectListAllNodeDeviceFlags = 16
	ConnectListNodeDevicesCapScsiHost     ConnectListAllNodeDeviceFlags = 32
	ConnectListNodeDevicesCapScsiTarget   ConnectListAllNodeDeviceFlags = 64
	ConnectListNodeDevicesCapScsi         ConnectListAllNodeDeviceFlags = 128
	ConnectListNodeDevicesCapStorage      ConnectListAllNodeDeviceFlags = 256
	ConnectListNodeDevicesCapFcHost       ConnectListAllNodeDeviceFlags = 512
	ConnectListNodeDevicesCapVports       ConnectListAllNodeDeviceFlags = 1024
	ConnectListNodeDevicesCapScsiGeneric  ConnectListAllNodeDeviceFlags = 2048
	ConnectListNodeDevicesCapDrm          ConnectListAllNodeDeviceFlags = 4096
	ConnectListNodeDevicesCapMdevTypes    ConnectListAllNodeDeviceFlags = 8192
	ConnectListNodeDevicesCapMdev         ConnectListAllNodeDeviceFlags = 16384
	ConnectListNodeDevicesCapCcwDev       ConnectListAllNodeDeviceFlags = 32768
)

// NodeDeviceEventID as declared in libvirt/libvirt-nodedev.h:155
type NodeDeviceEventID int32

// NodeDeviceEventID enumeration from libvirt/libvirt-nodedev.h:155
const (
	NodeDeviceEventIDLifecycle NodeDeviceEventID = iota
	NodeDeviceEventIDUpdate    NodeDeviceEventID = 1
)

// NodeDeviceEventLifecycleType as declared in libvirt/libvirt-nodedev.h:197
type NodeDeviceEventLifecycleType int32

// NodeDeviceEventLifecycleType enumeration from libvirt/libvirt-nodedev.h:197
const (
	NodeDeviceEventCreated NodeDeviceEventLifecycleType = iota
	NodeDeviceEventDeleted NodeDeviceEventLifecycleType = 1
)

// SecretUsageType as declared in libvirt/libvirt-secret.h:56
type SecretUsageType int32

// SecretUsageType enumeration from libvirt/libvirt-secret.h:56
const (
	SecretUsageTypeNone   SecretUsageType = iota
	SecretUsageTypeVolume SecretUsageType = 1
	SecretUsageTypeCeph   SecretUsageType = 2
	SecretUsageTypeIscsi  SecretUsageType = 3
	SecretUsageTypeTLS    SecretUsageType = 4
)

// ConnectListAllSecretsFlags as declared in libvirt/libvirt-secret.h:79
type ConnectListAllSecretsFlags int32

// ConnectListAllSecretsFlags enumeration from libvirt/libvirt-secret.h:79
const (
	ConnectListSecretsEphemeral   ConnectListAllSecretsFlags = 1
	ConnectListSecretsNoEphemeral ConnectListAllSecretsFlags = 2
	ConnectListSecretsPrivate     ConnectListAllSecretsFlags = 4
	ConnectListSecretsNoPrivate   ConnectListAllSecretsFlags = 8
)

// SecretEventID as declared in libvirt/libvirt-secret.h:140
type SecretEventID int32

// SecretEventID enumeration from libvirt/libvirt-secret.h:140
const (
	SecretEventIDLifecycle    SecretEventID = iota
	SecretEventIDValueChanged SecretEventID = 1
)

// SecretEventLifecycleType as declared in libvirt/libvirt-secret.h:182
type SecretEventLifecycleType int32

// SecretEventLifecycleType enumeration from libvirt/libvirt-secret.h:182
const (
	SecretEventDefined   SecretEventLifecycleType = iota
	SecretEventUndefined SecretEventLifecycleType = 1
)

// StoragePoolState as declared in libvirt/libvirt-storage.h:58
type StoragePoolState int32

// StoragePoolState enumeration from libvirt/libvirt-storage.h:58
const (
	StoragePoolInactive     StoragePoolState = iota
	StoragePoolBuilding     StoragePoolState = 1
	StoragePoolRunning      StoragePoolState = 2
	StoragePoolDegraded     StoragePoolState = 3
	StoragePoolInaccessible StoragePoolState = 4
)

// StoragePoolBuildFlags as declared in libvirt/libvirt-storage.h:66
type StoragePoolBuildFlags int32

// StoragePoolBuildFlags enumeration from libvirt/libvirt-storage.h:66
const (
	StoragePoolBuildNew         StoragePoolBuildFlags = iota
	StoragePoolBuildRepair      StoragePoolBuildFlags = 1
	StoragePoolBuildResize      StoragePoolBuildFlags = 2
	StoragePoolBuildNoOverwrite StoragePoolBuildFlags = 4
	StoragePoolBuildOverwrite   StoragePoolBuildFlags = 8
)

// StoragePoolDeleteFlags as declared in libvirt/libvirt-storage.h:71
type StoragePoolDeleteFlags int32

// StoragePoolDeleteFlags enumeration from libvirt/libvirt-storage.h:71
const (
	StoragePoolDeleteNormal StoragePoolDeleteFlags = iota
	StoragePoolDeleteZeroed StoragePoolDeleteFlags = 1
)

// StoragePoolCreateFlags as declared in libvirt/libvirt-storage.h:88
type StoragePoolCreateFlags int32

// StoragePoolCreateFlags enumeration from libvirt/libvirt-storage.h:88
const (
	StoragePoolCreateNormal               StoragePoolCreateFlags = iota
	StoragePoolCreateWithBuild            StoragePoolCreateFlags = 1
	StoragePoolCreateWithBuildOverwrite   StoragePoolCreateFlags = 2
	StoragePoolCreateWithBuildNoOverwrite StoragePoolCreateFlags = 4
)

// StorageVolType as declared in libvirt/libvirt-storage.h:130
type StorageVolType int32

// StorageVolType enumeration from libvirt/libvirt-storage.h:130
const (
	StorageVolFile    StorageVolType = iota
	StorageVolBlock   StorageVolType = 1
	StorageVolDir     StorageVolType = 2
	StorageVolNetwork StorageVolType = 3
	StorageVolNetdir  StorageVolType = 4
	StorageVolPloop   StorageVolType = 5
)

// StorageVolDeleteFlags as declared in libvirt/libvirt-storage.h:136
type StorageVolDeleteFlags int32

// StorageVolDeleteFlags enumeration from libvirt/libvirt-storage.h:136
const (
	StorageVolDeleteNormal        StorageVolDeleteFlags = iota
	StorageVolDeleteZeroed        StorageVolDeleteFlags = 1
	StorageVolDeleteWithSnapshots StorageVolDeleteFlags = 2
)

// StorageVolWipeAlgorithm as declared in libvirt/libvirt-storage.h:168
type StorageVolWipeAlgorithm int32

// StorageVolWipeAlgorithm enumeration from libvirt/libvirt-storage.h:168
const (
	StorageVolWipeAlgZero       StorageVolWipeAlgorithm = iota
	StorageVolWipeAlgNnsa       StorageVolWipeAlgorithm = 1
	StorageVolWipeAlgDod        StorageVolWipeAlgorithm = 2
	StorageVolWipeAlgBsi        StorageVolWipeAlgorithm = 3
	StorageVolWipeAlgGutmann    StorageVolWipeAlgorithm = 4
	StorageVolWipeAlgSchneier   StorageVolWipeAlgorithm = 5
	StorageVolWipeAlgPfitzner7  StorageVolWipeAlgorithm = 6
	StorageVolWipeAlgPfitzner33 StorageVolWipeAlgorithm = 7
	StorageVolWipeAlgRandom     StorageVolWipeAlgorithm = 8
	StorageVolWipeAlgTrim       StorageVolWipeAlgorithm = 9
)

// StorageVolInfoFlags as declared in libvirt/libvirt-storage.h:176
type StorageVolInfoFlags int32

// StorageVolInfoFlags enumeration from libvirt/libvirt-storage.h:176
const (
	StorageVolUseAllocation StorageVolInfoFlags = iota
	StorageVolGetPhysical   StorageVolInfoFlags = 1
)

// StorageXMLFlags as declared in libvirt/libvirt-storage.h:190
type StorageXMLFlags int32

// StorageXMLFlags enumeration from libvirt/libvirt-storage.h:190
const (
	StorageXMLInactive StorageXMLFlags = 1
)

// ConnectListAllStoragePoolsFlags as declared in libvirt/libvirt-storage.h:244
type ConnectListAllStoragePoolsFlags int32

// ConnectListAllStoragePoolsFlags enumeration from libvirt/libvirt-storage.h:244
const (
	ConnectListStoragePoolsInactive    ConnectListAllStoragePoolsFlags = 1
	ConnectListStoragePoolsActive      ConnectListAllStoragePoolsFlags = 2
	ConnectListStoragePoolsPersistent  ConnectListAllStoragePoolsFlags = 4
	ConnectListStoragePoolsTransient   ConnectListAllStoragePoolsFlags = 8
	ConnectListStoragePoolsAutostart   ConnectListAllStoragePoolsFlags = 16
	ConnectListStoragePoolsNoAutostart ConnectListAllStoragePoolsFlags = 32
	ConnectListStoragePoolsDir         ConnectListAllStoragePoolsFlags = 64
	ConnectListStoragePoolsFs          ConnectListAllStoragePoolsFlags = 128
	ConnectListStoragePoolsNetfs       ConnectListAllStoragePoolsFlags = 256
	ConnectListStoragePoolsLogical     ConnectListAllStoragePoolsFlags = 512
	ConnectListStoragePoolsDisk        ConnectListAllStoragePoolsFlags = 1024
	ConnectListStoragePoolsIscsi       ConnectListAllStoragePoolsFlags = 2048
	ConnectListStoragePoolsScsi        ConnectListAllStoragePoolsFlags = 4096
	ConnectListStoragePoolsMpath       ConnectListAllStoragePoolsFlags = 8192
	ConnectListStoragePoolsRbd         ConnectListAllStoragePoolsFlags = 16384
	ConnectListStoragePoolsSheepdog    ConnectListAllStoragePoolsFlags = 32768
	ConnectListStoragePoolsGluster     ConnectListAllStoragePoolsFlags = 65536
	ConnectListStoragePoolsZfs         ConnectListAllStoragePoolsFlags = 131072
	ConnectListStoragePoolsVstorage    ConnectListAllStoragePoolsFlags = 262144
)

// StorageVolCreateFlags as declared in libvirt/libvirt-storage.h:342
type StorageVolCreateFlags int32

// StorageVolCreateFlags enumeration from libvirt/libvirt-storage.h:342
const (
	StorageVolCreatePreallocMetadata StorageVolCreateFlags = 1
	StorageVolCreateReflink          StorageVolCreateFlags = 2
)

// StorageVolDownloadFlags as declared in libvirt/libvirt-storage.h:354
type StorageVolDownloadFlags int32

// StorageVolDownloadFlags enumeration from libvirt/libvirt-storage.h:354
const (
	StorageVolDownloadSparseStream StorageVolDownloadFlags = 1
)

// StorageVolUploadFlags as declared in libvirt/libvirt-storage.h:363
type StorageVolUploadFlags int32

// StorageVolUploadFlags enumeration from libvirt/libvirt-storage.h:363
const (
	StorageVolUploadSparseStream StorageVolUploadFlags = 1
)

// StorageVolResizeFlags as declared in libvirt/libvirt-storage.h:394
type StorageVolResizeFlags int32

// StorageVolResizeFlags enumeration from libvirt/libvirt-storage.h:394
const (
	StorageVolResizeAllocate StorageVolResizeFlags = 1
	StorageVolResizeDelta    StorageVolResizeFlags = 2
	StorageVolResizeShrink   StorageVolResizeFlags = 4
)

// StoragePoolEventID as declared in libvirt/libvirt-storage.h:430
type StoragePoolEventID int32

// StoragePoolEventID enumeration from libvirt/libvirt-storage.h:430
const (
	StoragePoolEventIDLifecycle StoragePoolEventID = iota
	StoragePoolEventIDRefresh   StoragePoolEventID = 1
)

// StoragePoolEventLifecycleType as declared in libvirt/libvirt-storage.h:476
type StoragePoolEventLifecycleType int32

// StoragePoolEventLifecycleType enumeration from libvirt/libvirt-storage.h:476
const (
	StoragePoolEventDefined   StoragePoolEventLifecycleType = iota
	StoragePoolEventUndefined StoragePoolEventLifecycleType = 1
	StoragePoolEventStarted   StoragePoolEventLifecycleType = 2
	StoragePoolEventStopped   StoragePoolEventLifecycleType = 3
	StoragePoolEventCreated   StoragePoolEventLifecycleType = 4
	StoragePoolEventDeleted   StoragePoolEventLifecycleType = 5
)

// StreamFlags as declared in libvirt/libvirt-stream.h:34
type StreamFlags int32

// StreamFlags enumeration from libvirt/libvirt-stream.h:34
const (
	StreamNonblock StreamFlags = 1
)

// StreamRecvFlagsValues as declared in libvirt/libvirt-stream.h:50
type StreamRecvFlagsValues int32

// StreamRecvFlagsValues enumeration from libvirt/libvirt-stream.h:50
const (
	StreamRecvStopAtHole StreamRecvFlagsValues = 1
)

// StreamEventType as declared in libvirt/libvirt-stream.h:238
type StreamEventType int32

// StreamEventType enumeration from libvirt/libvirt-stream.h:238
const (
	StreamEventReadable StreamEventType = 1
	StreamEventWritable StreamEventType = 2
	StreamEventError    StreamEventType = 4
	StreamEventHangup   StreamEventType = 8
)

// errorLevel as declared in libvirt/virterror.h:44
type errorLevel int32

// errorLevel enumeration from libvirt/virterror.h:44
const (
	errNone    errorLevel = iota
	errWarning errorLevel = 1
	errError   errorLevel = 2
)

// errorDomain as declared in libvirt/virterror.h:140
type errorDomain int32

// errorDomain enumeration from libvirt/virterror.h:140
const (
	fromNone           errorDomain = iota
	fromXen            errorDomain = 1
	fromXend           errorDomain = 2
	fromXenstore       errorDomain = 3
	fromSexpr          errorDomain = 4
	fromXML            errorDomain = 5
	fromDom            errorDomain = 6
	fromRPC            errorDomain = 7
	fromProxy          errorDomain = 8
	fromConf           errorDomain = 9
	fromQemu           errorDomain = 10
	fromNet            errorDomain = 11
	fromTest           errorDomain = 12
	fromRemote         errorDomain = 13
	fromOpenvz         errorDomain = 14
	fromXenxm          errorDomain = 15
	fromStatsLinux     errorDomain = 16
	fromLxc            errorDomain = 17
	fromStorage        errorDomain = 18
	fromNetwork        errorDomain = 19
	fromDomain         errorDomain = 20
	fromUml            errorDomain = 21
	fromNodedev        errorDomain = 22
	fromXenInotify     errorDomain = 23
	fromSecurity       errorDomain = 24
	fromVbox           errorDomain = 25
	fromInterface      errorDomain = 26
	fromOne            errorDomain = 27
	fromEsx            errorDomain = 28
	fromPhyp           errorDomain = 29
	fromSecret         errorDomain = 30
	fromCPU            errorDomain = 31
	fromXenapi         errorDomain = 32
	fromNwfilter       errorDomain = 33
	fromHook           errorDomain = 34
	fromDomainSnapshot errorDomain = 35
	fromAudit          errorDomain = 36
	fromSysinfo        errorDomain = 37
	fromStreams        errorDomain = 38
	fromVmware         errorDomain = 39
	fromEvent          errorDomain = 40
	fromLibxl          errorDomain = 41
	fromLocking        errorDomain = 42
	fromHyperv         errorDomain = 43
	fromCapabilities   errorDomain = 44
	fromURI            errorDomain = 45
	fromAuth           errorDomain = 46
	fromDbus           errorDomain = 47
	fromParallels      errorDomain = 48
	fromDevice         errorDomain = 49
	fromSSH            errorDomain = 50
	fromLockspace      errorDomain = 51
	fromInitctl        errorDomain = 52
	fromIdentity       errorDomain = 53
	fromCgroup         errorDomain = 54
	fromAccess         errorDomain = 55
	fromSystemd        errorDomain = 56
	fromBhyve          errorDomain = 57
	fromCrypto         errorDomain = 58
	fromFirewall       errorDomain = 59
	fromPolkit         errorDomain = 60
	fromThread         errorDomain = 61
	fromAdmin          errorDomain = 62
	fromLogging        errorDomain = 63
	fromXenxl          errorDomain = 64
	fromPerf           errorDomain = 65
	fromLibssh         errorDomain = 66
	fromResctrl        errorDomain = 67
)

// errorNumber as declared in libvirt/virterror.h:326
type errorNumber int32

// errorNumber enumeration from libvirt/virterror.h:326
const (
	errOk                     errorNumber = iota
	errInternalError          errorNumber = 1
	errNoMemory               errorNumber = 2
	errNoSupport              errorNumber = 3
	errUnknownHost            errorNumber = 4
	errNoConnect              errorNumber = 5
	errInvalidConn            errorNumber = 6
	errInvalidDomain          errorNumber = 7
	errInvalidArg             errorNumber = 8
	errOperationFailed        errorNumber = 9
	errGetFailed              errorNumber = 10
	errPostFailed             errorNumber = 11
	errHTTPError              errorNumber = 12
	errSexprSerial            errorNumber = 13
	errNoXen                  errorNumber = 14
	errXenCall                errorNumber = 15
	errOsType                 errorNumber = 16
	errNoKernel               errorNumber = 17
	errNoRoot                 errorNumber = 18
	errNoSource               errorNumber = 19
	errNoTarget               errorNumber = 20
	errNoName                 errorNumber = 21
	errNoOs                   errorNumber = 22
	errNoDevice               errorNumber = 23
	errNoXenstore             errorNumber = 24
	errDriverFull             errorNumber = 25
	errCallFailed             errorNumber = 26
	errXMLError               errorNumber = 27
	errDomExist               errorNumber = 28
	errOperationDenied        errorNumber = 29
	errOpenFailed             errorNumber = 30
	errReadFailed             errorNumber = 31
	errParseFailed            errorNumber = 32
	errConfSyntax             errorNumber = 33
	errWriteFailed            errorNumber = 34
	errXMLDetail              errorNumber = 35
	errInvalidNetwork         errorNumber = 36
	errNetworkExist           errorNumber = 37
	errSystemError            errorNumber = 38
	errRPC                    errorNumber = 39
	errGnutlsError            errorNumber = 40
	warNoNetwork              errorNumber = 41
	errNoDomain               errorNumber = 42
	errNoNetwork              errorNumber = 43
	errInvalidMac             errorNumber = 44
	errAuthFailed             errorNumber = 45
	errInvalidStoragePool     errorNumber = 46
	errInvalidStorageVol      errorNumber = 47
	warNoStorage              errorNumber = 48
	errNoStoragePool          errorNumber = 49
	errNoStorageVol           errorNumber = 50
	warNoNode                 errorNumber = 51
	errInvalidNodeDevice      errorNumber = 52
	errNoNodeDevice           errorNumber = 53
	errNoSecurityModel        errorNumber = 54
	errOperationInvalid       errorNumber = 55
	warNoInterface            errorNumber = 56
	errNoInterface            errorNumber = 57
	errInvalidInterface       errorNumber = 58
	errMultipleInterfaces     errorNumber = 59
	warNoNwfilter             errorNumber = 60
	errInvalidNwfilter        errorNumber = 61
	errNoNwfilter             errorNumber = 62
	errBuildFirewall          errorNumber = 63
	warNoSecret               errorNumber = 64
	errInvalidSecret          errorNumber = 65
	errNoSecret               errorNumber = 66
	errConfigUnsupported      errorNumber = 67
	errOperationTimeout       errorNumber = 68
	errMigratePersistFailed   errorNumber = 69
	errHookScriptFailed       errorNumber = 70
	errInvalidDomainSnapshot  errorNumber = 71
	errNoDomainSnapshot       errorNumber = 72
	errInvalidStream          errorNumber = 73
	errArgumentUnsupported    errorNumber = 74
	errStorageProbeFailed     errorNumber = 75
	errStoragePoolBuilt       errorNumber = 76
	errSnapshotRevertRisky    errorNumber = 77
	errOperationAborted       errorNumber = 78
	errAuthCancelled          errorNumber = 79
	errNoDomainMetadata       errorNumber = 80
	errMigrateUnsafe          errorNumber = 81
	errOverflow               errorNumber = 82
	errBlockCopyActive        errorNumber = 83
	errOperationUnsupported   errorNumber = 84
	errSSH                    errorNumber = 85
	errAgentUnresponsive      errorNumber = 86
	errResourceBusy           errorNumber = 87
	errAccessDenied           errorNumber = 88
	errDbusService            errorNumber = 89
	errStorageVolExist        errorNumber = 90
	errCPUIncompatible        errorNumber = 91
	errXMLInvalidSchema       errorNumber = 92
	errMigrateFinishOk        errorNumber = 93
	errAuthUnavailable        errorNumber = 94
	errNoServer               errorNumber = 95
	errNoClient               errorNumber = 96
	errAgentUnsynced          errorNumber = 97
	errLibssh                 errorNumber = 98
	errDeviceMissing          errorNumber = 99
	errInvalidNwfilterBinding errorNumber = 100
	errNoNwfilterBinding      errorNumber = 101
)
