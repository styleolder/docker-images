package factory

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/cache"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/scheduler"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	"github.com/golang/glog"
)

//SchedulerCache "wraps" a cache store and notifies the scheduler of
//certain caching events. For now, only the "delete" notification needs to be relayed
type SchedulerCache struct {
	store     cache.Store
	scheduler scheduler.Scheduler
}

func NewSchedulerCache(store cache.Store, scheduler scheduler.Scheduler) *SchedulerCache {
	return &SchedulerCache{store: store, scheduler: scheduler}
}

//Unfortunate Hack - so that scheduler can be set "later"
func NewSchedulerCacheUnsetScheduler(store cache.Store) *SchedulerCache {
	return &SchedulerCache{store: store}
}

func (c *SchedulerCache) SetScheduler(scheduler scheduler.Scheduler) {
	c.scheduler = scheduler
}

func (c *SchedulerCache) Add(id string, obj interface{}) {
	c.store.Add(id, obj)
}

func (c *SchedulerCache) Update(id string, obj interface{}) {
	c.store.Update(id, obj)
}

func (c *SchedulerCache) Delete(id string) {
	c.store.Delete(id)
	glog.V(0).Infof("the following pod has been deleted from cache: %s. notifying scheduler", id)

	if statefulScheduler, ok := c.scheduler.(scheduler.StatefulScheduler); ok {
		statefulScheduler.Delete(id)
	} else {
		glog.V(0).Infof("scheduler does not support deletes. this will likely lead to container leaks.")
	}
}

func (c *SchedulerCache) List() []interface{} {
	return c.store.List()
}

func (c *SchedulerCache) ContainedIDs() util.StringSet {
	return c.store.ContainedIDs()
}

func (c *SchedulerCache) Get(id string) (item interface{}, exists bool) {
	return c.store.Get(id)
}

func (c *SchedulerCache) Replace(idToObj map[string]interface{}) {
	c.store.Replace(idToObj)
}
