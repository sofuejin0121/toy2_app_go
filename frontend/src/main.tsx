import { Provider } from 'jotai';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { SWRConfig } from 'swr';
import App from './App.tsx';
import { AuthLoader } from './contexts/AuthContext';
import './index.css';

const rootEl = document.getElementById('root');
if (!rootEl) throw new Error('Root element not found');

createRoot(rootEl).render(
  <StrictMode>
    <BrowserRouter>
      <Provider>
        <SWRConfig value={{ dedupingInterval: 2000 }}>
          <AuthLoader>
            <App />
          </AuthLoader>
        </SWRConfig>
      </Provider>
    </BrowserRouter>
  </StrictMode>,
);
