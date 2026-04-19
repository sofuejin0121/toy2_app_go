/**
 * 通知設定（SettingsPage）で使う iOS 風トグル 1 行分。
 *
 * - 見た目の複雑な Tailwind（peer / after 擬似要素）を TRACK_CLASS に閉じ、ページ側は title と checked だけ渡す。
 * - 実体は非表示の checkbox（sr-only）＋ label でクリック領域を確保するパターン。
 */
const TRACK_CLASS =
  'w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-blue-600 after:content-[\'\'] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all';

interface SettingsToggleRowProps {
  title: string;
  description: string;
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
}

export default function SettingsToggleRow({
  title,
  description,
  checked,
  onCheckedChange,
}: SettingsToggleRowProps) {
  return (
    <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
      <div>
        <p className="font-medium text-gray-900 text-sm">{title}</p>
        <p className="text-xs text-gray-500 mt-0.5">{description}</p>
      </div>
      <label className="relative inline-flex items-center cursor-pointer">
        <input
          type="checkbox"
          checked={checked}
          onChange={(e) => onCheckedChange(e.target.checked)}
          className="sr-only peer"
        />
        <div className={TRACK_CLASS} />
      </label>
    </div>
  );
}
