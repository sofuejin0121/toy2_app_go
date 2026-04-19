import { useId, useState } from 'react';

type PasswordInputProps = {
  id?: string;
  name?: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  required?: boolean;
  autoComplete?: string;
  /** input に追加するクラス（ベースは他テキスト欄と同幅 w-full。枠線の重複は避ける） */
  inputClassName?: string;
  /** 入力欄の下に表示するヒント（例: 表示切替の案内） */
  helperText?: string;
};

/** メール欄などと同じ見た目の幅・枠。右端はトグル用に余白（pr-11） */
const inputBaseClass =
  'w-full border border-gray-300 rounded-lg px-3 py-2 pr-11 text-sm text-gray-900 outline-none placeholder:text-gray-400 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-500';

const toggleClass =
  'absolute right-1 top-1/2 z-10 flex h-9 w-9 -translate-y-1/2 items-center justify-center rounded-md border-0 bg-transparent text-gray-600 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-1';

/**
 * パスワードの表示/非表示トグル。
 * input は他フィールドと同じ w いっぱいの枠、トグルは枠内の右端に重ねて横に置く。
 */
export default function PasswordInput({
  id: idProp,
  name,
  value,
  onChange,
  required,
  autoComplete = 'current-password',
  inputClassName = '',
  helperText,
}: PasswordInputProps) {
  const uid = useId();
  const id = idProp ?? `password-${uid}`;
  const [visible, setVisible] = useState(false);

  const toggleLabel = visible ? 'パスワードを隠す' : 'パスワードを表示';

  const inputClass = [inputBaseClass, inputClassName].filter(Boolean).join(' ');

  return (
    <div className="w-full">
      {/* helperText は relative の外。内側だと absolute の top:50% の基準がずれる */}
      <div className="relative w-full">
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
          className={inputClass}
        />
        <button
          type="button"
          onClick={() => setVisible((v) => !v)}
          title={toggleLabel}
          className={toggleClass}
          aria-label={toggleLabel}
          aria-pressed={visible}
        >
          {visible ? <EyeOffIcon /> : <EyeIcon />}
        </button>
      </div>
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
