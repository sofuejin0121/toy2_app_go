import { useId, useState } from 'react';

type PasswordInputProps = {
  id?: string;
  name?: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  required?: boolean;
  autoComplete?: string;
  inputClassName?: string;
  /** 入力欄の下に表示するヒント（例: 表示切替の案内） */
  helperText?: string;
};

const defaultInputClass =
  'w-full border border-gray-300 rounded-lg px-3 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent';

/**
 * パスワードの表示/非表示トグル付き input。
 * ラベルは親で描画する想定（各ページのレイアウトを崩さないため）。
 */
export default function PasswordInput({
  id: idProp,
  name,
  value,
  onChange,
  required,
  autoComplete = 'current-password',
  inputClassName = defaultInputClass,
  helperText,
}: PasswordInputProps) {
  const uid = useId();
  const id = idProp ?? `password-${uid}`;
  const [visible, setVisible] = useState(false);

  const toggleLabel = visible ? 'パスワードを隠す' : 'パスワードを表示';

  return (
    <div className="relative">
      <input
        id={id}
        name={name}
        type={visible ? 'text' : 'password'}
        value={value}
        onChange={onChange}
        required={required}
        autoComplete={autoComplete}
        spellCheck={false}
        aria-describedby={helperText ? `${id}-hint` : undefined}
        className={inputClassName}
      />
      <button
        type="button"
        onClick={() => setVisible((v) => !v)}
        title={toggleLabel}
        className="absolute right-1.5 top-1/2 z-10 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-md border border-gray-200 bg-white text-gray-700 shadow-sm hover:bg-gray-50 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1"
        aria-label={toggleLabel}
        aria-pressed={visible}
      >
        {visible ? <EyeOffIcon /> : <EyeIcon />}
      </button>
      {helperText ? (
        <p className="mt-1 text-xs text-gray-500" id={`${id}-hint`}>
          {helperText}
        </p>
      ) : null}
    </div>
  );
}

function EyeIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden>
      <path
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M1 12s4-7 11-7 11 7 11 7-4 7-11 7-11-7-11-7z"
      />
      <circle cx="12" cy="12" r="3" strokeWidth="2" />
    </svg>
  );
}

function EyeOffIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden>
      <path
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M3 3l18 18M10.6 10.6a3 3 0 004.8 4.8M9.9 5.1A10.4 10.4 0 0112 5c7 0 11 7 11 7a21.5 21.5 0 01-4.1 5.1M6.2 6.2C3.6 8.1 2 11 2 11a21.3 21.3 0 0011 7c1.7 0 3.3-.3 4.8-.8M9.9 9.9A3 3 0 0012 15a3 3 0 002.1-5.1"
      />
    </svg>
  );
}
