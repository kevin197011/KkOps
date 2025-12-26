import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider } from './contexts/AuthContext';
import PrivateRoute from './components/PrivateRoute';
import MainLayout from './components/MainLayout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Projects from './pages/Projects';
import Environments from './pages/Environments';
import CloudPlatforms from './pages/CloudPlatforms';
import Hosts from './pages/Hosts';
import Users from './pages/Users';
import Roles from './pages/Roles';
import Permissions from './pages/Permissions';
import HostGroups from './pages/HostGroups';
import HostTags from './pages/HostTags';
import Deployments from './pages/Deployments';
import Audit from './pages/Audit';
import WebSSH from './pages/WebSSH';
import SSHKeys from './pages/SSHKeys';
import Settings from './pages/Settings';
import BatchOperations from './pages/BatchOperations';
import FormulaDeployment from './pages/FormulaDeployment';
import './App.css';

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <AuthProvider>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route
              path="/"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Dashboard />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/projects"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Projects />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/environments"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Environments />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/cloud-platforms"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <CloudPlatforms />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/hosts"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Hosts />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/webssh"
              element={
                <PrivateRoute>
                  <WebSSH />
                </PrivateRoute>
              }
            />
            <Route
              path="/ssh-keys"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <SSHKeys />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/deployments"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Deployments />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/users"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Users />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/roles"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Roles />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/permissions"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Permissions />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/host-groups"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <HostGroups />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/host-tags"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <HostTags />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/audit"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Audit />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/settings"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <Settings />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/batch-operations"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <BatchOperations />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route
              path="/formulas"
              element={
                <PrivateRoute>
                  <MainLayout>
                    <FormulaDeployment />
                  </MainLayout>
                </PrivateRoute>
              }
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Router>
      </AuthProvider>
    </ConfigProvider>
  );
};

export default App;
