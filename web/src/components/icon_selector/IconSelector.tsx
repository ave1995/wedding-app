import { useEffect, useState } from "react";
import { get } from "../../functions/fetch";
import { apiUrl } from "../../functions/api";
import { useApiErrorHandler } from "../../hooks/useApiErrorHandler";

export type BucketItem = {
  Name: string;
  URL: string;
};

interface IconSelectorProps {
  onSelect: (icon: BucketItem) => void;
  onClose: () => void;
}

export default function IconSelector({ onSelect, onClose }: IconSelectorProps) {
  const [svgs, setSvgs] = useState<BucketItem[]>([]);
  const { handleError } = useApiErrorHandler();

  useEffect(() => {
    async function fetchIcons() {
      const result = await get<BucketItem[]>(apiUrl("/bucket-urls"), {
        bucket: "wedding-user-icons",
        suffix: ".svg",
      });
      if (handleError(result.error, result.status)) return;

      const updatedSvgItems = result.data!.map((item) => ({
        ...item,
        URL: apiUrl(item.URL),
      }));

      setSvgs(updatedSvgItems);
    }
    fetchIcons();
  }, []);

  function selectIcon(icon: BucketItem) {
    onSelect(icon);
    onClose();
  }

  return (
    <div className="p-4 border-2 rounded-lg bg-white shadow-lg border-[#3D52D5]">
      <h2 className="mb-2 font-semibold">Vyber si svojí ikonku:</h2>
      <div className="grid grid-cols-4 gap-4 items-center justify-center">
        {svgs.map((svg) => (
          <img
            key={svg.Name}
            src={svg.URL}
            alt={svg.Name}
            className={`w-12 h-12 cursor-pointer border-2 rounded-lg border-gray-400 hover:border-[#3D52D5] active:scale-95 active:shadow-sm
        transition-all duration-150 ease-in-out`}
            onClick={() => selectIcon(svg)}
          />
        ))}
      </div>
    </div>
  );
}
