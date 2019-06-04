package a008

import (
	"fmt"
	"testing"

	checkup ".."
)

func TestA008Success(t *testing.T) {
	fmt.Println(t.Name())
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"/var/lib/postgresql/": FsItem{
		Fstype:     "ext4",
		Size:       "",
		Avail:      "",
		Used:       "130G",
		UsePercent: "45%",
		MountPoint: "/var/lib/postgresql/",
		Path:       "/var/lib/postgresql/",
		Device:     "/dev/sdba0",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if result.P1 ||
		result.P2 ||
		result.P3 ||
		!checkup.ResultInList(result.Conclusions, A008_SPACE_USAGE_NORMAL) ||
		len(result.Recommendations) > 0 {
		t.Fatal("TestA008Success failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestA008NfsFileSystem(t *testing.T) {
	fmt.Println(t.Name())
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{
		"nfs": FsItem{
			Fstype:     "nfs",
			Size:       "",
			Avail:      "",
			Used:       "130G",
			UsePercent: "43%",
			MountPoint: "/home/nfs",
			Path:       "/home/nfs",
			Device:     "/dev/nfs",
		},
		"nfs2": FsItem{
			Fstype:     "nfs2",
			Size:       "",
			Avail:      "",
			Used:       "130G",
			UsePercent: "43%",
			MountPoint: "/home/nfs2",
			Path:       "/home/nfs2",
			Device:     "/dev/nfs2",
		}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, A008_NFS_DISK) ||
		!checkup.ResultInList(result.Recommendations, A008_NFS_DISK) {
		t.Fatal("TestNfsFileSystem failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestA008NotRecommendedFileSystem(t *testing.T) {
	fmt.Println(t.Name())
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{
		"/dev/zfs": FsItem{
			Fstype:     "zfs",
			Size:       "",
			Avail:      "",
			Used:       "130G",
			UsePercent: "23%",
			MountPoint: "/home/zfs",
			Path:       "",
			Device:     "/dev/zfs",
		},
		"/dev/jfs": FsItem{
			Fstype:     "jfs",
			Size:       "",
			Avail:      "",
			Used:       "130G",
			UsePercent: "45%",
			MountPoint: "/home/jfs",
			Path:       "",
			Device:     "/dev/jfs",
		},
	}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P3 ||
		!checkup.ResultInList(result.Conclusions, A008_NOT_RECOMMENDED_FS) ||
		!checkup.ResultInList(result.Recommendations, A008_NOT_RECOMMENDED_FS) {
		t.Fatal("TestA008NotRecommendedFileSystem failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestA008DiskUsageCritical(t *testing.T) {
	fmt.Println(t.Name())
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"ext4": FsItem{
		Fstype:     "ext4",
		Size:       "",
		Avail:      "",
		Used:       "130G",
		UsePercent: "92%",
		MountPoint: "/home/ext4",
		Path:       "/home/ext4",
		Device:     "/dev/ext4",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P1 ||
		!checkup.ResultInList(result.Conclusions, A008_SPACE_USAGE_CRITICAL) ||
		!checkup.ResultInList(result.Recommendations, A008_SPACE_USAGE_CRITICAL) {
		t.Fatal("TestA008DiskUsageCritical failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}

func TestA008DiskUsageWarning(t *testing.T) {
	fmt.Println(t.Name())
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"ext4": FsItem{
		Fstype:     "ext4",
		Size:       "",
		Avail:      "",
		Used:       "130G",
		UsePercent: "72%",
		MountPoint: "/home/ext4",
		Path:       "/home/ext4",
		Device:     "/dev/ext4",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P2 ||
		!checkup.ResultInList(result.Conclusions, A008_SPACE_USAGE_WARNING) ||
		!checkup.ResultInList(result.Recommendations, A008_SPACE_USAGE_WARNING) {
		t.Fatal("TestA008DiskUsageWarning failed")
	}
	checkup.PrintResultConclusions(result)
	checkup.PrintResultRecommendations(result)
}
