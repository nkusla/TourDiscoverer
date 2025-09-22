import api from './api'

export const tourExecutionService = {
  // Get tours that can be executed by tourists
  getExecutableTours() {
    return api.get('/api/tour/executable')
  },

  // Start a tour execution
  startTourExecution(tourId, latitude, longitude) {
    return api.post('/api/tour/execution/start', {
      tour_id: tourId,
      latitude: latitude,
      longitude: longitude
    })
  },

  // Get active tour execution for the current tourist
  getActiveTourExecution() {
    return api.get('/api/tour/execution/active')
  },

  // End tour execution
  endTourExecution(executionId, status) {
    return api.put(`/api/tour/execution/${executionId}/end`, {
      status: status // 'completed' or 'abandoned'
    })
  },

  // Check proximity to key points
  checkProximity(executionId, latitude, longitude) {
    return api.post(`/api/tour/execution/${executionId}/check-proximity`, {
      latitude: latitude,
      longitude: longitude
    })
  },

  // Get tourist position from stakeholder service
  getTouristPosition() {
    return api.get('/api/stakeholder/position')
  },

  // Get specific tour details
  getTourById(tourId) {
    return api.get(`/api/tour/${tourId}`)
  }
}