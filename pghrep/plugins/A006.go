package main

import (
	"../src/pyraconv"
)

var Data map[string]interface{}

type prepare string

func compareHostsData(data map[string]interface{}) {
	hosts := pyraconv.ToInterfaceMap(data["hosts"])
	master := pyraconv.ToString(hosts["master"])
	replicas := pyraconv.ToStringArray(hosts["replicas"])
	diffData := make(map[string]interface{})

	results := pyraconv.ToInterfaceMap(data["results"])
	masterData := pyraconv.ToInterfaceMap(results[master])
	masterData = pyraconv.ToInterfaceMap(masterData["data"])

	pgSettingsData := pyraconv.ToInterfaceMap(masterData["pg_settings"])
	diffSettings := make(map[string]interface{})
	for settingName, value := range pgSettingsData {
		valueData := pyraconv.ToInterfaceMap(value)
		settingValue := pyraconv.ToString(valueData["setting"])
		settingUnit := pyraconv.ToString(valueData["unit"])
		var diffSetting []interface{}
		masterSetting := make(map[string]interface{})
		masterSetting["value"] = settingValue
		masterSetting["unit"] = settingUnit
		diffSetting = append(diffSetting, masterSetting)
		diff := false
		for _, replica := range replicas {
			rSettingValue, rSettingUnit := getReplicaSettingValue(data, replica, settingName)
			if rSettingValue != "null" && rSettingUnit != "null" {
				// Process only hosts which have data
				hostSetting := make(map[string]interface{})
				hostSetting["value"] = rSettingValue
				hostSetting["unit"] = rSettingUnit
				diffSetting = append(diffSetting, hostSetting)
				if (settingValue != rSettingValue) || (settingUnit != rSettingUnit) {
					diff = true
				}
			} else {
				hostSetting := make(map[string]interface{})
				hostSetting["value"] = ""
				hostSetting["unit"] = ""
				diffSetting = append(diffSetting, hostSetting)
			}
		}
		if diff {
			diffSettings[settingName] = diffSetting
		}
	}
	diffData["pg_settings"] = diffSettings

	pgConfigData := pyraconv.ToInterfaceMap(masterData["pg_config"])
	diffConfigs := make(map[string]interface{})
	for configName, value := range pgConfigData {
		valueData := pyraconv.ToInterfaceMap(value)
		settingValue := pyraconv.ToString(valueData["setting"])
		settingUnit := pyraconv.ToString(valueData["unit"])
		var diffConfig []interface{}
		masterSetting := make(map[string]interface{})
		masterSetting["value"] = settingValue
		diffConfig = append(diffConfig, masterSetting)
		diff := false
		for _, replica := range replicas {
			rSettingValue, rSettingUnit := getReplicaConfigValue(data, replica, configName)
			if rSettingValue != "null" && rSettingUnit != "null" {
				hostSetting := make(map[string]interface{})
				hostSetting["value"] = rSettingValue
				hostSetting["unit"] = rSettingUnit
				diffConfig = append(diffConfig, hostSetting)
				if (settingValue != rSettingValue) || (settingUnit != rSettingUnit) {
					diff = true
				}
			} else {
				hostSetting := make(map[string]interface{})
				hostSetting["value"] = ""
				hostSetting["unit"] = ""
				diffConfig = append(diffConfig, hostSetting)
			}
		}
		if diff {
			diffConfigs[configName] = diffConfig
		}
	}
	diffData["pg_configs"] = diffConfigs

	data["diffData"] = diffData
}

func getReplicaSettingValue(data map[string]interface{}, replica string, settingName string) (string, string) {
	results := pyraconv.ToInterfaceMap(data["results"])
	if results[replica] != nil {
		hostData := pyraconv.ToInterfaceMap(results[replica])
		hostData = pyraconv.ToInterfaceMap(hostData["data"])
		pgSettingsData := pyraconv.ToInterfaceMap(hostData["pg_settings"])
		pgSettingData := pyraconv.ToInterfaceMap(pgSettingsData[settingName])
		return pyraconv.ToString(pgSettingData["setting"]), pyraconv.ToString(pgSettingData["unit"])
	}
	return "null", "null"
}

func getReplicaConfigValue(data map[string]interface{}, replica string, settingName string) (string, string) {
	results := pyraconv.ToInterfaceMap(data["results"])
	if results[replica] != nil {
		hostData := pyraconv.ToInterfaceMap(results[replica])
		hostData = pyraconv.ToInterfaceMap(hostData["data"])
		pgConfigsData := pyraconv.ToInterfaceMap(hostData["pg_config"])
		pgСonfigData := pyraconv.ToInterfaceMap(pgConfigsData[settingName])
		return pyraconv.ToString(pgСonfigData["setting"]), pyraconv.ToString(pgСonfigData["unit"])
	}
	return "null", "null"
}


func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
	compareHostsData(data)
	return data
}

var Preparer prepare
