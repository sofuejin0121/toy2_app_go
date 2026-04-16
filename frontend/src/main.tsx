import { Provider } from 'jotai';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App.tsx';
import { AuthLoader } from './contexts/AuthContext';
import './index.css';

const rootEl = document.getElementById('root');
if (!rootEl) throw new Error('Root element not found');

createRoot(rootEl).render(
  <StrictMode>
    <BrowserRouter>
      {/* Provider は jotai の atom をグローバルに管理するためのラッパー */}
      <Provider>
        {/* AuthLoader は起動時に一度だけログイン状態を確認する */}
        <AuthLoader>
          <App />
        </AuthLoader>
      </Provider>
    </BrowserRouter>
  </StrictMode>,
);
