import { useEffect, useId, useState } from "react";
import Button from "./Button.tsx";

export type FullscreenImageButtonProps = {
  src: string;
  label: string;
  alt?: string;
};

export default function FullscreenImageButton({
  src,
  alt = "",
  label,
}: FullscreenImageButtonProps) {
  const [open, setOpen] = useState(false);
  const titleId = useId();

  useEffect(() => {
    if (!open) return;
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [open]);

  return (
    <>
      <Button label={label} onClick={() => setOpen(true)} type={"Blue"}/>

      {open && (
        <div
          role="dialog"
          aria-modal="true"
          aria-labelledby={titleId}
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
          onClick={() => setOpen(false)}
        >
          <div
            className="relative flex h-[100svh] w-[100svw] items-center justify-center"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="fixed top-10 right-10">
              <Button
                label="Zavřít"
                onClick={() => setOpen(false)}
              />
            </div>

            <img
              id={titleId}
              src={src}
              alt={alt}
              className="max-h-[98svh] max-w-[98svw] select-none rounded-lg object-contain shadow-2xl"
              draggable={false}
            />
          </div>
        </div>
      )}
    </>
  );
}
