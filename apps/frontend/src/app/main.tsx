import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { createBrowserRouter, RouterProvider} from 'react-router-dom';

import LoginPage from '../pages/auth/LoginPage';
import NotFoundPage from '../pages/notFound/NotFoundPage';

import './index.css'

import AuthProvider from 'react-auth-kit'
import createAuthStore from 'react-auth-kit/store/createAuthStore';

const store = createAuthStore('cookie', {
  authName: '_auth',
  cookieDomain: window.location.hostname,
  cookieSecure: true,
});

const router = createBrowserRouter([
  {path: '/', element: <App />},
  {path: '/login', element: <LoginPage />},
  {path: '*', element: <NotFoundPage />},
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider store={store}>
      <RouterProvider router={router}/>
    </AuthProvider>
  </StrictMode>,
)
