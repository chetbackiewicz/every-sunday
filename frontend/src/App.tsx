import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AppProvider } from './hooks/useAppState';
import { Header } from './components/layout/Header';
import { BudgetOverview } from './components/budget/BudgetOverview';
import { LoginForm } from './components/auth/LoginForm';
import './App.css';

// Simple Auth Guard (placeholder)
const PrivateRoute: React.FC<{ children: React.ReactElement }> = ({ children }) => {
  // For now, assume authenticated or skip check until auth is fully integrated
  const isAuthenticated = true; // TODO: Check from state
  return isAuthenticated ? children : <Navigate to="/login" />;
};

function App() {
  return (
    <AppProvider>
      <Router>
        <div className="min-h-screen bg-gray-100">
          <Routes>
            <Route path="/login" element={<LoginForm />} />
            <Route
              path="/*"
              element={
                <PrivateRoute>
                  <>
                    <Header />
                    <Routes>
                      <Route path="/" element={<Navigate to={`/budget/${new Date().toISOString().slice(0, 7)}`} />} />
                      <Route path="/budget/:month" element={<BudgetOverview />} />
                    </Routes>
                  </>
                </PrivateRoute>
              }
            />
          </Routes>
        </div>
      </Router>
    </AppProvider>
  );
}

export default App;
