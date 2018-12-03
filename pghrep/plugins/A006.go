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
        valueData := pyraconv.ToInterfaceMap(value);
        settingValue := pyraconv.ToString(valueData["setting"])
        settingUnit := pyraconv.ToString(valueData["unit"])
        diffSetting := make(map[string]interface{})
        masterSetting := make(map[string]interface{})
        masterSetting["value"] = settingValue
        masterSetting["unit"] = settingUnit
        diffSetting["master"] = masterSetting
        diff := false
        for _, replica := range replicas {
            rSettingValue, rSettingUnit := getReplicaSettingValue(data, replica, settingName)
            hostSetting := make(map[string]interface{})
            hostSetting["value"] = rSettingValue
            hostSetting["unit"] = rSettingUnit
            diffSetting[replica] = hostSetting
            if (settingValue != rSettingValue) || (settingUnit != rSettingUnit) {
                diff = true
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
        valueData := pyraconv.ToInterfaceMap(value);
        settingValue := pyraconv.ToString(valueData["setting"])
        settingUnit := pyraconv.ToString(valueData["unit"])
        diffConfig := make(map[string]interface{})
        masterSetting := make(map[string]interface{})
        masterSetting["value"] = settingValue
        diffConfig["master"] = masterSetting
        diff := false
        for _, replica := range replicas {
            rSettingValue, rSettingUnit := getReplicaConfigValue(data, replica, configName)
            hostSetting := make(map[string]interface{})
            hostSetting["value"] = rSettingValue
            hostSetting["unit"] = rSettingUnit
            diffConfig[replica] = hostSetting
            if (settingValue != rSettingValue) || (settingUnit != rSettingUnit) {
                diff = true
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
    hostData := pyraconv.ToInterfaceMap(results[replica])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    pgSettingsData := pyraconv.ToInterfaceMap(hostData["pg_settings"])
    pgSettingData := pyraconv.ToInterfaceMap(pgSettingsData[settingName])
    return pyraconv.ToString(pgSettingData["setting"]), pyraconv.ToString(pgSettingData["unit"])
}

func getReplicaConfigValue(data map[string]interface{}, replica string, settingName string) (string, string) {
    results := pyraconv.ToInterfaceMap(data["results"])
    hostData := pyraconv.ToInterfaceMap(results[replica])
    hostData = pyraconv.ToInterfaceMap(hostData["data"])
    pgConfigsData := pyraconv.ToInterfaceMap(hostData["pg_config"])
    pgСonfigData := pyraconv.ToInterfaceMap(pgConfigsData[settingName])
    return pyraconv.ToString(pgСonfigData["setting"]), pyraconv.ToString(pgСonfigData["unit"])
}

func (g prepare) Prepare(data map[string]interface{}) map[string]interface{} {
    compareHostsData(data)
    return data
}

var Preparer prepare