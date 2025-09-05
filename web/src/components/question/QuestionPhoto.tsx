interface QuestionPhoto {
  src: string;
  alt: string
}

export default function QuestionPhoto({ src, alt }: QuestionPhoto) {
  return (
    <img
      src={src}
      alt={alt}
      className="h-full w-full object-scale-down"
    />
  );
}
