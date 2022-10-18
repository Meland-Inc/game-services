package controller

type ServiceStatus string

const (
	ServiceStatusFree   ServiceStatus = "Free"
	ServiceStatusNormal ServiceStatus = "Normal"
	ServiceStatusBusy   ServiceStatus = "Busy"
)

type ServiceStatusRecord struct {
	recordByService map[string]ServiceStatus
	recordByStatus  map[ServiceStatus][]string
}

func NewServiceStatusRecord() *ServiceStatusRecord {
	return &ServiceStatusRecord{
		recordByService: make(map[string]ServiceStatus),
		recordByStatus:  make(map[ServiceStatus][]string),
	}
}

func (s *ServiceStatusRecord) calculateServiceStatus(online, maxOnline int32) ServiceStatus {
	if online == 0 || maxOnline == 0 {
		return ServiceStatusFree
	}

	percent := int32(float32(online) / float32(maxOnline) * 100)
	status := ServiceStatusNormal
	if percent <= 30 {
		status = ServiceStatusFree
	} else if percent >= 80 {
		status = ServiceStatusBusy
	}
	return status
}

func (s *ServiceStatusRecord) GetServicesByStatus(status ServiceStatus) []string {
	return s.recordByStatus[status]
}

func (s *ServiceStatusRecord) AddServiceStatusRecord(appId string, online, maxOnline int32) {
	status := s.calculateServiceStatus(online, maxOnline)
	curStatus, exist := s.recordByService[appId]
	if exist && curStatus == status {
		return
	}

	if exist {
		s.RemoveServiceStatusRecord(appId)
	}
	s.recordByService[appId] = status
	if _, exist := s.recordByStatus[status]; !exist {
		s.recordByStatus[status] = []string{appId}
	} else {
		s.recordByStatus[status] = append(s.recordByStatus[status], appId)
	}
}

func (s *ServiceStatusRecord) RemoveServiceStatusRecord(appId string) {
	status, exist := s.recordByService[appId]
	if !exist {
		return
	}

	idList := s.recordByStatus[status]
	for idx, id := range idList {
		if id == appId {
			idList = append(idList[:idx], idList[idx+1:]...)
		}
	}
	s.recordByStatus[status] = idList
	delete(s.recordByService, appId)
}
