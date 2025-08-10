import NotFound from "../assets/NotFound.png";

function NotFoundPage() {
  return (
    <div className="flex flex-col items-center place-content-center font-stretch-50%">
      <h1 className="text-6xl">Oops!</h1>
      <div className="flex items-center gap-3">
        <span className="text-9xl">4</span>
        <img src={NotFound} alt="NotFound" className="w-20 h-20"></img>
        <span className="text-9xl">4</span>
      </div>
      <h2 className="text-2xl font-semibold">PAGE NOT FOUND</h2>
    </div>
  );
}

export default NotFoundPage;
