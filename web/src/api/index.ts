import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
});

// 订阅 API
export const subscriptionApi = {
  getAll: () => api.get('/subscriptions'),
  add: (name: string, url: string) => api.post('/subscriptions', { name, url }),
  update: (id: string, data: any) => api.put(`/subscriptions/${id}`, data),
  delete: (id: string) => api.delete(`/subscriptions/${id}`),
  refresh: (id: string) => api.post(`/subscriptions/${id}/refresh`),
  refreshAll: () => api.post('/subscriptions/refresh-all'),
};

// 过滤器 API
export const filterApi = {
  getAll: () => api.get('/filters'),
  add: (data: any) => api.post('/filters', data),
  update: (id: string, data: any) => api.put(`/filters/${id}`, data),
  delete: (id: string) => api.delete(`/filters/${id}`),
};

// 规则 API
export const ruleApi = {
  getAll: () => api.get('/rules'),
  add: (data: any) => api.post('/rules', data),
  update: (id: string, data: any) => api.put(`/rules/${id}`, data),
  delete: (id: string) => api.delete(`/rules/${id}`),
};

// 规则组 API
export const ruleGroupApi = {
  getAll: () => api.get('/rule-groups'),
  update: (id: string, data: any) => api.put(`/rule-groups/${id}`, data),
};

// 设置 API
export const settingsApi = {
  get: () => api.get('/settings'),
  update: (data: any) => api.put('/settings', data),
};

// 配置 API
export const configApi = {
  generate: () => api.post('/config/generate'),
  preview: () => api.get('/config/preview'),
  apply: () => api.post('/config/apply'),
};

// 服务 API
export const serviceApi = {
  status: () => api.get('/service/status'),
  start: () => api.post('/service/start'),
  stop: () => api.post('/service/stop'),
  restart: () => api.post('/service/restart'),
  reload: () => api.post('/service/reload'),
};

// launchd API
export const launchdApi = {
  status: () => api.get('/launchd/status'),
  install: () => api.post('/launchd/install'),
  uninstall: () => api.post('/launchd/uninstall'),
};

// 监控 API
export const monitorApi = {
  system: () => api.get('/monitor/system'),
  logs: () => api.get('/monitor/logs'),
};

// 节点 API
export const nodeApi = {
  getAll: () => api.get('/nodes'),
  getCountries: () => api.get('/nodes/countries'),
  getByCountry: (code: string) => api.get(`/nodes/country/${code}`),
  parse: (url: string) => api.post('/nodes/parse', { url }),
};

// 手动节点 API
export const manualNodeApi = {
  getAll: () => api.get('/manual-nodes'),
  add: (data: any) => api.post('/manual-nodes', data),
  update: (id: string, data: any) => api.put(`/manual-nodes/${id}`, data),
  delete: (id: string) => api.delete(`/manual-nodes/${id}`),
};

export default api;
