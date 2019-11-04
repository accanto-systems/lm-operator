package alm

import (
	"strings"

	"github.com/go-logr/logr"
	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
)

type serviceDeploymentInfo struct {
	serviceName    string
	imageName      string
	imageVersion   string
	port           int32
	targetPort     int
	nodePort       int32
	numReplicas    int32
	cpuRequests    string
	cpuLimit       string
	memoryRequests string
	memoryLimits   string
	heap           string
}

type nimrodServiceDeploymentInfo struct {
	serviceDeploymentInfo
	themesConfigMap  string
	localesConfigMap string
}

type configuratorDeploymentInfo struct {
	serviceDeploymentInfo
	run bool
}

type deploymentInfo struct {
	configurator configuratorDeploymentInfo
	conductor    serviceDeploymentInfo
	apollo       serviceDeploymentInfo
	galileo      serviceDeploymentInfo
	talledega    serviceDeploymentInfo
	daytona      serviceDeploymentInfo
	nimrod       nimrodServiceDeploymentInfo
	ishtar       serviceDeploymentInfo
	relay        serviceDeploymentInfo
	watchtower   serviceDeploymentInfo
	doki         serviceDeploymentInfo
	brent        serviceDeploymentInfo
}

func createDeploymentInfo(instance *comv1alpha1.ALM, reqLogger logr.Logger) (deploymentInfo, error) {
	lmRelease, err := getLMRelease(instance.Spec.Release, reqLogger)
	if err != nil {
		return deploymentInfo{}, err
	}

	deploymentInfo := deploymentInfo{}
	deploymentInfo.configurator.serviceName = "lm-configurator"
	deploymentInfo.configurator.imageName = "lm-configurator"
	deploymentInfo.configurator.imageVersion = lmRelease.Configurator.Version
	deploymentInfo.configurator.numReplicas = int32(1)
	deploymentInfo.configurator.run = instance.Spec.Configurator.Run

	deploymentInfo.conductor.serviceName = "conductor"
	deploymentInfo.conductor.port = 8761
	deploymentInfo.conductor.targetPort = 8761
	deploymentInfo.conductor.nodePort = -1
	deploymentInfo.conductor.imageName = "conductor"
	deploymentInfo.conductor.imageVersion = lmRelease.Conductor.Version

	deploymentInfo.apollo.serviceName = "apollo"
	deploymentInfo.apollo.port = 8282
	deploymentInfo.apollo.targetPort = 8282
	deploymentInfo.apollo.nodePort = -1
	deploymentInfo.apollo.imageName = "apollo"
	deploymentInfo.apollo.imageVersion = lmRelease.Apollo.Version

	deploymentInfo.galileo.serviceName = "galileo"
	deploymentInfo.galileo.port = 8283
	deploymentInfo.galileo.targetPort = 8283
	deploymentInfo.galileo.nodePort = -1
	deploymentInfo.galileo.imageName = "galileo"
	deploymentInfo.galileo.imageVersion = lmRelease.Galileo.Version

	deploymentInfo.talledega.serviceName = "talledega"
	deploymentInfo.talledega.port = 8287
	deploymentInfo.talledega.targetPort = 8287
	deploymentInfo.talledega.nodePort = -1
	deploymentInfo.talledega.imageName = "talledega"
	deploymentInfo.talledega.imageVersion = lmRelease.Talledega.Version

	deploymentInfo.daytona.serviceName = "daytona"
	deploymentInfo.daytona.port = 8281
	deploymentInfo.daytona.targetPort = 8281
	deploymentInfo.daytona.nodePort = -1
	deploymentInfo.daytona.imageName = "daytona"
	deploymentInfo.daytona.imageVersion = lmRelease.Daytona.Version

	deploymentInfo.nimrod.serviceName = "nimrod"
	deploymentInfo.nimrod.port = 8290
	deploymentInfo.nimrod.targetPort = 8290
	deploymentInfo.nimrod.nodePort = -1
	deploymentInfo.nimrod.imageName = "nimrod"
	deploymentInfo.nimrod.imageVersion = lmRelease.Nimrod.Version

	deploymentInfo.ishtar.serviceName = "ishtar"
	deploymentInfo.ishtar.port = 8280
	deploymentInfo.ishtar.targetPort = 8280
	deploymentInfo.ishtar.nodePort = -1
	deploymentInfo.ishtar.imageName = "ishtar"
	deploymentInfo.ishtar.imageVersion = lmRelease.Ishtar.Version

	deploymentInfo.relay.serviceName = "relay"
	deploymentInfo.relay.port = 8285
	deploymentInfo.relay.targetPort = 8285
	deploymentInfo.relay.nodePort = -1
	deploymentInfo.relay.imageName = "relay"
	deploymentInfo.relay.imageVersion = lmRelease.Relay.Version

	deploymentInfo.watchtower.serviceName = "watchtower"
	deploymentInfo.watchtower.port = 8284
	deploymentInfo.watchtower.targetPort = 8284
	deploymentInfo.watchtower.nodePort = -1
	deploymentInfo.watchtower.imageName = "watchtower"
	deploymentInfo.watchtower.imageVersion = lmRelease.Watchtower.Version

	deploymentInfo.doki.serviceName = "doki"
	deploymentInfo.doki.port = 8288
	deploymentInfo.doki.targetPort = 8288
	deploymentInfo.doki.nodePort = -1
	deploymentInfo.doki.imageName = "doki"
	deploymentInfo.doki.imageVersion = lmRelease.Doki.Version

	deploymentInfo.brent.serviceName = "brent"
	deploymentInfo.brent.port = 8291
	deploymentInfo.brent.targetPort = 8291
	deploymentInfo.brent.nodePort = -1
	deploymentInfo.brent.imageName = "brent"
	deploymentInfo.brent.imageVersion = lmRelease.Brent.Version

	if strings.ToLower(instance.Spec.DeploymentType) == "ha" {
		deploymentInfo.conductor.numReplicas = int32(3)
		deploymentInfo.conductor.cpuRequests = "1"
		deploymentInfo.conductor.memoryRequests = "1024Mi"
		deploymentInfo.conductor.heap = "1G"

		deploymentInfo.apollo.numReplicas = int32(3)
		deploymentInfo.apollo.cpuRequests = "1"
		deploymentInfo.apollo.memoryRequests = "1024Mi"
		deploymentInfo.apollo.heap = "1G"

		deploymentInfo.galileo.numReplicas = int32(3)
		deploymentInfo.galileo.cpuRequests = "1"
		deploymentInfo.galileo.memoryRequests = "1024Mi"
		deploymentInfo.galileo.heap = "1G"

		deploymentInfo.talledega.numReplicas = int32(3)
		deploymentInfo.talledega.cpuRequests = "1"
		deploymentInfo.talledega.memoryRequests = "1024Mi"
		deploymentInfo.talledega.heap = "1G"

		deploymentInfo.daytona.numReplicas = int32(3)
		deploymentInfo.daytona.cpuRequests = "1"
		deploymentInfo.daytona.memoryRequests = "1024Mi"
		deploymentInfo.daytona.heap = "1G"

		deploymentInfo.nimrod.numReplicas = int32(3)
		deploymentInfo.nimrod.cpuRequests = "1"
		deploymentInfo.nimrod.memoryRequests = "1024Mi"
		deploymentInfo.nimrod.heap = "1G"

		deploymentInfo.ishtar.numReplicas = int32(3)
		deploymentInfo.ishtar.cpuRequests = "1"
		deploymentInfo.ishtar.memoryRequests = "1024Mi"
		deploymentInfo.ishtar.heap = "1G"

		deploymentInfo.relay.numReplicas = int32(3)
		deploymentInfo.relay.cpuRequests = "1"
		deploymentInfo.relay.memoryRequests = "1024Mi"
		deploymentInfo.relay.heap = "1G"

		deploymentInfo.watchtower.numReplicas = int32(3)
		deploymentInfo.watchtower.cpuRequests = "1"
		deploymentInfo.watchtower.memoryRequests = "1024Mi"
		deploymentInfo.watchtower.heap = "1G"

		deploymentInfo.doki.numReplicas = int32(3)
		deploymentInfo.doki.cpuRequests = "1"
		deploymentInfo.doki.memoryRequests = "1024Mi"
		deploymentInfo.doki.heap = "1G"

		deploymentInfo.brent.numReplicas = int32(3)
		deploymentInfo.brent.cpuRequests = "1"
		deploymentInfo.brent.memoryRequests = "1024Mi"
		deploymentInfo.brent.heap = "1G"
	} else if strings.ToLower(instance.Spec.DeploymentType) == "tiny" {
		deploymentInfo.conductor.numReplicas = int32(1)
		deploymentInfo.conductor.cpuRequests = "100m"
		deploymentInfo.conductor.memoryRequests = "156Mi"
		deploymentInfo.conductor.heap = "128m"

		deploymentInfo.apollo.numReplicas = int32(1)
		deploymentInfo.apollo.cpuRequests = "100m"
		deploymentInfo.apollo.memoryRequests = "156Mi"
		deploymentInfo.apollo.heap = "128m"

		deploymentInfo.galileo.numReplicas = int32(1)
		deploymentInfo.galileo.cpuRequests = "300m"
		deploymentInfo.galileo.memoryRequests = "300Mi"
		deploymentInfo.galileo.heap = "256m"

		deploymentInfo.talledega.numReplicas = int32(1)
		deploymentInfo.talledega.cpuRequests = "300m"
		deploymentInfo.talledega.memoryRequests = "300Mi"
		deploymentInfo.talledega.heap = "256m"

		deploymentInfo.daytona.numReplicas = int32(1)
		deploymentInfo.daytona.cpuRequests = "500m"
		deploymentInfo.daytona.memoryRequests = "156Mi"
		deploymentInfo.daytona.heap = "128m"

		deploymentInfo.nimrod.numReplicas = int32(1)
		deploymentInfo.nimrod.cpuRequests = "300m"
		deploymentInfo.nimrod.memoryRequests = "300Mi"
		deploymentInfo.nimrod.heap = "256m"

		deploymentInfo.ishtar.numReplicas = int32(1)
		deploymentInfo.ishtar.cpuRequests = "200m"
		deploymentInfo.ishtar.memoryRequests = "156Mi"
		deploymentInfo.ishtar.heap = "128m"

		deploymentInfo.relay.numReplicas = int32(1)
		deploymentInfo.relay.cpuRequests = "50m"
		deploymentInfo.relay.memoryRequests = "156Mi"
		deploymentInfo.relay.heap = "128m"

		deploymentInfo.watchtower.numReplicas = int32(1)
		deploymentInfo.watchtower.cpuRequests = "500m"
		deploymentInfo.watchtower.memoryRequests = "412Mi"
		deploymentInfo.watchtower.heap = "384m"

		deploymentInfo.doki.numReplicas = int32(1)
		deploymentInfo.doki.cpuRequests = "200m"
		deploymentInfo.doki.memoryRequests = "384Mi"
		deploymentInfo.doki.heap = "256m"

		deploymentInfo.brent.numReplicas = int32(3)
		deploymentInfo.brent.cpuRequests = "200m"
		deploymentInfo.brent.memoryRequests = "256Mi"
		deploymentInfo.brent.heap = "256m"
	} else if strings.ToLower(instance.Spec.DeploymentType) == "middle" {
		deploymentInfo.conductor.numReplicas = int32(1)
		deploymentInfo.conductor.cpuRequests = "100m"
		deploymentInfo.conductor.memoryRequests = "128Mi"
		deploymentInfo.conductor.heap = "128m"

		deploymentInfo.apollo.numReplicas = int32(1)
		deploymentInfo.apollo.cpuRequests = "100m"
		deploymentInfo.apollo.memoryRequests = "512Mi"
		deploymentInfo.apollo.heap = "512m"

		deploymentInfo.galileo.numReplicas = int32(1)
		deploymentInfo.galileo.cpuRequests = "500m"
		deploymentInfo.galileo.memoryRequests = "512Mi"
		deploymentInfo.galileo.heap = "512m"

		deploymentInfo.talledega.numReplicas = int32(1)
		deploymentInfo.talledega.cpuRequests = "500m"
		deploymentInfo.talledega.memoryRequests = "512Mi"
		deploymentInfo.talledega.heap = "512m"

		deploymentInfo.daytona.numReplicas = int32(1)
		deploymentInfo.daytona.cpuRequests = "500m"
		deploymentInfo.daytona.memoryRequests = "512Mi"
		deploymentInfo.daytona.heap = "512m"

		deploymentInfo.nimrod.numReplicas = int32(1)
		deploymentInfo.nimrod.cpuRequests = "300m"
		deploymentInfo.nimrod.memoryRequests = "256Mi"
		deploymentInfo.nimrod.heap = "256m"

		deploymentInfo.ishtar.numReplicas = int32(1)
		deploymentInfo.ishtar.cpuRequests = "300m"
		deploymentInfo.ishtar.memoryRequests = "256Mi"
		deploymentInfo.ishtar.heap = "256m"

		deploymentInfo.relay.numReplicas = int32(1)
		deploymentInfo.relay.cpuRequests = "100m"
		deploymentInfo.relay.memoryRequests = "128Mi"
		deploymentInfo.relay.heap = "128m"

		deploymentInfo.watchtower.numReplicas = int32(1)
		deploymentInfo.watchtower.cpuRequests = "500m"
		deploymentInfo.watchtower.memoryRequests = "512Mi"
		deploymentInfo.watchtower.heap = "512m"

		deploymentInfo.doki.numReplicas = int32(1)
		deploymentInfo.doki.cpuRequests = "250m"
		deploymentInfo.doki.memoryRequests = "512Mi"
		deploymentInfo.doki.heap = "512m"

		deploymentInfo.brent.numReplicas = int32(1)
		deploymentInfo.brent.cpuRequests = "500m"
		deploymentInfo.brent.memoryRequests = "512Mi"
		deploymentInfo.brent.heap = "512m"
	} else {
		deploymentInfo.conductor.numReplicas = int32(1)
		deploymentInfo.conductor.cpuRequests = "1"
		deploymentInfo.conductor.memoryRequests = "1024Mi"
		deploymentInfo.conductor.heap = "1G"

		deploymentInfo.apollo.numReplicas = int32(1)
		deploymentInfo.apollo.cpuRequests = "1"
		deploymentInfo.apollo.memoryRequests = "1024Mi"
		deploymentInfo.apollo.heap = "1G"

		deploymentInfo.galileo.numReplicas = int32(1)
		deploymentInfo.galileo.cpuRequests = "1"
		deploymentInfo.galileo.memoryRequests = "1024Mi"
		deploymentInfo.galileo.heap = "1G"

		deploymentInfo.talledega.numReplicas = int32(1)
		deploymentInfo.talledega.cpuRequests = "1"
		deploymentInfo.talledega.memoryRequests = "1024Mi"
		deploymentInfo.talledega.heap = "1G"

		deploymentInfo.daytona.numReplicas = int32(1)
		deploymentInfo.daytona.cpuRequests = "1"
		deploymentInfo.daytona.memoryRequests = "1024Mi"
		deploymentInfo.daytona.heap = "1G"

		deploymentInfo.nimrod.numReplicas = int32(1)
		deploymentInfo.nimrod.cpuRequests = "1"
		deploymentInfo.nimrod.memoryRequests = "1024Mi"
		deploymentInfo.nimrod.heap = "1G"

		deploymentInfo.ishtar.numReplicas = int32(1)
		deploymentInfo.ishtar.cpuRequests = "1"
		deploymentInfo.ishtar.memoryRequests = "1024Mi"
		deploymentInfo.ishtar.heap = "1G"

		deploymentInfo.relay.numReplicas = int32(1)
		deploymentInfo.relay.cpuRequests = "1"
		deploymentInfo.relay.memoryRequests = "1024Mi"
		deploymentInfo.relay.heap = "1G"

		deploymentInfo.watchtower.numReplicas = int32(1)
		deploymentInfo.watchtower.cpuRequests = "1"
		deploymentInfo.watchtower.memoryRequests = "1024Mi"
		deploymentInfo.watchtower.heap = "1G"

		deploymentInfo.doki.numReplicas = int32(1)
		deploymentInfo.doki.cpuRequests = "1"
		deploymentInfo.doki.memoryRequests = "1024Mi"
		deploymentInfo.doki.heap = "1G"

		deploymentInfo.brent.numReplicas = int32(1)
		deploymentInfo.brent.cpuRequests = "1"
		deploymentInfo.brent.memoryRequests = "1024Mi"
		deploymentInfo.brent.heap = "1G"
	}

	return deploymentInfo, nil
}
