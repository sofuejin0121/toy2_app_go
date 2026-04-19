/**
 * アプリケーションのエントリポイント（Vite が最初に読み込む TSX）。
 *
 * 処理の流れ:
 * 1. index.html の `<div id="root">` を取得する（無ければ即失敗させ、不具合に気づきやすくする）
 * 2. `createRoot(rootEl)` でその DOM を「React が描画・更新を担当するルート」に登録する
 * 3. `render(...)` に JSX ツリーを渡し、外側のコンポーネントから順にマウントする
 *
 * JSX の入れ子（外 → 内）の役割:
 * - StrictMode … 開発時のみの追加チェック（副作用の二重実行など）。本番向けの見た目用ではない。
 * - BrowserRouter … ブラウザの URL と React を同期。子孫で Routes / Link / useParams が使える。
 * - Provider (jotai) … atom（例: currentUserAtom）をアプリ全体で共有するコンテキスト。
 * - SWRConfig … `useSWR` 全体に効くデフォルトオプション。ここでは `dedupingInterval` のみ。
 * - AuthLoader … セッション確認など「ログイン状態を揃えてから」子を描画（実装は AuthContext.tsx）。
 * - App … ルート定義（URL とページの対応表）。
 *
 * SWR / SWRConfig（初心者向け補足）:
 * - SWR は「キー（多くは文字列）」と「そのキー用のデータ取得関数」を組み合わせ、結果をキャッシュする。
 * - `<SWRConfig>` で囲んだ子は、ここで渡したデフォルト設定を共有する。
 * - `dedupingInterval: 2000` … 同一キーへのフェッチが 2 秒以内に重なったとき、HTTP をまとめて重複を抑える。
 * - グローバル fetcher は設定していない。各 `useSWR(key, fetcher)` の第 2 引数に API 関数を渡す方針。
 */
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
  /* 開発時の厳格チェック（本番ビルドでも動くが、主に開発者向け） */
  <StrictMode>
    {/* URL と画面を結ぶ土台 → App.tsx の Routes が動く */}
    <BrowserRouter>
      {/* Jotai: useAtom で共有状態にアクセス */}
      <Provider>
        {/* SWR 共通設定（各フックの useSWR がこの範囲の設定を継承） */}
        <SWRConfig value={{ dedupingInterval: 2000 }}>
          {/* 起動時にログイン状態を確定させてからルートツリーを出す */}
          <AuthLoader>
            <App />
          </AuthLoader>
        </SWRConfig>
      </Provider>
    </BrowserRouter>
  </StrictMode>,
);
