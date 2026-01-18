// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { lazy, Suspense } from 'react'
import { Spin } from 'antd'
import Login from './pages/auth/Login'
import ProtectedRoute from './components/ProtectedRoute'
import MainLayout from './layouts/MainLayout'

// Lazy load pages for code splitting
const Dashboard = lazy(() => import('./pages/Dashboard'))
const UserList = lazy(() => import('./pages/users/UserList'))
const RoleList = lazy(() => import('./pages/roles/RoleList'))
const AssetList = lazy(() => import('./pages/assets/AssetList'))
const AssetDetail = lazy(() => import('./pages/assets/AssetDetail'))
const ProjectList = lazy(() => import('./pages/projects/ProjectList'))
const EnvironmentList = lazy(() => import('./pages/environments/EnvironmentList'))
const CloudPlatformList = lazy(() => import('./pages/cloudPlatforms/CloudPlatformList'))
const CategoryList = lazy(() => import('./pages/categories/CategoryList'))
const TagList = lazy(() => import('./pages/tags/TagList'))
const TemplateList = lazy(() => import('./pages/executions/TemplateList'))
const ExecutionOperatorPage = lazy(() => import('./pages/executions/ExecutionOperatorPage'))
const ExecutionHistoryPage = lazy(() => import('./pages/executions/TaskExecutionList'))
const ExecutionLogsPage = lazy(() => import('./pages/executions/TaskExecutionLogs'))
const DeploymentModuleList = lazy(() => import('./pages/deployments/DeploymentModuleList'))
const ScheduledTaskList = lazy(() => import('./pages/tasks/ScheduledTaskList'))
const SSHKeyList = lazy(() => import('./pages/ssh/SSHKeyList'))
const WebSSHTerminal = lazy(() => import('./pages/ssh/WebSSHTerminal'))
const AuditLogList = lazy(() => import('./pages/audit/AuditLogList'))

// Loading component
const PageLoading = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' }}>
    <Spin size="large" />
  </div>
)

function App() {
  return (
    <BrowserRouter
      future={{
        v7_startTransition: true,
        v7_relativeSplatPath: true,
      }}
    >
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <Dashboard />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/users"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <UserList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/roles"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <RoleList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/assets"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <AssetList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/assets/:id"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <AssetDetail />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <ProjectList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/environments"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <EnvironmentList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/cloud-platforms"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <CloudPlatformList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/categories"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <CategoryList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/tags"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <TagList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/templates"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <TemplateList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/executions"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <ExecutionOperatorPage />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/executions/:executionId/history"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <ExecutionHistoryPage />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/execution-records/:id/logs"
          element={
            <ProtectedRoute>
              <Suspense fallback={<PageLoading />}>
                <ExecutionLogsPage />
              </Suspense>
            </ProtectedRoute>
          }
        />
        <Route
          path="/deployments"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <DeploymentModuleList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/tasks"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <ScheduledTaskList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/ssh/keys"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <SSHKeyList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/ssh/terminal"
          element={
            <ProtectedRoute>
              <Suspense fallback={<PageLoading />}>
                <WebSSHTerminal />
              </Suspense>
            </ProtectedRoute>
          }
        />
        <Route
          path="/audit-logs"
          element={
            <ProtectedRoute>
              <MainLayout>
                <Suspense fallback={<PageLoading />}>
                  <AuditLogList />
                </Suspense>
              </MainLayout>
            </ProtectedRoute>
          }
        />
        <Route path="/" element={<Navigate to="/login" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
