import { apiUrl } from "../../functions/api";

interface QuestionPhoto {
  path: string;
}

type QuestionPhotoImg = {
  src: string;
  name: string;
};

function pathToBucketQuery(path: string): QuestionPhotoImg | null {
  if (!path) return null;

  const parts = path.split("/");
  const name = parts.pop() || "";
  const bucket = parts.join("/");

  const src = apiUrl(
    `/bucket-data?bucket=${encodeURIComponent(
      bucket
    )}&name=${encodeURIComponent(name)}`
  );

  return { src, name };
}

export default function QuestionPhoto({ path }: QuestionPhoto) {
  const questionPhotoImg = pathToBucketQuery(path);

  if (!questionPhotoImg) {
    return null;
  }

  return (
    <img
      src={questionPhotoImg.src}
      alt={questionPhotoImg.name}
      className="w-full h-full object-contain"
    />
  );
}
