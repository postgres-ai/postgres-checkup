package main

import (
	"testing"
)

func TestSucess(t *testing.T) {
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"/var/lib/postgresql/": FsItem{
		Fstype:     "ext4",
		Size:       "",
		Avail:      "",
		Used:       "",
		UsePercent: "45%",
		MountPoint: "/var/lib/postgresql/",
		Path:       "/var/lib/postgresql/",
		Device:     "/dev/sdba0",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if result.P1 || result.P2 || result.P3 {
		t.Fatal("TestSucess failed")
	}
}

func TestNfsFileSystem(t *testing.T) {
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"nfs": FsItem{
		Fstype:     "nfs",
		Size:       "",
		Avail:      "",
		Used:       "",
		UsePercent: "",
		MountPoint: "/home/nfs",
		Path:       "/home/nfs",
		Device:     "/dev/nfs",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P1 {
		t.Fatal("TestNfsFileSystem failed")
	}
}

func TestNotReqFileSystem(t *testing.T) {
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{
		"/dev/zfs": FsItem{
			Fstype:     "zfs",
			Size:       "",
			Avail:      "",
			Used:       "",
			UsePercent: "",
			MountPoint: "/home/zfs",
			Path:       "",
			Device:     "/dev/zfs",
		},
		"/dev/jfs": FsItem{
			Fstype:     "jfs",
			Size:       "",
			Avail:      "",
			Used:       "",
			UsePercent: "",
			MountPoint: "/home/jfs",
			Path:       "",
			Device:     "/dev/jfs",
		},
	}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P3 {
		t.Fatal("TestNotReqFileSystem failed")
	}
}

func TestUserPercentCritical(t *testing.T) {
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"ext4": FsItem{
		Fstype:     "ext4",
		Size:       "",
		Avail:      "",
		Used:       "",
		UsePercent: "92%",
		MountPoint: "/home/ext4",
		Path:       "/home/ext4",
		Device:     "/dev/ext4",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P1 {
		t.Fatal("TestUserPercentCritical failed")
	}
}

func TestUserPercentProblem(t *testing.T) {
	var report A008Report
	var hostResult A008ReportHostResult
	hostResult.Data.DbData = map[string]FsItem{"ext4": FsItem{
		Fstype:     "ext4",
		Size:       "",
		Avail:      "",
		Used:       "",
		UsePercent: "72%",
		MountPoint: "/home/ext4",
		Path:       "/home/ext4",
		Device:     "/dev/ext4",
	}}
	report.Results = A008ReportHostsResults{"test-host": hostResult}
	result := A008Process(report)
	if !result.P2 {
		t.Fatal("TestUserPercentProblem failed")
	}
}
