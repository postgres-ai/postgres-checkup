package a008

import (
	"fmt"
	"testing"

	checkup ".."
)

func printConclusions(result checkup.ReportOutcome) {
	for _, conclusion := range result.Conclusions {
		fmt.Println("C:  ", conclusion)
	}
}

func printReccomendations(result checkup.ReportOutcome) {
	for _, recommendation := range result.Recommendations {
		fmt.Println("R:  ", recommendation)
	}
}

func TestA008Sucess(t *testing.T) {
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
	if result.P1 || result.P2 || result.P3 {
		t.Fatal("TestSucess failed")
	}
	printConclusions(result)
	printReccomendations(result)
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
	if !result.P1 {
		t.Fatal("TestNfsFileSystem failed")
	}
	printConclusions(result)
	printReccomendations(result)
}

func TestA008NotReqFileSystem(t *testing.T) {
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
	if !result.P3 {
		t.Fatal("TestNotReqFileSystem failed")
	}
	printConclusions(result)
	printReccomendations(result)
}

func TestA008UserPercentCritical(t *testing.T) {
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
	if !result.P1 {
		t.Fatal("TestUserPercentCritical failed")
	}
	printConclusions(result)
	printReccomendations(result)
}

func TestA008UserPercentProblem(t *testing.T) {
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
	if !result.P2 {
		t.Fatal("TestUserPercentProblem failed")
	}
	printConclusions(result)
	printReccomendations(result)
}
