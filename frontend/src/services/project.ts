import api from '../config/api';

export interface Project {
  id: number;
  name: string;
  description?: string;
  status: 'active' | 'disabled';
  created_by: number;
  created_at: string;
  updated_at: string;
  creator?: {
    id: number;
    username: string;
    display_name?: string;
  };
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  status?: string;
}

export interface ListProjectsResponse {
  projects: Project[];
  total: number;
  page: number;
  page_size: number;
}

export interface ProjectMember {
  user_id: number;
  role: 'owner' | 'admin' | 'member' | 'viewer';
  user?: {
    id: number;
    username: string;
    display_name?: string;
  };
}

export const projectService = {
  list: async (page: number = 1, pageSize: number = 10, filters?: any): Promise<ListProjectsResponse> => {
    const response = await api.get('/projects', {
      params: { page, page_size: pageSize, ...filters },
    });
    return response.data;
  },

  get: async (id: number): Promise<{ project: Project }> => {
    const response = await api.get(`/projects/${id}`);
    return response.data;
  },

  create: async (data: CreateProjectRequest): Promise<{ project: Project }> => {
    const response = await api.post('/projects', data);
    return response.data;
  },

  update: async (id: number, data: UpdateProjectRequest): Promise<{ project: Project }> => {
    const response = await api.put(`/projects/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/projects/${id}`);
  },

  addMember: async (projectId: number, userId: number, role: string): Promise<void> => {
    await api.post(`/projects/${projectId}/members`, {
      user_id: userId,
      role: role,
    });
  },

  removeMember: async (projectId: number, userId: number): Promise<void> => {
    await api.delete(`/projects/${projectId}/members`, {
      data: { user_id: userId },
    });
  },

  getMembers: async (projectId: number): Promise<{ members: ProjectMember[] }> => {
    const response = await api.get(`/projects/${projectId}/members`);
    return response.data;
  },
};

