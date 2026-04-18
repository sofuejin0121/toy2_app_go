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

/**
 * アプリのエントリポイント。
 *
 * SWRConfig について（初心者向け）:
 * - SWR は「キー（文字列）」と「そのキーでデータを取ってくる関数」を組み合わせてキャッシュするライブラリです。
 * - <SWRConfig> で囲んだ範囲の子コンポーネントは、同じデフォルト設定（再取得の間隔など）を共有します。
 * - ここでは dedupingInterval のみ指定しています。同じキーへのリクエストが 2 秒以内に重なった場合、
 *   実際の HTTP は 1 回にまとめられる（重複抑制）イメージです。
 * - グローバルな fetcher は設定していません。各 useSWR 呼び出しの第 2 引数で API 関数を渡す方針です。
 */
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
