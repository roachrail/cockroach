// Code generated by "stringer"; DO NOT EDIT.

package clusterversion

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[V21_2-0]
	_ = x[Start22_1-1]
	_ = x[TargetBytesAvoidExcess-2]
	_ = x[TraceIDDoesntImplyStructuredRecording-3]
	_ = x[AlterSystemTableStatisticsAddAvgSizeCol-4]
	_ = x[MVCCAddSSTable-5]
	_ = x[InsertPublicSchemaNamespaceEntryOnRestore-6]
	_ = x[UnsplitRangesInAsyncGCJobs-7]
	_ = x[ValidateGrantOption-8]
	_ = x[PebbleFormatBlockPropertyCollector-9]
	_ = x[ProbeRequest-10]
	_ = x[SelectRPCsTakeTracingInfoInband-11]
	_ = x[PreSeedTenantSpanConfigs-12]
	_ = x[SeedTenantSpanConfigs-13]
	_ = x[PublicSchemasWithDescriptors-14]
	_ = x[EnsureSpanConfigReconciliation-15]
	_ = x[EnsureSpanConfigSubscription-16]
	_ = x[EnableSpanConfigStore-17]
	_ = x[ScanWholeRows-18]
	_ = x[SCRAMAuthentication-19]
	_ = x[UnsafeLossOfQuorumRecoveryRangeLog-20]
	_ = x[AlterSystemProtectedTimestampAddColumn-21]
	_ = x[EnableProtectedTimestampsForTenant-22]
	_ = x[DeleteCommentsWithDroppedIndexes-23]
	_ = x[RemoveIncompatibleDatabasePrivileges-24]
	_ = x[AddRaftAppliedIndexTermMigration-25]
	_ = x[PostAddRaftAppliedIndexTermMigration-26]
	_ = x[DontProposeWriteTimestampForLeaseTransfers-27]
	_ = x[EnablePebbleFormatVersionBlockProperties-28]
	_ = x[DisableSystemConfigGossipTrigger-29]
	_ = x[MVCCIndexBackfiller-30]
	_ = x[EnableLeaseHolderRemoval-31]
	_ = x[BackupResolutionInJob-32]
	_ = x[LooselyCoupledRaftLogTruncation-33]
	_ = x[ChangefeedIdleness-34]
	_ = x[BackupDoesNotOverwriteLatestAndCheckpoint-35]
	_ = x[EnableDeclarativeSchemaChanger-36]
	_ = x[RowLevelTTL-37]
	_ = x[PebbleFormatSplitUserKeysMarked-38]
	_ = x[IncrementalBackupSubdir-39]
	_ = x[DateStyleIntervalStyleCastRewrite-40]
	_ = x[EnableNewStoreRebalancer-41]
	_ = x[ClusterLocksVirtualTable-42]
	_ = x[AutoStatsTableSettings-43]
	_ = x[ForecastStats-44]
	_ = x[SuperRegions-45]
	_ = x[EnableNewChangefeedOptions-46]
	_ = x[SpanCountTable-47]
	_ = x[PreSeedSpanCountTable-48]
	_ = x[SeedSpanCountTable-49]
	_ = x[V22_1-50]
	_ = x[Start22_2-51]
	_ = x[LocalTimestamps-52]
	_ = x[EnsurePebbleFormatVersionRangeKeys-53]
	_ = x[EnablePebbleFormatVersionRangeKeys-54]
	_ = x[TrigramInvertedIndexes-55]
	_ = x[RemoveGrantPrivilege-56]
	_ = x[MVCCRangeTombstones-57]
	_ = x[UpgradeSequenceToBeReferencedByID-58]
	_ = x[SampledStmtDiagReqs-59]
	_ = x[AddSSTableTombstones-60]
	_ = x[SystemPrivilegesTable-61]
	_ = x[EnablePredicateProjectionChangefeed-62]
	_ = x[AlterSystemSQLInstancesAddLocality-63]
}

const _Key_name = "V21_2Start22_1TargetBytesAvoidExcessTraceIDDoesntImplyStructuredRecordingAlterSystemTableStatisticsAddAvgSizeColMVCCAddSSTableInsertPublicSchemaNamespaceEntryOnRestoreUnsplitRangesInAsyncGCJobsValidateGrantOptionPebbleFormatBlockPropertyCollectorProbeRequestSelectRPCsTakeTracingInfoInbandPreSeedTenantSpanConfigsSeedTenantSpanConfigsPublicSchemasWithDescriptorsEnsureSpanConfigReconciliationEnsureSpanConfigSubscriptionEnableSpanConfigStoreScanWholeRowsSCRAMAuthenticationUnsafeLossOfQuorumRecoveryRangeLogAlterSystemProtectedTimestampAddColumnEnableProtectedTimestampsForTenantDeleteCommentsWithDroppedIndexesRemoveIncompatibleDatabasePrivilegesAddRaftAppliedIndexTermMigrationPostAddRaftAppliedIndexTermMigrationDontProposeWriteTimestampForLeaseTransfersEnablePebbleFormatVersionBlockPropertiesDisableSystemConfigGossipTriggerMVCCIndexBackfillerEnableLeaseHolderRemovalBackupResolutionInJobLooselyCoupledRaftLogTruncationChangefeedIdlenessBackupDoesNotOverwriteLatestAndCheckpointEnableDeclarativeSchemaChangerRowLevelTTLPebbleFormatSplitUserKeysMarkedIncrementalBackupSubdirDateStyleIntervalStyleCastRewriteEnableNewStoreRebalancerClusterLocksVirtualTableAutoStatsTableSettingsForecastStatsSuperRegionsEnableNewChangefeedOptionsSpanCountTablePreSeedSpanCountTableSeedSpanCountTableV22_1Start22_2LocalTimestampsEnsurePebbleFormatVersionRangeKeysEnablePebbleFormatVersionRangeKeysTrigramInvertedIndexesRemoveGrantPrivilegeMVCCRangeTombstonesUpgradeSequenceToBeReferencedByIDSampledStmtDiagReqsAddSSTableTombstonesSystemPrivilegesTableEnablePredicateProjectionChangefeedAlterSystemSQLInstancesAddLocality"

var _Key_index = [...]uint16{0, 5, 14, 36, 73, 112, 126, 167, 193, 212, 246, 258, 289, 313, 334, 362, 392, 420, 441, 454, 473, 507, 545, 579, 611, 647, 679, 715, 757, 797, 829, 848, 872, 893, 924, 942, 983, 1013, 1024, 1055, 1078, 1111, 1135, 1159, 1181, 1194, 1206, 1232, 1246, 1267, 1285, 1290, 1299, 1314, 1348, 1382, 1404, 1424, 1443, 1476, 1495, 1515, 1536, 1571, 1605}

func (i Key) String() string {
	if i < 0 || i >= Key(len(_Key_index)-1) {
		return "Key(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Key_name[_Key_index[i]:_Key_index[i+1]]
}
