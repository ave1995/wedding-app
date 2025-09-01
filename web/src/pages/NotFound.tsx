import NotFound from "../assets/NotFound.png";

function NotFoundPage() {
  return (
    <div className="flex flex-col items-center place-content-center font-stretch-50%  w-96 h-screen p-6">
      <h1 className="text-6xl">Ups!</h1>
      <h2 className="text-4xl">Styď se! Tohle jsi neměl najít!</h2>
      <div className="flex items-center gap-3">
        <span className="text-9xl">4</span>
        <img src={NotFound} alt="NotFound" className="w-20 h-20"></img>
        <span className="text-9xl">4</span>
      </div>
      <h3 className="text-2xl font-semibold">STRÁNKA NENALEZENA</h3>
    </div>
  );
}

export default NotFoundPage;
