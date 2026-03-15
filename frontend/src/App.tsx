import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { ThemeProvider } from '@/contexts/ThemeContext'
import { AuthProvider } from '@/contexts/AuthContext'
import ProtectedRoute from '@/components/layout/ProtectedRoute'
import AppLayout from '@/components/layout/AppLayout'
import LoginPage from '@/pages/LoginPage'
import DashboardPage from '@/pages/DashboardPage'
import UsersPage from '@/pages/UsersPage'
import ServicesPage from '@/pages/ServicesPage'
import PermissionsPage from '@/pages/PermissionsPage'
import RolesPage from '@/pages/RolesPage'
import AuditPage from '@/pages/AuditPage'
import ConsumerLoginPage from '@/pages/auth/ConsumerLoginPage'
import TokenCallbackPage from '@/pages/auth/TokenCallbackPage'

function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <AuthProvider>
          <Routes>
            {/* Consumer-facing login (OTP + Google OAuth) — used by all IIT products */}
            <Route path="/auth/login" element={<ConsumerLoginPage />} />
            <Route path="/auth/callback" element={<TokenCallbackPage />} />

            {/* Admin login */}
            <Route path="/login" element={<LoginPage />} />

            {/* Protected */}
            <Route
              element={
                <ProtectedRoute>
                  <AppLayout />
                </ProtectedRoute>
              }
            >
              <Route path="/" element={<DashboardPage />} />
              <Route
                path="/users"
                element={
                  <ProtectedRoute requiredRole="admin">
                    <UsersPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/services"
                element={
                  <ProtectedRoute requiredRole="admin">
                    <ServicesPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/permissions"
                element={
                  <ProtectedRoute requiredRole="admin">
                    <PermissionsPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/roles"
                element={
                  <ProtectedRoute requiredRole="admin">
                    <RolesPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/audit"
                element={
                  <ProtectedRoute requiredRole="admin">
                    <AuditPage />
                  </ProtectedRoute>
                }
              />
            </Route>
          </Routes>
          <Toaster
            position="top-right"
            toastOptions={{
              className: 'dark:bg-gray-800 dark:text-gray-100',
              duration: 4000,
            }}
          />
        </AuthProvider>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
