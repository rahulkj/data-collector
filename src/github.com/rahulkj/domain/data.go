package domain

import (
  "code.google.com/p/log4go"
)


// book model
type Data struct {
  Environment  string `json:"environment"`
  Id int `json:"id"`
  VsphereInfo VsphereInfo `json:"vSphereInfo"`
  NetworkInfo NetworkInfo `json:"networkInfo"`
  EnvInfo EnvInfo `json:"envInfo"`
  MandatoryInfo MandatoryInfo `json:"mandatoryInfo"`
  OptionalInfo OptionalInfo `json:"optionalInfo"`
}

type VsphereInfo struct {
  VsphereAddress string `json:"vSphereAddress"`
  VsphereUserName string `json:"vSphereUserName"`
  VspherePassword string `json:"vSpherePassword"`
  SvrUserName string `json:"svrUserName"`
  SvrPassword string `json:"svrPassword"`
}

type NetworkInfo struct {
  NetMask string `json:"netMask"`
  DefaultGateway string `json:"defaultGateway"`
  DnsServers string `json:"dnsServers"`
  NtpServers string `json:"ntpServers"`
  VsphereSubnet string `json:"vSphereSubnet"`
}

type EnvInfo struct {
  DataCenterName string `json:"dataCenterName"`
  ClusterName string `json:"clusterName"`
  DataStoreNames string `json:"dataStoreNames"`
  ResourcePoolName string `json:"resourcePoolName"`
  NetworkName string `json:"networkName"`
}

type MandatoryInfo struct {
  OpsMgrIPAddress string `json:"opsMgrIPAddress"`
  OpsMgrUserName string `json:"opsMgrUserName"`
  OpsMgrPassword string `json:"opsMgrPassword"`
  ExternalLoadBalancer bool `json:"externalLoadBalancer"`
  HAproxyIPs string `json:"haproxyIPs"`
  ExternalLBIPs string `json:"externalLBIPs"`
  ExternalAppsDomain string `json:"externalAppsDomain"`
  ExternalSystemDomain string `json:"externalSystemDomain"`
  RouterIPs string `json:"routerIPs"`
  ExcludedIPRanges string `json:"excludedIPRanges"`
  SystemDomain string `json:"systemDomain"`
  ApplicationDomain string `json:"applicationDomain"`
  PublicCert string `json:"publicCert"`
  PrivateCert string `json:"privateCert"`
}

type OptionalInfo struct {
  CcDBEncryptionKey string `json:"ccDBEncryptionKey"`
  MaxFileSize string `json:"maxFileSize"`
  SsoURL string `json:"ssoURL"`
  Email Email `json:"email"`
  AppMemory string `json:"appMemory"`
  ServiceInstances string `json:"serviceInstances"`
}

type Email struct {
  ReplyToEmail string `json:"replyToEmail"`
  FromEmail string `json:"fromEmail"`
  SmtpServerAddress string `json:"smtpServerAddress"`
  SmtpServerPort string `json:"smtpServerPort"`
  HeloDomain string `json:"heloDomain"`
  SmtpAuthRequired bool `json:"smtpAuthRequired"`
  SmtpServerUsername string `json:"smtpServerUsername"`
  SmtpServerPassword string `json:"smtpServerPassword"`
}

var id = 0
func getNextId() int {
  id += 1
  return id
}

func SaveData(data Data) Data {
  log := make(log4go.Logger)
  log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

  if data.Id == 0 {
    log.Debug("Generating ID as thats missing\n");
    data.Id = getNextId()
  }

  log.Debug("Save Data Invoked %v\n",data);

  return data
}
