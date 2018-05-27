package scheduler

import (
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/hortonworks/gohadoop/hadoop_common/security"
	"github.com/hortonworks/gohadoop/hadoop_yarn"
	"github.com/hortonworks/gohadoop/hadoop_yarn/conf"
	"github.com/hortonworks/gohadoop/hadoop_yarn/yarn_client"
)

const (
	allocateIntervalMs       = 1000 //1 second - keep this low to avoid excessive allocation delays
	maxCallbackNotifications = 16
)

type YARNScheduler struct {
	yarnClient          *yarn_client.YarnClient
	rmClient            *yarn_client.AMRMClientAsync
	podsToContainersMap map[string]*hadoop_yarn.ContainerIdProto
	handler             *yarnSchedulerCallbackHandler
}

type yarnSchedulerCallbackHandler struct {
	completedContainers chan []*hadoop_yarn.ContainerStatusProto
	allocatedContainers chan []*hadoop_yarn.ContainerProto
}

func NewYARNScheduler() Scheduler {
	handler := newYarnSchedulerCallbackHandler()
	yarnC, rmC := YARNInit(handler)
	podsToContainers := make(map[string]*hadoop_yarn.ContainerIdProto)

	return &YARNScheduler{
		yarnClient:          yarnC,
		rmClient:            rmC,
		podsToContainersMap: podsToContainers,
		handler:             handler}
}

func YARNInit(handler *yarnSchedulerCallbackHandler) (*yarn_client.YarnClient, *yarn_client.AMRMClientAsync) {
	var err error

	hadoopConfDir := os.Getenv("HADOOP_CONF_DIR")

	if hadoopConfDir == "" {
		log.Fatal("HADOOP_CONF_DIR not set! Unable to initialize YARNScheduler!")
	}

	// Create YarnConfiguration
	conf, _ := conf.NewYarnConfiguration()

	// Create YarnClient
	yarnClient, _ := yarn_client.CreateYarnClient(conf)

	// Create new application to get ApplicationSubmissionContext
	_, asc, _ := yarnClient.CreateNewApplication()

	// Some useful information
	queue := "default"
	appName := "kubernetes"
	appType := "PAAS"
	unmanaged := true
	clc := hadoop_yarn.ContainerLaunchContextProto{}

	// Setup ApplicationSubmissionContext for the application
	asc.AmContainerSpec = &clc
	asc.ApplicationName = &appName
	asc.Queue = &queue
	asc.ApplicationType = &appType
	asc.UnmanagedAm = &unmanaged

	// Submit!
	err = yarnClient.SubmitApplication(asc)
	if err != nil {
		log.Fatal("yarnClient.SubmitApplication ", err)
	}
	log.Println("Successfully submitted unmanaged application: ", asc.ApplicationId)
	time.Sleep(1 * time.Second)

	appReport, err := yarnClient.GetApplicationReport(asc.ApplicationId)
	if err != nil {
		log.Fatal("yarnClient.GetApplicationReport ", err)
	}
	appState := appReport.GetYarnApplicationState()
	for appState != hadoop_yarn.YarnApplicationStateProto_ACCEPTED {
		log.Println("Application in state ", appState)
		time.Sleep(1 * time.Second)
		appReport, err = yarnClient.GetApplicationReport(asc.ApplicationId)
		appState = appReport.GetYarnApplicationState()
		if appState == hadoop_yarn.YarnApplicationStateProto_FAILED || appState == hadoop_yarn.YarnApplicationStateProto_KILLED {
			log.Fatal("Application in state ", appState)
		}
	}

	amRmToken := appReport.GetAmRmToken()
	for amRmToken == nil {
		log.Println("AMRMToken is nil. sleeping before trying again.")
		time.Sleep(1 * time.Second)
		appReport, err = yarnClient.GetApplicationReport(asc.ApplicationId)
		if err != nil {
			log.Println("Failed to get application report! error: ", err)
			return nil, nil
		}
		amRmToken = appReport.GetAmRmToken()
	}

	if amRmToken != nil {
		savedAmRmToken := *amRmToken
		service, _ := conf.GetRMSchedulerAddress()
		log.Println("Saving token with address: ", service)
		security.GetCurrentUser().AddUserTokenWithAlias(service, &savedAmRmToken)
	}

	log.Println("Application in state ", appState)

	// Create AMRMClient
	var attemptId int32
	attemptId = 1
	applicationAttemptId := hadoop_yarn.ApplicationAttemptIdProto{ApplicationId: asc.ApplicationId, AttemptId: &attemptId}

	rmClient, _ := yarn_client.CreateAMRMClientAsync(conf, allocateIntervalMs, *handler)

	log.Println("Created RM client: ", rmClient)

	// Wait for ApplicationAttempt to be in Launched state
	appAttemptReport, err := yarnClient.GetApplicationAttemptReport(&applicationAttemptId)
	appAttemptState := appAttemptReport.GetYarnApplicationAttemptState()
	for appAttemptState != hadoop_yarn.YarnApplicationAttemptStateProto_APP_ATTEMPT_LAUNCHED {
		log.Println("ApplicationAttempt in state ", appAttemptState)
		time.Sleep(1 * time.Second)
		appAttemptReport, err = yarnClient.GetApplicationAttemptReport(&applicationAttemptId)
		appAttemptState = appAttemptReport.GetYarnApplicationAttemptState()
	}
	log.Println("ApplicationAttempt in state ", appAttemptState)

	// Register with ResourceManager
	log.Println("About to register application master.")
	err = rmClient.RegisterApplicationMaster("", -1, "")
	if err != nil {
		log.Fatal("rmClient.RegisterApplicationMaster ", err)
	}
	log.Println("Successfully registered application master.")

	return yarnClient, rmClient
}

func (yarnScheduler *YARNScheduler) Delete(id string) error {
	log.Printf("attempting to delete pod with id: %s", id)
	rmClient := yarnScheduler.rmClient
	containerId, found := yarnScheduler.podsToContainersMap[id]

	if !found {
		log.Println("attempting to delete a pod that doesn't have an associated YARN container!")
		return errors.New("attempting to delete a pod that doesn't have an associated YARN container!")
	}

	rmClient.ReleaseAssignedContainer(containerId)

	const maxAttempts = int(5)
	releaseAttempts := 0

	for releaseAttempts < maxAttempts {
		select {
		case completedContainers := <-yarnScheduler.handler.completedContainers:
			log.Println("received container completion status: ", completedContainers[0])
			return nil
		case <-time.After(3 * time.Second):
			// Sleep for a while before trying again
			releaseAttempts++
			log.Println("no response to release request. sleeping...")
			time.Sleep(3 * time.Second)
			continue
		}
	}

	log.Println("failed to delete container with id: ", id)
	return errors.New("failed to release container!")
}

//allocates one "fat" container.
func (yarnScheduler *YARNScheduler) Schedule(pod api.Pod, minionLister MinionLister) (string, error) {
	rmClient := yarnScheduler.rmClient

	// Add resource requests
	// Even though we only allocate one container for now, this could change in the future.

	const numContainers = int32(1)
	memory := int32(128)
	resource := hadoop_yarn.ResourceProto{Memory: &memory}
	numAllocatedContainers := int32(0)
	const maxAttempts = int(5)
	allocationAttempts := 0
	//allocatedContainers := make([]*hadoop_yarn.ContainerProto, numContainers, numContainers)

	rmClient.AddRequest(1, "*", &resource, numContainers)

	for numAllocatedContainers < numContainers && allocationAttempts < maxAttempts {
		var allocatedContainers []*hadoop_yarn.ContainerProto

		select {
		case allocatedContainers = <-yarnScheduler.handler.allocatedContainers:
			break
		case <-time.After(3 * time.Second):
			// Sleep for a while before trying again
			allocationAttempts++
			log.Println("Sleeping...")
			time.Sleep(3 * time.Second)
			log.Println("Sleeping... done!")
			continue
		}

		for _, container := range allocatedContainers {
			allocatedContainers[numAllocatedContainers] = container
			numAllocatedContainers++
			log.Println("#containers allocated so far: ", numAllocatedContainers)

			log.Printf("pod=%s YARN container=%v", pod.Name, container.GetId())
			yarnScheduler.podsToContainersMap[pod.Name] = container.GetId()
			host := *container.NodeId.Host
			port := *container.NodeId.Port
			log.Println("allocated container on: ", host)

			//launch a "sleep" container so that the "loop" is completed and YARN doesn't deallocate the container
			containerLaunchContext := hadoop_yarn.ContainerLaunchContextProto{Command: []string{"while :; do sleep 1; done"}}
			nmClient, err := yarn_client.CreateAMNMClient(host, int(port))
			if err != nil {
				log.Fatal("Failed to create AMNMClient! ", err)
			}
			log.Printf("Successfully created nmClient: %v", nmClient)
			log.Printf("Attempting to start container on %s", host)

			err = nmClient.StartContainer(container, &containerLaunchContext)
			if err != nil {
				log.Fatal("failed to start container! ", err)
			}

			//We have the hostname available. return from here.
			//This many change in case we allocate more than one container in the future.
			return findMinionForHost(host, minionLister)
		}

		log.Println("#containers allocated: ", len(allocatedContainers))
		log.Println("Total #containers allocated so far: ", numAllocatedContainers)
	}

	log.Println("Final #containers allocated: ", numAllocatedContainers)

	return "<invalid_host>", errors.New("unable to schedule pod! YARN didn't allocate a container")
}

/* YARN returns hostnames, but minions maybe using IPs.
TODO: This is an expensive mechanism to find the right minion corresponding to the YARN node.
Find a better mechanism if possible (at a minimum - caching could be added in some form)
*/
func findMinionForHost(host string, minionLister MinionLister) (string, error) {
	hostIPs, err := net.LookupIP(host)

	if err != nil {
		return "<invalid_host>", errors.New("unable to lookup IPs for YARN host: " + host)
	}

	for _, hostIP := range hostIPs {
		minions, err := minionLister.List()
		if err != nil {
			return "<invalid_host>", errors.New("update to list minions")
		}

		for _, minion := range minions.Items {
			minionStr := minion.Name
			minionIPs, err := net.LookupIP(minionStr)

			if err != nil {
				return "<invalid_host>", errors.New("unable to lookup IPs for minion: " + minionStr)
			}

			for _, minionIP := range minionIPs {
				if hostIP.Equal(minionIP) {
					log.Printf("YARN node %s maps to minion: %s", host, minionStr)
					return minionStr, nil
				}
			}
		}
	}

	return "<invalid_host>", errors.New("unable to find minion for YARN host: " + host)
}

func newYarnSchedulerCallbackHandler() *yarnSchedulerCallbackHandler {
	completedContainers := make(chan []*hadoop_yarn.ContainerStatusProto, maxCallbackNotifications)
	allocatedContainers := make(chan []*hadoop_yarn.ContainerProto, maxCallbackNotifications)

	return &yarnSchedulerCallbackHandler{completedContainers: completedContainers, allocatedContainers: allocatedContainers}
}

func (ch yarnSchedulerCallbackHandler) OnContainersCompleted(completedContainers []*hadoop_yarn.ContainerStatusProto) {
	log.Println("received completed containers notification. writing to channel: ", completedContainers)
	ch.completedContainers <- completedContainers
}

func (ch yarnSchedulerCallbackHandler) OnContainersAllocated(allocatedContainers []*hadoop_yarn.ContainerProto) {
	log.Println("received allocated containers notification. writing to channel: ", allocatedContainers)
	ch.allocatedContainers <- allocatedContainers
}

func (ch yarnSchedulerCallbackHandler) OnShutDownRequest() {
	log.Println("received shutdown request!")
}

func (ch yarnSchedulerCallbackHandler) OnNodesUpdated(updatedNodes []*hadoop_yarn.NodeReportProto) {
	log.Println("nodes updated. this operation is not suppored. doing nothing.")
}

//not supported currently
func (ch yarnSchedulerCallbackHandler) GetProgress() float64 {
	return 0.001
}

func (ch yarnSchedulerCallbackHandler) OnError(err error) {
	log.Println("allocation error! error: ", err)
}
