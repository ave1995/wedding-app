import { useEffect, useState } from "react";
import { get } from "../../functions/fetch";

export interface SvgItem {
  Name: string;
  URL: string;
}

interface IconSelectorProps {
  onSelect: (icon: SvgItem) => void;
  onClose: () => void;
}

export default function IconSelector({ onSelect, onClose }: IconSelectorProps) {
  const API_BASE_URL = import.meta.env.VITE_API_URL;
  const [svgs, setSvgs] = useState<SvgItem[]>([]);

  useEffect(() => {
    async function fetchIcons() {
      const svgItems = await get<SvgItem[]>(`${API_BASE_URL}/user-svgs`);
      setSvgs(svgItems);
    }
    fetchIcons();
  }, []);

  function selectIcon(icon: SvgItem) {
    onSelect(icon);
    onClose();
  }

  return (
    <div className="p-4 border rounded bg-white shadow-lg">
      <h2 className="mb-2 font-bold">Choose an icon:</h2>
      <div className="flex gap-4 flex-wrap items-center justify-center">
        {svgs.map((svg) => (
          <img
            key={svg.Name}
            src={svg.URL}
            alt={svg.Name}
            className={`w-12 h-12 cursor-pointer border-2 rounded-lg`}
            onClick={() => selectIcon(svg)}
          />
        ))}
      </div>
    </div>
  );
}
