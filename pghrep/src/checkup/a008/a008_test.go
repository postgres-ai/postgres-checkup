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
		!checkup.InList(result.Conclusions, "No significant risks of out-of-disk-space problem have been detected.") ||
		!checkup.InList(result.Recommendations, "No recommendations.") {
		t.Fatal("TestA008Success failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
		!checkup.InList(result.Conclusions, "[P1] `/home/nfs`, `/home/nfs2` on host `test-host` are located on an NFS drive. This might lead to serious issues with Postgres, including downtime and data corruption.") ||
		!checkup.InList(result.Recommendations, "[P1] Do not use NFS for Postgres.") {
		t.Fatal("TestNfsFileSystem failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
		!checkup.InList(result.Conclusions, "[P3] `/home/zfs`, `/home/jfs` on host `test-host` are located on drives where the following filesystems are used: `zfs`, `jfs` respectively. This might mean that Postgres performance and reliability characteristics are worse than it could be in case of use of more popular filesystems (such as ext4).") ||
		!checkup.InList(result.Recommendations, "[P3] Consider using ext4 for all Postgres directories.") {
		t.Fatal("TestA008NotRecommendedFileSystem failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
		!checkup.InList(result.Conclusions, "[P1] Disk `/home/ext4` on `test-host` space usage is 130G, it exceeds 90%. There are significant risks of out-of-disk-space problem. In this case, PostgreSQL will stop working and manual fix will be required.") ||
		!checkup.InList(result.Recommendations, "[P1] Add more disk space to `/home/ext4` on `test-host` as soon as possible to prevent outage.") {
		t.Fatal("TestA008DiskUsageCritical failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
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
		!checkup.InList(result.Conclusions, "[P2] Disk `/home/ext4` on `test-host` space usage is 130G, it exceeds 70%. There are some risks of out-of-disk-space problem.") ||
		!checkup.InList(result.Recommendations, "[P2] Add more disk space to `/home/ext4` on `test-host`. It is recommended to keep free disk space more than 30% to reduce risks of out-of-disk-space problem.") {
		t.Fatal("TestA008DiskUsageWarning failed")
	}
	checkup.PrintConclusions(result)
	checkup.PrintRecommendations(result)
}
