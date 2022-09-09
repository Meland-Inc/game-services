package controller

type ServiceStatus string

const (
	ServiceStatusFree   ServiceStatus = "Free"
	ServiceStatusNormal ServiceStatus = "Normal"
	ServiceStatusBusy   ServiceStatus = "Busy"
)

type ServiceStatusRecord struct {
	recordByService map[int64]ServiceStatus
	recordByStatus  map[ServiceStatus][]int64
}

func NewServiceStatusRecord() *ServiceStatusRecord {
	return &ServiceStatusRecord{
		recordByService: make(map[int64]ServiceStatus),
		recordByStatus:  make(map[ServiceStatus][]int64),
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

func (s *ServiceStatusRecord) GetServicesByStatus(status ServiceStatus) []int64 {
	return s.recordByStatus[status]
}

func (s *ServiceStatusRecord) AddServiceStatusRecord(serviceId int64, online, maxOnline int32) {
	status := s.calculateServiceStatus(online, maxOnline)
	curStatus, exist := s.recordByService[serviceId]
	if exist && curStatus == status {
		return
	}

	if exist {
		s.RemoveServiceStatusRecord(serviceId)
	}
	s.recordByService[serviceId] = status
	if _, exist := s.recordByStatus[status]; !exist {
		s.recordByStatus[status] = []int64{serviceId}
	} else {
		s.recordByStatus[status] = append(s.recordByStatus[status], serviceId)
	}
}

func (s *ServiceStatusRecord) RemoveServiceStatusRecord(serviceId int64) {
	status, exist := s.recordByService[serviceId]
	if !exist {
		return
	}

	idList := s.recordByStatus[status]
	for idx, id := range idList {
		if id == serviceId {
			idList = append(idList[:idx], idList[idx+1:]...)
		}
	}
	s.recordByStatus[status] = idList
	delete(s.recordByService, serviceId)
}
